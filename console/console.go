package console

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"syscall"

	"github.com/mitchellh/go-homedir"
)

var (
	onlyWhitespace = regexp.MustCompile(`^\s*$`)
	exit           = regexp.MustCompile(`^\s*exit\s*;*\s*$`)
)

const (
	// HistoryFile is the file within the data directory to store input scrollback.
	HistoryFile = "._history"

	// DefaultPrompt is the default prompt line prefix to use for user input querying.
	DefaultPrompt = "> "
)

type Executor interface {
	NewHandler() (handler Handler, err error)
	WordCompleter() []string
}

type Handler interface {
	RunCommand(input string) (err error)
	Close()
}

type Console struct {
	prompt   string       // Input prompt prefix string
	prompter UserPrompter // Input prompter to allow interactive user feedback
	histPath string       // Absolute path to the console scrollback history
	history  []string     // Scroll history maintained by the console
	handler  Handler
	abort    chan os.Signal
}

func New(exec Executor, opts ...func(*Console)) (c *Console, err error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	handler, err := exec.NewHandler()
	if err != nil {
		return nil, err
	}
	// Initialize the console and return
	console := &Console{
		prompt:   DefaultPrompt,
		prompter: Stdin,
		histPath: filepath.Join(home, HistoryFile),
		handler:  handler,
	}
	for _, opt := range opts {
		opt(console)
	}
	// Configure the console's input prompter for scrollback and tab completion
	if content, err := ioutil.ReadFile(console.histPath); err != nil {
		console.prompter.SetHistory(nil)
	} else {
		console.history = strings.Split(string(content), "\n")
		console.prompter.SetHistory(console.history)
	}
	console.SetWordCompleter(exec.WordCompleter())
	return console, nil
}

func WithPrompt(prompt string) func(*Console) {
	return func(console *Console) {
		console.prompt = prompt
	}
}

func WithHistoryFile(file string) func(*Console) {
	return func(console *Console) {
		home, _ := homedir.Dir()
		console.histPath = filepath.Join(home, file)
	}
}

func (c *Console) SetWordCompleter(words []string) {
	c.prompter.SetWordCompleter(func(line string, pos int) (head string, completions []string, tail string) {
		if len(line) == 0 || pos == 0 {
			for _, word := range words {
				if strings.Index(word, ".") == -1 {
					completions = append(completions, word)
				}
			}
			return head, completions, tail
		}
		for _, word := range words {
			if strings.HasPrefix(strings.ToLower(word), strings.ToLower(string([]rune(line)[:pos]))) {
				completions = append(completions, word)
			}
		}
		return head, completions, string([]rune(line)[pos:])
	})
}

// Interactive starts an interactive user session, where input is propted from
// the configured user prompter.
func (c *Console) Interactive() {
	defer c.SaveHistory()
	defer c.handler.Close()
	var scheduler = make(chan string) // Channel to send the next prompt on and receive the input

	// Start a goroutine to listen for prompt requests and send back inputs
	go func() {
		for {
			// Read the next user input
			line, err := c.prompter.PromptInput(<-scheduler)
			if err != nil {
				// In case of an error, either clear the prompt or fail
				if err == abortedErr { // ctrl-C
					scheduler <- ""
					continue
				}
				close(scheduler)
				return
			}
			// User input retrieved, send for interpretation and loop
			scheduler <- line
		}
	}()
	// Monitor Ctrl-C too in case the input is empty and we need to bail
	c.abort = make(chan os.Signal, 1)
	signal.Notify(c.abort, syscall.SIGINT, syscall.SIGTERM)

	// Start sending prompts to the user and reading back inputs
	for {
		// Send the next prompt, triggering an input read and process the result
		scheduler <- c.prompt
		select {
		case <-c.abort:
			close(scheduler)
			// User forcefully quite the console
			fmt.Println("\n exiting . . .")
			return

		case input, ok := <-scheduler:
			// User input was returned by the prompter, handle special cases
			if !ok || exit.MatchString(input) {
				return
			}
			if onlyWhitespace.MatchString(input) {
				continue
			}
			if err := c.Execute(input); err != nil {
				fmt.Printf("input:[ %s ], execute err:%s", input, err.Error())
			}
		}
	}
}

// Evaluate executes func and pretty prints the result to the specified output stream.
func (c *Console) Execute(input string) error {
	appendHistory := func() {
		if len(input) > 0 && input[0] != ' ' {
			if command := strings.TrimSpace(input); len(c.history) == 0 || command != c.history[len(c.history)-1] {
				c.history = append(c.history, command)
				if c.prompter != nil {
					c.prompter.AppendHistory(command)
				}
			}
		}
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("input:[ %s ], fatal err: %v\n", input, r)
		} else {
			appendHistory()
		}
	}()
	return c.handler.RunCommand(input)
}

func (c *Console) Close() {
	if c.abort != nil {
		select {
		case <-c.abort:
			return
		default:
			c.abort <- syscall.SIGQUIT
		}
	}
}

func (c *Console) SaveHistory() error {
	if err := ioutil.WriteFile(c.histPath, []byte(strings.Join(c.history, "\n")), 0600); err != nil {
		return err
	}
	if err := os.Chmod(c.histPath, 0600); err != nil { // Force 0600, even if it was different previously
		return err
	}
	return nil
}

func (c *Console) ClearHistory() {
	c.history = nil
	c.prompter.ClearHistory()
	if err := os.Remove(c.histPath); err != nil {
		fmt.Println("can't delete history file:", err)
	} else {
		fmt.Println("history file deleted")
	}
}

func ParseInputData(input string) (field, method string, params []reflect.Value) {
	var t int

	for i, k := range input {
		if k == '.' {
			field = input[t:i]
			t = i + 1
			continue
		}
		if k == '(' {
			method = input[t:i]
			t = i + 1
			continue
		}
		if k == ',' || k == ')' {
			if param := input[t:i]; param != "" {
				params = append(params, reflect.ValueOf(param))
				t = i + 1
			}
		}
	}
	return
}

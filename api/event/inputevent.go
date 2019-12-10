package event

import (
	"bufio"
	"os"
)

type inputEvent struct {
	done chan bool
}

func newInputEvent() inputEvent {
	return inputEvent{done: make(chan bool)}
}

// WaitForEnter waits until the user presses Enter
func (c *inputEvent) WaitForEnter() chan bool {
	go c.readFromCLI()
	return c.done
}

func (c *inputEvent) readFromCLI() {
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	c.done <- true
}

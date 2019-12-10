package task

// MultiTask contains a set of Tasks to be invoked synchronously
type MultiTask struct {
	tasks       []Task
	completedCB func()
}

// New creates a new multi task
func NewMultiTask(completedCB func()) *MultiTask {
	return &MultiTask{
		completedCB: completedCB,
	}
}

// Add adds a task
func (m *MultiTask) Add(task Task) {
	m.tasks = append(m.tasks, task)
}

// Invoke invokes the task
func (m *MultiTask) Invoke() {
	defer m.completedCB()

	for _, task := range m.tasks {
		task.Invoke()
	}
}

// Attempts returns the number of invocation attempts that were made
// in order to achieve a successful response
func (m *MultiTask) Attempts() int {
	var attempts int
	for _, t := range m.tasks {
		attempts += t.Attempts()
	}
	return attempts
}

// LastError returns the last error that was recorder
func (m *MultiTask) LastError() error {
	for _, t := range m.tasks {
		if t.LastError() != nil {
			return t.LastError()
		}
	}
	return nil
}

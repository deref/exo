package task

type Definition interface {
	Name() string
	RunTask(task *Task) error
}

func NewTaskFunc(name string, run func(task *Task) error) Definition {
	return &taskFunc{
		name: name,
		run:  run,
	}
}

type taskFunc struct {
	name string
	run  func(*Task) error
}

func (tf *taskFunc) Name() string {
	return tf.name
}

func (tf *taskFunc) RunTask(task *Task) error {
	return tf.run(task)
}

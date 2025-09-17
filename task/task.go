package task

type Task interface {
	Run() error
}

type RootPrivilegeTask struct {
}

func (t *RootPrivilegeTask) Run() error {
	return nil
}

package exec

type ExecEngine interface {
	DeleteApp(app string, namespace string) error
	DeleteGroup(group string, namespace string) error
	DeployApp(app string, appDir string, namespace string) error
	DeployGroup(group string, appGroupDir string, namespace string) error
}

type Opts struct {
	ExtraArgs []string
}

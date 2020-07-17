package template

const (
	stdinPath = "-"
)

type TemplateEngine interface {
	Template(name string, inputFilePath string, outputFilePath string) error
	Lint(name string, inputFilePath string) error
}

type Opts struct {
	Environment string
	ExtraArgs   []string
}

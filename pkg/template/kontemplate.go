package template

import (
	"github.com/DataReply/korgi/pkg"
	"github.com/codeskyblue/go-sh"
)

type KontemplateEngine struct {
	Opts Opts
}

func NewKontemplateEngine(Opts Opts) *KontemplateEngine {
	return &KontemplateEngine{Opts}
}
func (e *KontemplateEngine) Template(name string, inputFilePath string, outputFilePath string) error {

	inputArgs := pkg.ExplodeArg(append([]string{"template", inputFilePath, "-o", outputFilePath}, e.Opts.ExtraArgs...))

	err := sh.Command("kontemplate", inputArgs...).Run()
	if err != nil {
		return err
	}
	return nil
}
func (e *KontemplateEngine) Lint(name string, inputFilePath string) error {

	return nil
}

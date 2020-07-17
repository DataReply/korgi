package template

import (
	"fmt"

	"github.com/DataReply/kapply/pkg"
	"github.com/codeskyblue/go-sh"
)

type HelmFileEngine struct {
	Opts Opts
}

func NewHelmFileEngine(Opts Opts) *HelmFileEngine {
	return &HelmFileEngine{Opts}
}
func (e *HelmFileEngine) Template(name string, inputFilePath string, outputFilePath string) error {
	inputArgs := pkg.ExplodeArg(append([]string{"--environment", e.Opts.Environment,
		"--file", inputFilePath, "--state-values-set", fmt.Sprintf("app=%s", name),
		"template", "--output-dir", outputFilePath}, e.Opts.ExtraArgs...))

	err := sh.Command("helmfile", inputArgs...).Run()
	if err != nil {
		return err
	}
	return nil
}
func (e *HelmFileEngine) Lint(name string, inputFilePath string) error {

	inputArgs := pkg.ExplodeArg(append([]string{"--environment", e.Opts.Environment,
		"--file", inputFilePath, "--state-values-set", fmt.Sprintf("app=%s", name), "lint"}, e.Opts.ExtraArgs...))

	err := sh.Command("helmfile", inputArgs...).Run()
	if err != nil {
		return err
	}
	return nil
}

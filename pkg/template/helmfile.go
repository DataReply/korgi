package template

import (
	"fmt"

	"github.com/codeskyblue/go-sh"
)

type HelmFileEngine struct {
	genericOpts GenericOpts
}

func NewHelmFileEngine(genericOpts GenericOpts) *HelmFileEngine {
	return &HelmFileEngine{genericOpts}
}
func (e *HelmFileEngine) Template(name string, inputFilePath string, outputFilePath string) error {

	err := sh.Command("helmfile", "--environment", e.genericOpts.Environment,
		"--file", inputFilePath, "--state-values-set", fmt.Sprintf("app=%s", name), "template", "--output-dir", outputFilePath).Run()
	if err != nil {
		return err
	}
	return nil
}
func (e *HelmFileEngine) Lint(name string, inputFilePath string) error {

	// helmfile --environment "$ENVIRONMENT" --file "$app" --state-values-set app="$_app" template --output-dir "$app_group_dir/$_app"

	err := sh.Command("helmfile", "--environment", e.genericOpts.Environment,
		"--file", inputFilePath, "--state-values-set", fmt.Sprintf("app=%s", name), "lint").Run()
	if err != nil {
		return err
	}
	return nil
}

/*
Copyright © 2020  Artyom Topchyan a.topchyan@reply.de

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package template

import (
	"fmt"

	"github.com/DataReply/korgi/pkg"
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

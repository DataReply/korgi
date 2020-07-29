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
	"os/exec"

	"github.com/DataReply/korgi/pkg/utils"
	"github.com/go-logr/logr"
)

type HelmFileEngine struct {
	Opts Opts
	log  logr.Logger
}

func NewHelmFileEngine(Opts Opts, log logr.Logger) *HelmFileEngine {
	return &HelmFileEngine{Opts, log}
}

func (e *HelmFileEngine) Template(name string, inputFilePath string, outputFilePath string) error {
	inputArgs := append([]string{"--environment", e.Opts.Environment,
		"--file", inputFilePath, "--state-values-set", fmt.Sprintf("app=%s", name),
		"template", "--output-dir", outputFilePath}, e.Opts.ExtraArgs...)

	cmd := exec.Command("helmfile", inputArgs...)
	print := func(in string) {
		e.log.Info(in)
	}

	err := utils.ExecWithOutput(cmd, print, print)

	if err != nil {
		return err
	}

	return nil
}
func (e *HelmFileEngine) Lint(name string, inputFilePath string) error {

	inputArgs := append([]string{"--environment", e.Opts.Environment,
		"--file", inputFilePath, "--state-values-set", fmt.Sprintf("app=%s", name), "lint"}, e.Opts.ExtraArgs...)

	cmd := exec.Command("helmfile", inputArgs...)

	print := func(in string) {
		e.log.Info(in)
	}

	err := utils.ExecWithOutput(cmd, print, print)

	if err != nil {
		return err
	}

	return nil
}

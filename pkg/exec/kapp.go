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
package exec

import (
	"fmt"
	"os/exec"

	"github.com/DataReply/korgi/pkg/utils"
	"github.com/go-logr/logr"
)

type KappEngine struct {
	Opts Opts
	log  logr.Logger
}

func NewKappEngine(opts Opts, log logr.Logger) *KappEngine {
	return &KappEngine{opts, log}
}

func (e *KappEngine) DeleteApp(app string, namespace string) error {

	inputArgs := append([]string{"delete", "-a", app, "-n", namespace, fmt.Sprintf("--diff-run=%t", e.Opts.DiffRun)}, e.Opts.ExtraArgs...)

	cmd := exec.Command("kapp", inputArgs...)
	print := func(in string) {
		e.log.Info(in)
	}

	err := utils.ExecWithOutput(cmd, print, print)

	if err != nil {
		return err
	}

	return nil

}

func (e *KappEngine) DeleteGroup(group string, namespace string) error {

	inputArgs := append([]string{"app-group", "delete", "-n", namespace, "-g", group, fmt.Sprintf("--diff-run=%t", e.Opts.DiffRun)}, e.Opts.ExtraArgs...)

	cmd := exec.Command("kapp", inputArgs...)
	print := func(in string) {
		e.log.Info(in)
	}

	err := utils.ExecWithOutput(cmd, print, print)

	if err != nil {
		return err
	}

	return nil
}

func (e *KappEngine) DeployApp(app string, appDir string, namespace string) error {

	inputArgs := append([]string{"deploy", "-a", app, "-f", appDir, "-n", namespace, fmt.Sprintf("--diff-run=%t", e.Opts.DiffRun)}, e.Opts.ExtraArgs...)

	cmd := exec.Command("kapp", inputArgs...)
	print := func(in string) {
		e.log.Info(in)
	}

	err := utils.ExecWithOutput(cmd, print, print)

	if err != nil {
		return err
	}

	return nil

}

func (e *KappEngine) DeployGroup(group string, appGroupDir string, namespace string) error {

	inputArgs := append([]string{"app-group", "deploy", "-d", appGroupDir, "-n", namespace, "-g", group, fmt.Sprintf("--diff-run=%t", e.Opts.DiffRun)}, e.Opts.ExtraArgs...)

	cmd := exec.Command("kapp", inputArgs...)
	print := func(in string) {
		e.log.Info(in)
	}

	err := utils.ExecWithOutput(cmd, print, print)

	if err != nil {
		return err
	}

	return nil
}

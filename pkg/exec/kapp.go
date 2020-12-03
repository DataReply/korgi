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

var inputArgs []string

func NewKappEngine(opts Opts, log logr.Logger) *KappEngine {
	if !opts.AskForConfirmation {
		inputArgs = append(inputArgs, "-y")
	}
	return &KappEngine{opts, log}
}

func (e *KappEngine) DeleteApp(app string, namespace string) error {
	return e.exec(append(inputArgs, []string{"delete", "-a", app, "-n", namespace, fmt.Sprintf("--diff-run=%t", e.Opts.DiffRun)}...))
}

func (e *KappEngine) DeleteGroup(group string, namespace string) error {
	return e.exec(append(inputArgs, []string{"app-group", "delete", "-n", namespace, "-g", group, fmt.Sprintf("--diff-run=%t", e.Opts.DiffRun)}...))
}

func (e *KappEngine) DeployApp(app string, appDir string, namespace string) error {
	return e.exec(append(inputArgs, []string{"deploy", "-a", app, "-f", appDir, "-n", namespace, fmt.Sprintf("--diff-run=%t", e.Opts.DiffRun)}...))
}

func (e *KappEngine) DeployGroup(group string, appGroupDir string, namespace string) error {
	return e.exec(append(inputArgs, []string{"app-group", "deploy", "-d", appGroupDir, "-n", namespace, "-g", group, fmt.Sprintf("--diff-run=%t", e.Opts.DiffRun)}...))
}

func (e *KappEngine) exec(inputArgs []string) error {

	cmd := exec.Command("kapp", append(inputArgs, e.Opts.ExtraArgs...)...)
	print := func(in string) {
		e.log.Info(in)
	}

	var err error
	if e.Opts.AskForConfirmation {
		err = utils.ExecWithStdInOut(cmd)
	} else {
		err = utils.ExecWithOutput(cmd, print, print)
	}

	if err != nil {
		return err
	}

	return nil
}

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
	"github.com/DataReply/korgi/pkg"
	"github.com/codeskyblue/go-sh"
)

type KappEngine struct {
	Opts Opts
}

func NewKappEngine(opts Opts) *KappEngine {
	return &KappEngine{Opts: opts}
}

func (e *KappEngine) DeleteApp(app string, namespace string) error {

	inputArgs := pkg.ExplodeArg(append([]string{"-y", "delete", "-a", app, "-n", namespace}, e.Opts.ExtraArgs...))

	err := sh.Command("kapp", inputArgs...).Run()
	if err != nil {
		return err
	}
	return nil
}

func (e *KappEngine) DeleteGroup(group string, namespace string) error {

	inputArgs := pkg.ExplodeArg(append([]string{"-y", "app-group", "delete", "-n", namespace, "-g", group}, e.Opts.ExtraArgs...))

	err := sh.Command("kapp", inputArgs...).Run()
	if err != nil {
		return err
	}
	return nil
}

func (e *KappEngine) DeployApp(app string, appDir string, namespace string) error {

	inputArgs := pkg.ExplodeArg(append([]string{"-y", "deploy", "-a", app, "-f", appDir, "-n", namespace}, e.Opts.ExtraArgs...))

	err := sh.Command("kapp", inputArgs...).Run()
	if err != nil {
		return err
	}
	return nil
}

func (e *KappEngine) DeployGroup(group string, appGroupDir string, namespace string) error {

	inputArgs := pkg.ExplodeArg(append([]string{"-y", "app-group", "deploy", "-d", appGroupDir, "-n", namespace, "-g", group}, e.Opts.ExtraArgs...))

	err := sh.Command("kapp", inputArgs...).Run()
	if err != nil {
		return err
	}
	return nil
}

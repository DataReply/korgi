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

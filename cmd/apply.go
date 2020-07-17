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
package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/DataReply/kapply/pkg"
	"github.com/DataReply/kapply/pkg/template"
	"github.com/spf13/cobra"
)

var engine template.TemplateEngine

var targetNamespaceDir string
var namespaceDir string

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply resources to k8s",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		dir, _ := os.Getwd()

		realmDir := pkg.ConcatDir(dir, "realm")
		namespacesDir := pkg.ConcatDir(realmDir, "namespaces")
		namespaceDir = pkg.ConcatDir(namespacesDir, *namespace)

		if _, err := os.Stat(realmDir); os.IsNotExist(err) {
			return errors.New("realm directory does not exist")
		}
		if _, err := os.Stat(namespacesDir); os.IsNotExist(err) {
			return errors.New("namespaces directory does not exist")
		}
		if _, err := os.Stat(namespaceDir); os.IsNotExist(err) {
			return fmt.Errorf("%s directory does not exist", namespaceDir)
		}

		parentDir := fmt.Sprintf("/tmp/kapp/%s", time.Now().Format("2006-01-02/15-04:05"))

		targetNamespaceDir = pkg.ConcatDir(parentDir, *namespace)

		err := os.MkdirAll(parentDir, os.ModePerm)
		if err != nil {
			return errors.New("cannot create working directory")
		}

		_engineName, _ := cmd.Flags().GetString("template-engine")

		switch e := _engineName; e {
		case "helmfile":
			engine = template.NewHelmFileEngine(template.GenericOpts{
				Environment: *environment,
				Namespace:   *environment,
			})
		case "kontemplate":
			engine = template.NewHelmFileEngine(template.GenericOpts{
				Environment: *environment,
				Namespace:   *environment,
			})

		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.PersistentFlags().StringP("template-engine", "t", "helmfile", "Template engine")
	applyCmd.PersistentFlags().BoolP("lint", "l", false, "Lint temlate")
	applyCmd.PersistentFlags().BoolP("dry-run", "d", false, "Dry Run")

}

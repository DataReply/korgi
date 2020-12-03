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
	"os"
	"path/filepath"

	"github.com/DataReply/korgi/pkg/utils"
	"github.com/spf13/cobra"
)

// namespaceCmd represents the namespace command
var namespaceCmd = &cobra.Command{
	Use:   "apply-namespace",
	Short: "Namespace apply",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		lint, _ := cmd.Flags().GetBool("lint")

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		appFilter, _ := cmd.Flags().GetString("app")

		outputDir, _ := cmd.Flags().GetString("output-dir")

		isolated, _ := cmd.Flags().GetBool("isolate")

		askForConfirmation, _ := cmd.Flags().GetBool("ask-for-confirmation")

		namespace := args[0]

		namespaceDir := utils.GetNamespaceDir(namespace)
		if _, err := os.Stat(namespaceDir); os.IsNotExist(err) {
			return errors.New("namespaces directory does not exist")
		}

		err := filepath.Walk(namespaceDir, func(path string, info os.FileInfo, err error) error {

			if info.IsDir() && path != namespaceDir {
				group := filepath.Base(path)
				err := applyAppGroup(group, namespace, getFinalOutputDir(outputDir, isolated), appFilter, lint, dryRun, askForConfirmation)
				if err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(namespaceCmd)

	namespaceCmd.Flags().BoolP("lint", "l", false, "Lint temlate")
	namespaceCmd.Flags().BoolP("dry-run", "d", false, "Dry Run")

}

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
	"fmt"
	"os"
	"path/filepath"

	"github.com/DataReply/kapply/pkg"
	"github.com/spf13/cobra"
)

// deleteNamespaceCmd represents the namespace command
var deleteNamespaceCmd = &cobra.Command{
	Use:   "delete-namespace",
	Short: "Namespace apply",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filter, _ := cmd.Flags().GetString("filter")

		namespace := args[0]

		namespaceDir := pkg.GetNamespaceDir(namespace)
		if _, err := os.Stat(namespaceDir); os.IsNotExist(err) {
			return fmt.Errorf("%s directory does not exist", namespaceDir)
		}

		err := filepath.Walk(namespaceDir, func(path string, info os.FileInfo, err error) error {

			if info.IsDir() && path != namespaceDir {
				group := filepath.Base(path)
				err := deleteAppGroup(group, namespace, filter)
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
	rootCmd.AddCommand(deleteNamespaceCmd)

	deleteNamespaceCmd.Flags().BoolP("lint", "l", false, "Lint temlate")
	deleteNamespaceCmd.Flags().BoolP("dry-run", "d", false, "Dry Run")

}

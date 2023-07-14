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
	"github.com/spf13/cobra"
)

func deleteAppGroup(group string, namespace string, appFilter string) error {
	if appFilter != "" {
		err := kappEngine.DeleteApp(group+"-"+appFilter, namespace)
		if err != nil {
			return fmt.Errorf("kapp app delete: %w", err)
		}
		return nil
	}

	err := kappEngine.DeleteGroup(group, namespace)
	if err != nil {
		return fmt.Errorf("kapp group delete: %w", err)
	}
	return nil
}

// deleteCmd represents the delete command
var deleteGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "Delete app group or app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		group := args[0]

		namespace, _ := cmd.Flags().GetString("namespace")

		appFilter, _ := cmd.Flags().GetString("app")

		err := deleteAppGroup(group, namespace, appFilter)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	deleteCmd.AddCommand(deleteGroupCmd)
}

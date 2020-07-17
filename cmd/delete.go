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

	"github.com/DataReply/kapply/pkg/kapp"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete app group or app",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		group := args[0]

		filter, _ := cmd.Flags().GetString("app")
		if filter != "" {
			err := kapp.DeleteApp(group+"-"+filter, *namespace, nil)
			if err != nil {
				return fmt.Errorf("kapp app delete: %w", err)
			}
			return nil
		}

		err := kapp.DeleteGroup(group, *namespace, nil)
		if err != nil {
			return fmt.Errorf("kapp group delete: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringP("app", "a", "", "Filter single app from this group")

}

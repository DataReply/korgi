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
	"github.com/spf13/cobra"
	"os"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:              "apply",
	Short:            "Apply resources to k8s",
	Args:             cobra.ExactArgs(1),
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		namespace, _ := cmd.Flags().GetString("namespace")
		lint, _ := cmd.Flags().GetBool("lint")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		appFilter, _ := cmd.Flags().GetString("app")
		outputDir, _ := cmd.Flags().GetString("output-dir")
		isolated, _ := cmd.Flags().GetBool("isolate")
		askForConfirmation, _ := cmd.Flags().GetBool("ask-for-confirmation")

		if !askForConfirmation {
			toContinue, errAsking := delYN(os.Stdin)
			switch {
			case errAsking != nil:
				return errAsking
			case !toContinue:
				os.Exit(0)
			}
		}

		err := applyAppGroup(args[0], namespace, getFinalOutputDir(outputDir, isolated), appFilter, lint, dryRun, defaultMatcher, askForConfirmation)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.PersistentFlags().BoolP("lint", "l", false, "Lint temlate")
	applyCmd.PersistentFlags().BoolP("dry-run", "d", false, "Dry Run")
	applyCmd.PersistentFlags().StringP("namespace", "n", "", "Target namespace")
	applyCmd.MarkFlagRequired("namespace")

}

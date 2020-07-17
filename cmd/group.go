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
	"github.com/DataReply/kapply/pkg/kapp"
	"github.com/spf13/cobra"
)

var filterApp *string

func templateApp(app string, inputFilePath string, appGroupDir string, lint bool) error {

	targeAppDir := pkg.ConcatDir(appGroupDir, app)

	err := os.MkdirAll(targeAppDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("create app dir: %w", err)
	}
	if lint {
		err = engine.Lint(app, inputFilePath)
		if err != nil {
			return fmt.Errorf("linting: %w", err)
		}
	}
	err = engine.Template(app, inputFilePath, targeAppDir)
	if err != nil {
		return fmt.Errorf("templating: %w", err)
	}

	return nil
}

func iterateOnAppGroup(group string, lint bool, dryRun bool) error {
	appGroupDir := pkg.ConcatDir(namespaceDir, group)
	targetAppGroupDir := pkg.ConcatDir(targetNamespaceDir, group)

	err := os.MkdirAll(targetAppGroupDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("create group directory: %w", err)
	}

	matches, _ := filepath.Glob(fmt.Sprintf("%s/*", appGroupDir))

	for _, matchedAppFile := range matches {
		appFile := filepath.Base(matchedAppFile)
		if appFile != "_app_group.yaml" {
			app := pkg.SanitzeAppName(appFile)
			if *filterApp != "" {

				if app != *filterApp {
					continue
				}
			}

			err = templateApp(app, matchedAppFile, targetAppGroupDir, lint)
			if err != nil {
				return fmt.Errorf("apply app: %w", err)
			}

		}

	}
	if !dryRun {
		if *filterApp != "" {
			err = kapp.DeployApp(group+"-"+*filterApp, targetAppGroupDir+"/"+*filterApp, *namespace, nil)
			if err != nil {
				return fmt.Errorf("kapp app deploy: %w", err)
			}
			return nil
		}

		err = kapp.DeployGroup(group, targetAppGroupDir, *namespace, nil)
		if err != nil {
			return fmt.Errorf("kapp group deploy: %w", err)
		}
	}

	return nil
}

// groupCmd represents the group command
var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "App group command",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		lint, _ := cmd.Flags().GetBool("lint")

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		err := iterateOnAppGroup(args[0], lint, dryRun)
		if err != nil {
			return fmt.Errorf("iterate app group: %w", err)
		}

		return nil
	},
}

func init() {
	applyCmd.AddCommand(groupCmd)

	filterApp = groupCmd.Flags().StringP("app", "a", "", "Filter single app from this group")

}

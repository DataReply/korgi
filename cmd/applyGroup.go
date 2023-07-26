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
	"regexp"

	"github.com/DataReply/korgi/pkg/utils"
	"github.com/spf13/cobra"
)

func templateApp(app string, namespace string, inputFilePath string, appGroupDir string, lint bool) error {

	targetAppDir := utils.ConcatDirs(appGroupDir, app)

	err := os.MkdirAll(targetAppDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("creating app dir: %w", err)
	}
	if lint {
		err = helmfileEngine.Lint(app, inputFilePath)
		if err != nil {
			return err
		}
	}
	err = helmfileEngine.Template(app, namespace, inputFilePath, targetAppDir)
	if err != nil {
		return err
	}

	return nil
}

func getFinalOutputDir(outputDir string, isolated bool) string {
	if isolated {
		return utils.ConcatDirs(outputDir, execTime.Format("2006-01-02/15-04:05"))
	}
	return outputDir
}

func applyAppGroup(group string, namespace string, outputDir string, appFilter string, lint bool, dryRun bool, match string, askForConfirmation bool) error {

	log.V(0).Info("applying", "group", group, "namespace", namespace, "app", appFilter, "lint", lint, "dry", dryRun)
	namespaceDir := utils.GetNamespaceDir(namespace)
	if _, err := os.Stat(namespaceDir); os.IsNotExist(err) {
		return fmt.Errorf("%s directory does not exist", namespaceDir)
	}

	appGroupDir := utils.ConcatDirs(namespaceDir, group)

	targetAppGroupDir := utils.ConcatDirs(outputDir, namespace, group)

	err := os.MkdirAll(targetAppGroupDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("creating group directory: %w", err)
	}

	matches, err := filepath.Glob(appGroupDir + "/*")
	if err != nil {
		return fmt.Errorf("listing group directory: %w", err)
	}

	matchers, _ := regexp.Compile(match)

	for _, matchedAppFile := range matches {
		appFile := filepath.Base(matchedAppFile)
		if appFile != "_app_group.yaml" {
			app := utils.SanitizeAppName(appFile, match)
			if appFilter != "" {

				if app != appFilter {
					continue
				}
			}

			if matchers.MatchString(appFile) {
				err = templateApp(app, namespace, matchedAppFile, targetAppGroupDir, lint)
				if err != nil {
					return fmt.Errorf("templating app: %w", err)
				}
			}

		}

	}
	if !dryRun {
		if appFilter != "" {
			err = kappEngine.DeployApp(group+"-"+appFilter, utils.ConcatDirs(targetAppGroupDir, appFilter), namespace)
			if err != nil {
				return fmt.Errorf("running kapp deploy with appFilter: %w", err)
			}
			return nil
		}

		err = kappEngine.DeployGroup(group, targetAppGroupDir, namespace)
		if err != nil {
			return fmt.Errorf("running kapp deploy: %w", err)
		}
	}

	return nil
}

func runApplyWithMatch(cmd *cobra.Command, args []string, match string) error {

	namespace, _ := cmd.Flags().GetString("namespace")
	lint, _ := cmd.Flags().GetBool("lint")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	appFilter, _ := cmd.Flags().GetString("app")
	outputDir, _ := cmd.Flags().GetString("output-dir")
	isolated, _ := cmd.Flags().GetBool("isolate")
	askForConfirmation, _ := cmd.Flags().GetBool("ask-for-confirmation")

	err := applyAppGroup(args[0], namespace, getFinalOutputDir(outputDir, isolated), appFilter, lint, dryRun, match, askForConfirmation)
	if err != nil {
		return err
	}
	return nil
}

// applyCmd represents the apply command
var applyGroupCmd = &cobra.Command{
	Use:              "group",
	Short:            "App-group scoped apply",
	Args:             cobra.ExactArgs(1),
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		err := runApplyWithMatch(cmd, args, defaultMatcher)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	applyCmd.AddCommand(applyGroupCmd)
}

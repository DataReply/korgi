/*
Copyright © 2020 Artyom Topchyan a.topchyan@reply.de

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
	"time"

	"github.com/DataReply/korgi/pkg/exec"
	"github.com/DataReply/korgi/pkg/template"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var execTime time.Time
var templateEngine template.TemplateEngine
var execEngine exec.ExecEngine

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "korgi",
	Short: "DRY Kubernetes Deployments with kapp, helmfile and kontemplate",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		_engineName, _ := cmd.Flags().GetString("template-engine")
		templateExtraArgs, _ := cmd.Flags().GetStringArray("template-engine-args")
		execExtraArgs, _ := cmd.Flags().GetStringArray("exec-engine-args")

		environment, _ := cmd.Flags().GetString("environment")

		execEngine = exec.NewKappEngine(exec.Opts{
			ExtraArgs: execExtraArgs,
		})

		switch e := _engineName; e {
		case "helmfile":
			templateEngine = template.NewHelmFileEngine(template.Opts{
				Environment: environment,
				ExtraArgs:   templateExtraArgs,
			})
		case "kontemplate":
			templateEngine = template.NewKontemplateEngine(template.Opts{
				Environment: environment,
				ExtraArgs:   templateExtraArgs,
			})
		default:
			return fmt.Errorf("%s template engine is not supported", _engineName)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %s \n", err)
		os.Exit(1)
	}
}

func init() {

	execTime = time.Now()
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("environment", "e", "", "Target environment")

	rootCmd.PersistentFlags().StringP("working-dir", "w", "/tmp/kapp", "Working directory")
	rootCmd.PersistentFlags().StringP("filter", "f", "", "Filter to include a single app")
	rootCmd.PersistentFlags().StringP("template-engine", "t", "helmfile", "Template engine")

	rootCmd.PersistentFlags().StringArray("exec-engine-args", []string{}, "Execution engine extra args(only kapp is supported)s")
	rootCmd.PersistentFlags().StringArray("template-engine-args", []string{}, "Template engine extra args")

	rootCmd.MarkFlagRequired("environment")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".korgi" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".korgi")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

package kapp

import (
	"github.com/codeskyblue/go-sh"
)

func explodeArg(args []string, extraArgs []interface{}) []interface{} {
	new := make([]interface{}, len(args)+len(extraArgs))

	for i, v := range args {
		new[i] = v
	}
	return append(new, extraArgs...)
}

func DeleteApp(app string, namespace string, extraArgs []interface{}) error {

	inputArgs := explodeArg([]string{"-y", "delete", "-a", app, "-n", namespace}, extraArgs)

	err := sh.Command("kapp", inputArgs...).Run()
	if err != nil {
		return err
	}
	return nil
}

func DeleteGroup(group string, namespace string, extraArgs []interface{}) error {

	inputArgs := explodeArg([]string{"-y", "app-group", "delete", "-n", namespace, "-g", group}, extraArgs)

	err := sh.Command("kapp", inputArgs...).Run()
	if err != nil {
		return err
	}
	return nil
}

func DeployApp(app string, appDir string, namespace string, extraArgs []interface{}) error {

	inputArgs := explodeArg([]string{"-y", "deploy", "-a", app, "-f", appDir, "-n", namespace}, extraArgs)

	err := sh.Command("kapp", inputArgs...).Run()
	if err != nil {
		return err
	}
	return nil
}

func DeployGroup(group string, appGroupDir string, namespace string, extraArgs []interface{}) error {

	inputArgs := explodeArg([]string{"-y", "app-group", "deploy", "-d", appGroupDir, "-n", namespace, "-g", group}, extraArgs)

	err := sh.Command("kapp", inputArgs...).Run()
	if err != nil {
		return err
	}
	return nil
}

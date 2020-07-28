package exec

import (
	"github.com/DataReply/korgi/pkg"
	"github.com/codeskyblue/go-sh"
)

type KappEngine struct {
	Opts Opts
}

func NewKappEngine(opts Opts) *KappEngine {
	return &KappEngine{Opts: opts}
}

func (e *KappEngine) DeleteApp(app string, namespace string) error {

	inputArgs := pkg.ExplodeArg(append([]string{"-y", "delete", "-a", app, "-n", namespace}, e.Opts.ExtraArgs...))

	err := sh.Command("kapp", inputArgs...).Run()
	if err != nil {
		return err
	}
	return nil
}

func (e *KappEngine) DeleteGroup(group string, namespace string) error {

	inputArgs := pkg.ExplodeArg(append([]string{"-y", "app-group", "delete", "-n", namespace, "-g", group}, e.Opts.ExtraArgs...))

	err := sh.Command("kapp", inputArgs...).Run()
	if err != nil {
		return err
	}
	return nil
}

func (e *KappEngine) DeployApp(app string, appDir string, namespace string) error {

	inputArgs := pkg.ExplodeArg(append([]string{"-y", "deploy", "-a", app, "-f", appDir, "-n", namespace}, e.Opts.ExtraArgs...))

	err := sh.Command("kapp", inputArgs...).Run()
	if err != nil {
		return err
	}
	return nil
}

func (e *KappEngine) DeployGroup(group string, appGroupDir string, namespace string) error {

	inputArgs := pkg.ExplodeArg(append([]string{"-y", "app-group", "deploy", "-d", appGroupDir, "-n", namespace, "-g", group}, e.Opts.ExtraArgs...))

	err := sh.Command("kapp", inputArgs...).Run()
	if err != nil {
		return err
	}
	return nil
}

package conf

import (
	"os"

	"github.com/bitwormhole/starter-configen/tools/configen2"
	"github.com/bitwormhole/starter/application"
	"github.com/bitwormhole/starter/markup"
)

type theConfigenBootstrap struct {
	markup.Component
	instance   *ConfigenBootstrap  `initMethod:"Start"`
	AppContext application.Context `inject:"context"`
}

type ConfigenBootstrap struct {
	AppContext application.Context
}

func (inst *ConfigenBootstrap) Start() error {
	app := inst.AppContext
	args := app.GetArguments().Export()

	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	app.GetEnvironment().SetEnv("PWD", dir)
	os.Setenv("PWD", dir)

	return configen2.Run(app, args)
}

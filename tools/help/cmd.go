package help

import (
	"github.com/bitwormhole/starter-configen/cmd"
	"github.com/bitwormhole/starter/application"
)

type CommandHandler struct {
}

func (inst *CommandHandler) _impl_() cmd.CommandHandler {
	return inst
}

func (inst *CommandHandler) Init(context application.Context) error {
	return nil
}

func (inst *CommandHandler) Execute(cmd string, args []string) error {
	if cmd == "about" {
		return RunAbout(args)
	} else if cmd == "version" {
		return RunVersion(args)
	}
	return RunHelpDetail(args)
}

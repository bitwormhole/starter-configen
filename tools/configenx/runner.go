package configenx

import (
	"errors"
	"os"

	"github.com/bitwormhole/starter-configen/cmd"
	"github.com/bitwormhole/starter-configen/tools/configen1"
	"github.com/bitwormhole/starter-configen/tools/configen2"
	"github.com/bitwormhole/starter/application"
	"github.com/bitwormhole/starter/collection"
	"github.com/bitwormhole/starter/io/fs"
)

func Run(ctx application.Context, args []string) error {

	pwd, ok := os.LookupEnv("PWD")
	if !ok {
		return errors.New("no env: PWD")
	}

	dir := fs.Default().GetPath(pwd)
	propsFile := dir.GetChild("configen.properties")
	text, err := propsFile.GetIO().ReadText()
	if err != nil {
		return err
	}

	props, err := collection.ParseProperties(text, nil)
	if err != nil {
		return err
	}

	ver, err := props.GetPropertyRequired("configen.version")
	if err != nil {
		return err
	}

	if ver == "1" {
		return configen1.Run(ctx, args)
	} else if ver == "2" {
		return configen2.Run(ctx, args)
	}

	return errors.New("unsupported configen version:" + ver)
}

////////////////////////////////////////////////////////////////////////////////

type CommandHandler struct {
	context application.Context
}

func (inst *CommandHandler) _impl_() cmd.CommandHandler {
	return inst
}

func (inst *CommandHandler) Init(context application.Context) error {
	inst.context = context
	return nil
}

func (inst *CommandHandler) Execute(cmd string, args []string) error {
	return Run(inst.context, args)
}

////////////////////////////////////////////////////////////////////////////////

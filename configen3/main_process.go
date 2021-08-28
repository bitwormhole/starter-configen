package configen3

import (
	"github.com/bitwormhole/starter-configen/configen3/configen"
)

type MainProcess struct {
}

func (inst *MainProcess) _Impl() configen.Process {
	return inst
}

func (inst *MainProcess) Run(ctx *configen.Context) error {

	err := inst.doFindNodes(ctx)
	if err != nil {
		return err
	}

	err = inst.doClean(ctx)
	if err != nil {
		return err
	}

	err = inst.doScan(ctx)
	if err != nil {
		return err
	}

	err = inst.doGenerate1(ctx)
	if err != nil {
		return err
	}

	err = inst.doGenerate2(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (inst *MainProcess) doFindNodes(ctx *configen.Context) error {
	p := &FindNodeProcess{}
	return p.Run(ctx)
}

func (inst *MainProcess) doClean(ctx *configen.Context) error {
	p := &CleanProcess{}
	return p.Run(ctx)
}

func (inst *MainProcess) doScan(ctx *configen.Context) error {
	p := &ScanProcess{}
	return p.Run(ctx)
}

func (inst *MainProcess) doGenerate1(ctx *configen.Context) error {
	return nil
}

func (inst *MainProcess) doGenerate2(ctx *configen.Context) error {
	p := &gen2Process{}
	return p.Run(ctx)
}

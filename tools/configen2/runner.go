package configen2

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/bitwormhole/starter/application"
	"github.com/bitwormhole/starter/io/fs"
)

func Run(ac application.Context, args []string) error {
	r := &runner{}
	return r.Run(ac)
}

type runner struct {
}

func (inst *runner) Run(ac application.Context) error {

	context := &Context{AppContext: ac}

	err := inst.doInit(context)
	if err != nil {
		return err
	}

	err = inst.doScan(context)
	if err != nil {
		return err
	}

	err = inst.doBuildDom2(context)
	if err != nil {
		return err
	}

	err = inst.doBuildCode(context)
	if err != nil {
		return err
	}

	err = inst.doGolangCodePostProcess(context)
	if err != nil {
		return err
	}

	err = inst.doSaveCode(context)
	if err != nil {
		return err
	}

	return nil
}

func (inst *runner) doInit(ctx *Context) error {

	pwd, ok := os.LookupEnv("PWD")
	if !ok {
		return errors.New("no env: PWD")
	}

	dir := fs.Default().GetPath(pwd)

	codeBuilder := &golangCodeBuilder{}
	err := codeBuilder.init(ctx)
	if err != nil {
		return err
	}

	ctx.PWD = dir
	ctx.InputFile = dir.GetChild("configen.properties")
	ctx.OutputFile = dir.GetChild("auto_generated_by_starter_configen.go")

	ctx.ComTable = (&defaultComInfoTable{}).init(ctx)
	ctx.DirectoryScanner = (&defaultDirectoryScanner{}).init(ctx)
	ctx.SourceFileScanner = (&defaultSourceFileScanner{}).init(ctx)
	ctx.Imports = (&defaultImportManager{}).init(ctx)
	ctx.CodeBuilder = codeBuilder
	ctx.Dom2Builder = (&defaultDom2Builder{}).init(ctx)

	ctx.DOM1.Documents = make(map[string]*Dom1doc)
	ctx.DOM1.Imports = make(map[string]bool)
	ctx.DOM1.PackageName = ""

	ctx.DOM2.Components = make(map[string]*Dom2component)
	ctx.DOM2.Imports = make(map[string]string)

	return nil
}

func (inst *runner) doScan(ctx *Context) error {
	dir := ctx.PWD
	return ctx.DirectoryScanner.Scan(dir)
}

func (inst *runner) doBuildDom2(ctx *Context) error {
	err := ctx.Dom2Builder.Build(ctx)
	if err != nil {
		return err
	}

	// print dom2 json
	dom2 := &ctx.DOM2
	data, err := json.MarshalIndent(dom2, "", "\t")
	if err == nil {
		fmt.Println("dom2:", string(data))
	}

	return nil
}

func (inst *runner) doGolangCodePostProcess(ctx *Context) error {
	gcpp := &golangCodePostProcessor{}
	gcpp.Init(ctx)
	return gcpp.ReImportPackages()
}

func (inst *runner) doBuildCode(ctx *Context) error {
	dom := &ctx.DOM2
	code, err := ctx.CodeBuilder.Build(dom)
	ctx.OutputSourceCode = code
	return err
}

func (inst *runner) doSaveCode(ctx *Context) error {

	file := ctx.OutputFile
	code := ctx.OutputSourceCode
	fmt.Println("Output Source Code:", code)

	if ctx.DryRun {
		return nil
	}

	fmt.Println("Save code to file:", file.Path())
	return file.GetIO().WriteText(code, nil)
}

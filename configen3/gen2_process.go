package configen3

import (
	"strings"

	"github.com/bitwormhole/starter-configen/configen3/configen"
	"github.com/bitwormhole/starter-configen/configen3/data/model"
	"github.com/bitwormhole/starter-configen/configen3/generator/gen2"
	"github.com/bitwormhole/starter-configen/configen3/generator/helper"
	"github.com/bitwormhole/starter/io/fs"
	"github.com/bitwormhole/starter/vlog"
)

type gen2Process struct {
	// depthLimit int
}

func (inst *gen2Process) _Impl() configen.Process {
	return inst
}

func (inst *gen2Process) Run(ctx *configen.Context) error {
	children := ctx.Store.RootNode.Children
	for _, child := range children {
		err := inst.handleChild(child, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (inst *gen2Process) handleChild(child *model.Child, ctx *configen.Context) error {

	dir := child.Path
	if dir.IsFile() {
		dir = dir.Parent()
	}

	thisPkgPath, err := (&helper.ThisPackagePathComputer{}).Compute(ctx, dir)
	if err != nil {
		return err
	}

	sc := &scanningContext{}
	sc.context = ctx
	sc.thisPackagePath = thisPkgPath
	sc.mixdown.Package.Alias = dir.Name()

	err = inst.scanSourceFilesInDir(dir, sc)
	if err != nil {
		return err
	}

	code, err := inst.buildCode(child, sc)
	if err != nil {
		return err
	}

	err = inst.saveCode(code, child)
	if err != nil {
		return err
	}

	return nil
}

func (inst *gen2Process) buildCode(child *model.Child, sc *scanningContext) (string, error) {
	source := &sc.mixdown
	res := sc.context.Resources
	templ := gen2.NewTemplate(res, "configen3/gen2/templates")
	return templ.Build(source)
}

func (inst *gen2Process) saveCode(code string, child *model.Child) error {
	const outputFileName = "auto_generated_config_by_starter_configen.go"
	dir := child.Path
	file := dir
	if file.IsFile() {
		dir = file.Parent()
		if !strings.HasSuffix(file.Name(), ".go") {
			file = dir.GetChild(outputFileName)
		}
	} else if dir.IsDir() {
		file = dir.GetChild(outputFileName)
	} else {
		file = dir.GetChild(outputFileName)
	}
	vlog.Info("Write generated configuration to file ", file.Path())
	return file.GetIO().WriteText(code, nil, false)
}

func (inst *gen2Process) scanSourceFilesInDir(dir fs.Path, ctx *scanningContext) error {
	list := dir.ListItems()
	for _, file := range list {
		if !file.IsFile() {
			continue
		}
		name := file.Name()
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		err := inst.loadSourceFile(file, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (inst *gen2Process) loadSourceFile(file fs.Path, ctx *scanningContext) error {
	loader := &GoSourceLoader{}
	src, err := loader.Load(file, ctx.thisPackagePath)
	if err != nil {
		return err
	}
	return inst.handleSource(src, ctx)
}

func (inst *gen2Process) handleSource(source *model.GoSource, ctx *scanningContext) error {

	src := source.StructList
	dst := ctx.mixdown.StructList

	for _, item := range src {
		if item == nil {
			continue
		}
		if !item.IsComponentDefine() {
			continue
		}
		dst = append(dst, item)
	}

	ctx.mixdown.StructList = dst
	return nil
}

////////////////////////////////////////////////////////////////////////////////

type scanningContext struct {
	context         *configen.Context
	thisPackagePath string
	mixdown         model.GoSource
}

package configen3

import (
	"errors"
	"fmt"
	"os"

	"github.com/bitwormhole/starter-configen/configen3/configen"
	"github.com/bitwormhole/starter-configen/configen3/data/model"
	"github.com/bitwormhole/starter/io/fs"
	"github.com/bitwormhole/starter/vlog"
)

type FindNodeProcess struct {
	configenPropertiesFileName string
}

func (inst *FindNodeProcess) _Impl() configen.Process {
	return inst
}

func (inst *FindNodeProcess) Run(ctx *configen.Context) error {

	inst.configenPropertiesFileName = "configen.properties"

	pwd, err := inst.getPWD(ctx)
	if err != nil {
		return err
	}
	vlog.Debug("PWD=", pwd)

	goMod, err := inst.findGoModFile(ctx, pwd)
	if err != nil {
		return err
	}
	vlog.Debug("go.mod=", goMod)

	root, err := inst.getRootNode(ctx, goMod)
	if err != nil {
		return err
	}
	vlog.Debug("node.root=", root)

	err = inst.loadNodeTree(ctx, root)
	if err != nil {
		return err
	}

	return nil
}

func (inst *FindNodeProcess) loadNodeTree(ctx *configen.Context, root fs.Path) error {
	err := inst.loadNodeAndChild(ctx, root, nil, 0)
	if err != nil {
		return err
	}
	nodelist := ctx.Store.Nodes.List
	if nodelist != nil {
		if len(nodelist) > 0 {
			root := nodelist[0]
			ctx.Store.RootNode = root
		}
	}
	return err
}

func (inst *FindNodeProcess) loadNodeAndChild(ctx *configen.Context, file fs.Path, parent *model.Node, depth int) error {

	if depth > 9 {
		return errors.New("the path is too deep, path=" + file.Path())
	}

	filename := inst.configenPropertiesFileName
	dir := file.Parent()

	if file.IsFile() {
		// NOP
	} else if file.IsDir() {
		dir = file
		file = dir.GetChild(filename)
	} else if !file.Exists() && dir.IsDir() {
		file = dir.GetChild(filename)
	} else {
		return errors.New("bad file path: " + file.Path())
	}

	if !file.IsFile() {
		return errors.New("the path is not a file, path=" + file.Path())
	}

	loader := &nodeLoader{}
	node, err := loader.Load(file)
	if err != nil {
		return err
	}

	node.Parent = parent
	all := &ctx.Store.Nodes
	all.List = append(all.List, node)

	list := node.Children
	if list == nil {
		return nil
	}

	for _, child := range list {
		child.Path = dir.GetChild(child.Href)
		err := inst.loadNodeAndChild(ctx, child.Path, node, depth+1)
		if err != nil {
			return err
		}
	}

	return nil
}

func (inst *FindNodeProcess) getRootNode(ctx *configen.Context, goModFile fs.Path) (fs.Path, error) {
	name := inst.configenPropertiesFileName
	dir := goModFile.Parent()
	file := dir.GetChild(name)
	if file.IsFile() {
		return file, nil
	}
	str := fmt.Sprint("no root node file [" + name + "] in path [" + dir.Path() + "]")
	return nil, errors.New(str)
}

func (inst *FindNodeProcess) findGoModFile(ctx *configen.Context, pwd fs.Path) (fs.Path, error) {
	const name = "go.mod"
	dir := pwd
	for ; dir != nil; dir = dir.Parent() {
		file := dir.GetChild(name)
		if file.IsFile() {
			return file, nil
		}
	}
	return nil, errors.New("cannot find file [go.mod] in path [" + pwd.Path() + "]")
}

func (inst *FindNodeProcess) getPWD(ctx *configen.Context) (fs.Path, error) {
	pwd := ctx.Store.PWD
	if pwd == nil {
		dir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		pwd = fs.Default().GetPath(dir)
		ctx.Store.PWD = pwd
	}
	return pwd, nil
}

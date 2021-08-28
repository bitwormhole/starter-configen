package configen3

import (
	"errors"
	"strings"

	"github.com/bitwormhole/starter-configen/configen3/configen"
	"github.com/bitwormhole/starter-configen/configen3/data/model"
	"github.com/bitwormhole/starter/io/fs"
	"github.com/bitwormhole/starter/vlog"
)

// CleanProcess 清理过程
type CleanProcess struct {
	depthLimit int
}

func (inst *CleanProcess) _Impl() configen.Process {
	return inst
}

// Run 运行
func (inst *CleanProcess) Run(ctx *configen.Context) error {
	list := ctx.Store.Nodes.List
	for _, node := range list {
		err := inst.cleanNode(node)
		if err != nil {
			return err
		}
	}
	return nil
}

func (inst *CleanProcess) cleanNode(node *model.Node) error {
	dir := node.Path
	if dir.IsFile() {
		dir = dir.Parent()
	} else if !dir.Exists() {
		dir = dir.Parent()
	}
	return inst.cleanDir(dir)
}

func (inst *CleanProcess) cleanDir(dir fs.Path) error {
	if !dir.IsDir() {
		return errors.New("the dir is not exists, path=" + dir.Path())
	}
	vlog.Info("clean dir " + dir.Path())
	items := dir.ListItems()
	for _, item := range items {
		if inst.acceptToDelete(item) {
			vlog.Debug("rm " + item.Name())
			err := item.Delete()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (inst *CleanProcess) acceptToDelete(item fs.Path) bool {
	const prefix = "generated_"
	const suffix = ".go"
	name := item.Name()
	return strings.HasPrefix(name, prefix) && strings.HasSuffix(name, suffix)
}

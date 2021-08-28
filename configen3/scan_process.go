package configen3

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/bitwormhole/starter-configen/configen3/configen"
	"github.com/bitwormhole/starter-configen/configen3/data/model"
	"github.com/bitwormhole/starter-configen/configen3/generator/helper"
	"github.com/bitwormhole/starter-configen/configen3/scanning"
	"github.com/bitwormhole/starter/io/fs"
	"github.com/bitwormhole/starter/vlog"
)

// ScanProcess 扫描过程
type ScanProcess struct {
	depthLimit int
}

func (inst *ScanProcess) _Impl() configen.Process {
	return inst
}

// Run 运行
func (inst *ScanProcess) Run(ctx *configen.Context) error {

	inst.depthLimit = 30

	nodes := ctx.Store.Nodes.List
	for _, node := range nodes {
		err := inst.scanNode(ctx, node)
		if err != nil {
			return err
		}
	}
	return nil
}

func (inst *ScanProcess) scanNode(ctx *configen.Context, node *model.Node) error {

	targets := node.ScanTargets
	for _, tar := range targets {
		err := inst.scanTarget(ctx, tar)
		if err != nil {
			return err
		}
	}
	return nil
}

func (inst *ScanProcess) scanTarget(ctx *configen.Context, target *model.Scan) error {

	vlog.Debug("scan dir, path=", target.Path)

	scanner := &scanning.Scanner{}
	working := scanner.Init(ctx, target)
	ctx.Store.CurrentScanning = working

	err := inst.scanDir(ctx, target.Path, target.R, 0)
	if err != nil {
		return err
	}

	err = working.WriteResultTo(target.WriteTo)
	if err != nil {
		return err
	}

	return nil
}

func (inst *ScanProcess) scanDir(ctx *configen.Context, dir fs.Path, r bool, depth int) error {

	if depth > inst.depthLimit {
		return errors.New("the path is too deep, path=" + dir.Path())
	}

	if !dir.IsDir() {
		return errors.New("the path is not a dir, path=" + dir.Path())
	}

	vlog.Trace("scan dir: ", dir.Path())

	items := dir.ListItems()
	for _, item := range items {
		if item.IsDir() {
			if r {
				err := inst.scanDir(ctx, item, r, depth+1)
				if err != nil {
					return err
				}
			}
		} else if item.IsFile() {
			err := inst.scanFile(ctx, item)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (inst *ScanProcess) acceptSourceFile(ctx *configen.Context, file fs.Path) bool {
	dir := file.Parent()
	configenfile := dir.GetChild("configen.properties")
	return !configenfile.Exists()
}

func (inst *ScanProcess) scanFile(ctx *configen.Context, file fs.Path) error {

	const srcFileNameSuffix = ".go"
	filename := file.Name()
	if !strings.HasSuffix(filename, srcFileNameSuffix) {
		return nil // skip
	}
	if !inst.acceptSourceFile(ctx, file) {
		return nil // skip
	}

	vlog.Trace("scan source: ", file.Name())

	thisPkgPath, err := (&helper.ThisPackagePathComputer{}).Compute(ctx, file)
	if err != nil {
		return err
	}

	loader := &GoSourceLoader{}
	src, err := loader.Load(file, thisPkgPath)
	if err != nil {
		return err
	}

	inst.handleSource(ctx, src)

	return nil
}

func (inst *ScanProcess) handleSource(ctx *configen.Context, src *model.GoSource) error {

	if vlog.Default().IsTraceEnabled() {
		data, err := json.MarshalIndent(src, "", "\t")
		if err != nil {
			return err
		}
		vlog.Trace("source.json:", string(data))
	}
	ctx.Store.CurrentScanning.AddSource(src)
	return nil
}

////////////////////////////////////////////////////////////////////////////////

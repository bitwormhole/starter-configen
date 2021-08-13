package configen2

import (
	"crypto/sha1"
	"errors"
	"strings"

	"github.com/bitwormhole/starter/util"
)

type ImportManager interface {
	AddImport(path string)
	AddImportWithAlias(alias string, path string)
	AddImportWithoutHash(path string)
	FindAliasByPath(path string) (string, error)
}

////////////////////////////////////////////////////////////////////////////////

type defaultImportManager struct {
	context       *Context
	markupNoAlias string
}

func (inst *defaultImportManager) init(ctx *Context) ImportManager {
	inst.context = ctx
	inst.markupNoAlias = "_"
	return inst
}

func (inst *defaultImportManager) FindAliasByPath(path string) (string, error) {
	table := inst.context.DOM2.Imports
	alias := table[path]
	if alias == "" {
		return "", errors.New("no import from path:" + path)
	}
	return alias, nil
}

func (inst *defaultImportManager) AddImport(path string) {
	alias := inst.makeAliasByPath(path, true)
	inst.innerAddImport(alias, path)
}

func (inst *defaultImportManager) AddImportWithAlias(alias string, path string) {
	if alias == "" {
		alias = inst.makeAliasByPath(path, true)
	}
	inst.innerAddImport(alias, path)
}

func (inst *defaultImportManager) AddImportWithoutHash(path string) {
	alias := inst.makeAliasByPath(path, false)
	inst.innerAddImport(alias, path)
}

func (inst *defaultImportManager) innerAddImport(alias string, path string) {
	table := inst.context.DOM2.Imports
	older := table[path]
	if older != "" {
		// skip
		return
	}
	table[path] = alias
}

func (inst *defaultImportManager) makeAliasByPath(path string, hash bool) string {

	// name
	name := path
	index := strings.LastIndexByte(path, '/')
	if index > 0 {
		name = path[index+1:]
	}

	if !hash {
		return name
	}

	// sum
	sha1sum := sha1.Sum([]byte(path))
	sum := util.StringifyBytes(sha1sum[:])

	// mix
	return name + "_" + sum[0:6]
}

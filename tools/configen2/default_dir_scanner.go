package configen2

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/bitwormhole/starter/io/fs"
)

type defaultDirectoryScanner struct {
	Context *Context
}

func (inst *defaultDirectoryScanner) init(ctx *Context) DirectoryScanner {
	inst.Context = ctx
	return inst
}

func (inst *defaultDirectoryScanner) Scan(dir fs.Path) error {

	if !dir.IsDir() {
		return errors.New("the path is not a dir: " + dir.Path())
	}

	fmt.Println("Scan directory", dir.Path())
	list := dir.ListNames()
	sort.Strings(list)
	fileScanner := inst.Context.SourceFileScanner

	for index := range list {
		name := list[index]
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		file := dir.GetChild(name)
		doc, err := fileScanner.Scan(file)
		if err != nil {
			return err
		}
		inst.handleResultDoc(doc, file)
	}

	return nil
}

func (inst *defaultDirectoryScanner) handleResultDoc(doc *Dom1doc, path fs.Path) error {

	filename := path.Name()
	inst.Context.DOM1.PackageName = doc.PackageName
	inst.Context.DOM1.Documents[filename] = doc
	imports := doc.ImportList

	if imports != nil {
		for index := range imports {
			imp := imports[index]
			imp_path := imp.Path
			imp_alias := imp.Alias
			if imp_alias == "" {
				imp_alias = inst.getAliasByImportPath(imp_path)
				imp.Alias = imp_alias
			}
			inst.Context.DOM1.Imports[imp_path] = true
		}
	}

	/////// print json
	data, err := json.MarshalIndent(doc, "", "    ")
	if err == nil {
		fmt.Println(filename, ":", string(data))
	}

	return nil
}

func (inst *defaultDirectoryScanner) getAliasByImportPath(path string) string {
	index := strings.LastIndexByte(path, '/')
	if index < 0 {
		return path
	}
	return path[index+1:]
}

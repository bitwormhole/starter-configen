package configen1

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/bitwormhole/starter/io/fs"
)

type configenRunner struct {
	context *configenContext
}

func (inst *configenRunner) run() error {

	err := inst.init()
	if err != nil {
		return err
	}

	err = inst.loadConfigFile()
	if err != nil {
		return err
	}

	err = inst.scanPWD()
	if err != nil {
		return err
	}

	err = inst.scanSourceFiles()
	if err != nil {
		return err
	}

	err = inst.makeDOM()
	if err != nil {
		return err
	}

	err = inst.buildCode()
	if err != nil {
		return err
	}

	err = inst.output()
	if err != nil {
		return err
	}

	return nil
}

func (inst *configenRunner) init() error {
	inst.context.OutputFileName = "auto_generated_by_starter_configen.go"
	inst.context.ConfigFileName = "configen.properties"
	return nil
}

func (inst *configenRunner) loadConfigFile() error {

	context := inst.context
	filename := context.ConfigFileName
	pwd := context.PWD
	file := pwd.GetChild(filename)

	data, err := file.GetIO().ReadText()
	if err == nil {
		strings.TrimSpace(data) // NOP
	}
	return err
}

func (inst *configenRunner) scanPWD() error {

	dir := inst.context.PWD
	items := dir.ListItems()
	sources := []fs.Path{}
	fmt.Println("scan dir " + dir.Path())

	for index := range items {
		item := items[index]
		if item == nil {
			continue
		}
		if !item.IsFile() {
			continue
		}
		name := item.Name()
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		sources = append(sources, item)
		fmt.Println("find go source file " + item.Path())
	}

	inst.context.GoSourceFiles = sources
	return nil
}

func (inst *configenRunner) scanSourceFiles() error {
	files := inst.context.GoSourceFiles
	for index := range files {
		file := files[index]
		err := inst.scanSourceFile(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (inst *configenRunner) scanSourceFile(file fs.Path) error {

	scanner := &golangSourceFileScanner{}
	result, err := scanner.scan(file)
	if err != nil {
		return err
	}

	context := inst.context
	src := result.injectionList
	dst := context.DomInjectionTable

	if dst == nil {
		dst = make(map[string]*DomInjection)
		context.DomInjectionTable = dst
	}

	if src != nil {
		for index := range src {
			item := src[index]
			dst[item.Name] = item
		}
	}

	inst.context.PackageName = result.packageName
	return nil
}

func (inst *configenRunner) makeDOM() error {

	src := inst.context.DomInjectionTable
	dst := inst.context.ComConfigList
	keys := make([]string, 0)

	if src == nil {
		return errors.New("no src")
	}

	if dst == nil {
		dst = make([]*ComConfigInfo, 0)
	}

	for key := range src {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for index := range keys {
		key := keys[index]
		item1 := src[key]
		item2 := &ComConfigInfo{}
		err := item2.init(item1)
		if err != nil {
			// return err
			fmt.Println("skip func ", key, "()")
			continue
		}
		dst = append(dst, item2)
	}

	inst.context.ComConfigList = dst
	return nil
}

func (inst *configenRunner) buildCode() error {

	list := inst.context.ComConfigList
	js, err := json.MarshalIndent(list, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println("all_components:", string(js))

	builder := NewSourceCodeBuilder()
	builder.Begin(inst.context.PackageName)
	for index := range list {
		item := list[index]
		builder.AppendComponent(item)
	}
	builder.End()

	code, err := builder.Create()
	if err != nil {
		return err
	}

	inst.context.OutputCode = code
	return nil
}

func (inst *configenRunner) output() error {

	context := inst.context

	filename := context.OutputFileName
	dir := context.PWD
	file := dir.GetChild(filename)
	code := context.OutputCode

	fmt.Println("output.code: ", code)
	fmt.Println("write to file ", file.Path())
	return file.GetIO().WriteText(code, nil)
}

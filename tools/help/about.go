package help

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

////////////////////////////////////////////////////////////////////////////////
// AboutInfo

type AboutInfo struct {
	FactoryName     string
	ApplicationName string
	Version         string
	Web             string
}

func (inst *AboutInfo) init() {
	inst.FactoryName = "bitwormhole.com"
	inst.ApplicationName = "Bitwormhole Starter-Tools"
	inst.Version = "0.0.1_20210516_1957"
	inst.Web = "https://help.bitwormhole.com/starter/tools"
}

func (inst *AboutInfo) toProperties() map[string]string {
	table := make(map[string]string)
	table["application.name"] = inst.ApplicationName
	table["application.version"] = inst.Version
	table["factory.name"] = inst.FactoryName
	table["factory.web"] = inst.Web
	return table
}

func (inst *AboutInfo) toJSON() string {
	text, err := json.MarshalIndent(inst, "", "\t")
	if err != nil {
		return err.Error()
	}
	return string(text)
}

////////////////////////////////////////////////////////////////////////////////
// AboutRunner

type AboutRunner struct {
}

func (inst *AboutRunner) logList(list []string, title string) {

	fmt.Println(title)

	if list == nil {
		return
	}

	for index := range list {
		item := list[index]
		fmt.Println("\t", item)
	}
}

func (inst *AboutRunner) logTable(table map[string]string, title string) {

	fmt.Println(title)

	if table == nil {
		return
	}

	keys := make([]string, 0)
	for key := range table {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for index := range keys {
		key := keys[index]
		value := table[key]
		fmt.Println("\t", key, "=", value)
	}
}

func (inst *AboutRunner) run(args []string) error {

	environ := os.Environ()

	about := &AboutInfo{}
	about.init()

	inst.logList(os.Args, "CLI.Arguments")
	inst.logList(args, "Runner.Arguments")
	inst.logList(environ, "Environment")
	inst.logTable(about.toProperties(), "About")

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func RunAbout(args []string) error {
	runner := &AboutRunner{}
	return runner.run(args)
}

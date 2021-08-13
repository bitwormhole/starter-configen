package templates

import (
	"sort"
	"strings"

	"github.com/bitwormhole/starter-configen/tools/configen2"
)

const NL = "\n"
const TAB = "\t"
const SPACE = " "
const QUOTE = "\""

type MainTemplate struct {
	baseGroupTemplate

	templateForCominfobuilderItem configen2.CodeTemplate
}

func (inst *MainTemplate) _impl() innerTemplate {
	return inst
}

func (inst *MainTemplate) InitWithGroup(group *innerTemplateGroup) error {

	err := inst.InitTemplate(group, "configen2/templates/main.template")
	if err != nil {
		return err
	}

	inst.AddComputed("package.name", func(ctx *configen2.BuildingContext) (string, error) { return inst.getPackageName(ctx) })
	inst.AddComputed("import.items", func(ctx *configen2.BuildingContext) (string, error) { return inst.getImportItems(ctx) })
	inst.AddComputed("cominfobuilder.items", func(ctx *configen2.BuildingContext) (string, error) { return inst.getComInfoBuilderItems(ctx) })
	inst.AddComputed("component.methods", func(ctx *configen2.BuildingContext) (string, error) { return inst.getComMethods(ctx) })

	return nil
}

func (inst *MainTemplate) Build(ctx *configen2.BuildingContext) (string, error) {
	return inst.BuildFromTemplate(ctx)
}

func (inst *MainTemplate) getPackageName(ctx *configen2.BuildingContext) (string, error) {
	return ctx.DOM.PackageName, nil
}

func (inst *MainTemplate) getImportItems(ctx *configen2.BuildingContext) (string, error) {

	builder := &strings.Builder{}
	table := ctx.DOM.Imports

	for path := range table {
		alias := table[path]
		builder.WriteString(TAB)
		builder.WriteString(alias)
		builder.WriteString(SPACE)
		builder.WriteString(QUOTE)
		builder.WriteString(path)
		builder.WriteString(QUOTE)
		builder.WriteString(NL)
	}

	return builder.String(), nil
}

func (inst *MainTemplate) getComInfoBuilderItems(ctx *configen2.BuildingContext) (string, error) {

	builder := &strings.Builder{}
	template := inst.group.comInfoTemplate
	comTab := ctx.DOM.Components
	keys := []string{}

	for key := range comTab {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for index := range keys {
		key := keys[index]
		com := comTab[key]
		child := ctx.NewChild()
		child.Component = com
		code, err := template.Build(child)
		if err != nil {
			return "", err
		}
		builder.WriteString(code)
	}

	return builder.String(), nil
}

func (inst *MainTemplate) getComMethods(ctx *configen2.BuildingContext) (string, error) {

	builder := &strings.Builder{}
	template := inst.group.injectionAdapterTemplate
	comTab := ctx.DOM.Components

	keys := []string{}
	for key := range comTab {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for index := range keys {
		key := keys[index]
		com := comTab[key]
		child := ctx.NewChild()
		child.Component = com
		code, err := template.Build(child)
		if err != nil {
			return "", err
		}
		builder.WriteString(code)
	}

	return builder.String(), nil
}

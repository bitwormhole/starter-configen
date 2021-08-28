package gen2

import (
	"strconv"
	"strings"
	"text/template"

	"github.com/bitwormhole/starter-configen/configen3/data/model"
	"github.com/bitwormhole/starter-configen/configen3/mark"
)

type mainTemplate interface {
	// CreateTemplate(ctx *gen2Context) (string, error)
	Template
}

type mainTemplateData struct {
	PackageShortName        string
	ImportListItems         string
	ComItemsInConfigFunc    string
	ComItemsAfterConfigFunc string
}

type mainTemplateImpl struct {
	context *gen2Context
	templ   *template.Template
	data    mainTemplateData
}

func (inst *mainTemplateImpl) init(ctx *gen2Context) mainTemplate {

	fm := make(template.FuncMap)
	// fm[ "" ] = func ()  { 	}

	inst.context = ctx
	inst.templ = ctx.loadTemplate("main.template", fm)
	return inst
}

func (inst *mainTemplateImpl) buildImportListInner() string {
	builder := strings.Builder{}
	items := inst.context.importsBuilder.Create()
	if items == nil {
		return ""
	}
	for _, item := range items {
		if item.Path == "" {
			continue
		}
		builder.WriteString(mark.TAB)
		builder.WriteString(item.Alias)
		builder.WriteString(mark.SPACE)
		builder.WriteString(mark.QUOTE2)
		builder.WriteString(item.Path)
		builder.WriteString(mark.QUOTE2)
		builder.WriteString(mark.NL)
	}
	return builder.String()
}

func (inst *mainTemplateImpl) checkComInfo(info *ComInfo, index int) error {
	id := info.ComID
	if id == "" {
		id = "com" + strconv.Itoa(index) + "-" + info.InstanceStructName
		info.ComID = id
	}
	return nil
}

func (inst *mainTemplateImpl) buildComItems(source *model.GoSource) (string, string, error) {

	ctx := inst.context
	tmp1 := ctx.templateFactory.createComItemInConfigFuncTemplate(ctx)
	tmp2 := ctx.templateFactory.createComItemAfterConfigFuncTemplate(ctx)
	builder1 := &strings.Builder{}
	builder2 := &strings.Builder{}
	list := source.StructList

	for index, item := range list {
		info := &ComInfo{}
		err := info.init(item)
		if err != nil {
			return "", "", err
		}

		err = inst.checkComInfo(info, index)
		if err != nil {
			return "", "", err
		}

		str1, err := tmp1.BuildComponent(info)
		if err != nil {
			return "", "", err
		}

		builder1.WriteString(str1)
		builder1.WriteString(mark.NL)
		str2, err := tmp2.BuildComponent(info)
		if err != nil {
			return "", "", err
		}
		builder2.WriteString(str2)
		builder2.WriteString(mark.NL)
	}

	return builder1.String(), builder2.String(), nil
}

func (inst *mainTemplateImpl) Build(source *model.GoSource) (string, error) {

	comItemsInFn, comItemsAfterFn, err := inst.buildComItems(source)
	if err != nil {
		return "", err
	}

	data := &inst.data
	data.PackageShortName = source.Package.Alias
	data.ImportListItems = inst.buildImportListInner()
	data.ComItemsInConfigFunc = comItemsInFn
	data.ComItemsAfterConfigFunc = comItemsAfterFn

	str, err := executeTemplate(inst.templ, data)
	if err != nil {
		return "", err
	}

	str = "// (todo:gen2.template) " + mark.NL + str
	return str, nil
}

func executeTemplate(templ *template.Template, params interface{}) (string, error) {
	templ, err := templ.Clone()
	if err != nil {
		return "", err
	}
	builder := &strings.Builder{}
	err = templ.Execute(builder, params)
	if err != nil {
		return "", err
	}
	return builder.String(), nil
}

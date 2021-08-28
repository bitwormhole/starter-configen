package gen2

import (
	"text/template"
)

type comItemInConfigFuncTemplate interface {
	BuildComponent(info *ComInfo) (string, error)
}

type comItemInConfigFuncTemplateData struct {
	ComID      string
	ComScope   string
	ComClass   string
	ComAliases string

	ComFactoryStructName string
}

type comItemInConfigFuncTemplateImpl struct {
	context *gen2Context
	templ   *template.Template
	// data    mainTemplateData
}

func (inst *comItemInConfigFuncTemplateImpl) init(ctx *gen2Context) comItemInConfigFuncTemplate {

	tmp := ctx.loadTemplate("com-item-in-config-func.template", nil)

	inst.context = ctx
	inst.templ = tmp
	return inst
}

func (inst *comItemInConfigFuncTemplateImpl) BuildComponent(info *ComInfo) (string, error) {

	data := &comItemInConfigFuncTemplateData{}
	data.ComID = info.ComID
	data.ComClass = info.ComClass
	data.ComAliases = info.ComAliases
	data.ComScope = info.ComScope
	data.ComFactoryStructName = info.FactoryStructName

	return executeTemplate(inst.templ, data)
}

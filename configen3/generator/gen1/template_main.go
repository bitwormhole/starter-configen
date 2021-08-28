package gen1

import (
	"strings"
	"text/template"

	"github.com/bitwormhole/starter-configen/configen3/data/model"
)

type mainTemplate interface {
	Template
}

type mainTemplateImpl struct {
	context  *gen1Context
	template *template.Template
}

func (o *mainTemplateImpl) init(ctx *gen1Context) mainTemplate {
	t := ctx.loadTemplate("main.template")
	o.context = ctx
	o.template = t
	return o
}

func (o *mainTemplateImpl) Build(s *model.GoSource) (string, error) {

	ctx := o.context
	factory := ctx.templateFactory
	importListTempl := factory.createImportListTemplate(ctx)
	structListTempl := factory.createStructListTemplate(ctx)

	structListText, err := structListTempl.BuildStructs(s)
	if err != nil {
		return "", err
	}

	importListText, err := importListTempl.BuildImports(s)
	if err != nil {
		return "", err
	}

	builder := &strings.Builder{}
	o.writeComment(s, builder)
	o.writePackageDefine(s, builder)
	builder.WriteString(importListText)
	builder.WriteString(structListText)
	return builder.String(), nil
}

func (o *mainTemplateImpl) writeComment(s *model.GoSource, builder *strings.Builder) {
	const nl = "\n"
	const msg1 = "// 这个配置文件是由 starter-configen 工具自动生成的。"
	const msg2 = "// 任何时候，都不要手工修改这里面的内容！！！"
	builder.WriteString(msg1 + nl)
	builder.WriteString(msg2 + nl)
}

func (o *mainTemplateImpl) writePackageDefine(s *model.GoSource, builder *strings.Builder) {
	const nl = "\n"
	pkgName := s.Package.Alias
	builder.WriteString(nl)
	builder.WriteString("package " + pkgName + nl)
	builder.WriteString(nl)
}

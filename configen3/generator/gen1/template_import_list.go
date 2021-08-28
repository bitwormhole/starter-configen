package gen1

import (
	"strings"
	"text/template"

	"github.com/bitwormhole/starter-configen/configen3/data/model"
	"github.com/bitwormhole/starter-configen/configen3/mark"
)

type importListTemplate interface {
	BuildImports(s *model.GoSource) (string, error)
}

type importListTemplateImpl struct {
	context  *gen1Context
	template *template.Template
}

func (inst *importListTemplateImpl) init(ctx *gen1Context) importListTemplate {
	inst.context = ctx
	return inst
}

func (inst *importListTemplateImpl) BuildImports(s *model.GoSource) (string, error) {

	items := inst.context.importsBuilder.Create()

	if items == nil {
		return "", nil
	}

	builder := strings.Builder{}
	builder.WriteString("import (" + mark.NL)

	for _, item := range items {
		builder.WriteString(mark.TAB + item.Alias + mark.SPACE)
		builder.WriteString(mark.QUOTE2 + item.Path + mark.QUOTE2)
		builder.WriteString(mark.NL)
	}

	builder.WriteString(")" + mark.NL)
	return builder.String(), nil
}

package templates

import (
	"errors"
	"strings"

	"github.com/bitwormhole/starter-configen/tools/configen2"
)

type injectionGetPropertyTemplate struct {
	baseGroupTemplate

	typeToMethodTable map[string]string
}

func (inst *injectionGetPropertyTemplate) _impl() innerTemplate {
	return inst
}

func (inst *injectionGetPropertyTemplate) InitWithGroup(group *innerTemplateGroup) error {

	err := inst.InitTemplate(group, "configen2/templates/injection_get_property.template")
	if err != nil {
		return err
	}

	inst.AddComputed("injection.getter.readmethod", func(ctx *configen2.BuildingContext) (string, error) { return inst.propertyGetterReadMethod(ctx) })
	//	inst.AddComputed("injection.getter.name", func(ctx *configen2.BuildingContext) (string, error) { return inst.propertyInjectionGetterName(ctx) })
	//	inst.AddComputed("injection.field.type", func(ctx *configen2.BuildingContext) (string, error) { return inst.propertyInjectionFieldType(ctx) })

	return nil
}

func (inst *injectionGetPropertyTemplate) Build(ctx *configen2.BuildingContext) (string, error) {

	return inst.BuildFromTemplate(ctx)
}

func (inst *injectionGetPropertyTemplate) getTypeToMethodMap() map[string]string {
	table := inst.typeToMethodTable
	if table == nil {
		table = map[string]string{
			"string":  "ReadString",
			"bool":    "ReadBool",
			"int":     "ReadInt",
			"int32":   "ReadInt32",
			"int64":   "ReadInt64",
			"float32": "ReadFloat32",
			"float64": "ReadFloat64",
		}
		inst.typeToMethodTable = table
	}
	return table
}

func (inst *injectionGetPropertyTemplate) propertyGetterReadMethod(ctx *configen2.BuildingContext) (string, error) {

	field := ctx.Injection
	fieldType := strings.TrimSpace(field.FieldType)
	table := inst.getTypeToMethodMap()
	method := table[fieldType]

	if method != "" {
		return method, nil
	}

	msg := &strings.Builder{}
	msg.WriteString("Unsupported fieldType, ")
	msg.WriteString(" selector:[" + field.Selector + "]")
	msg.WriteString(" struct:[" + ctx.Component.StructName + "]")
	msg.WriteString(" field:[" + field.FieldName + "]")
	msg.WriteString(" fieldType:[" + fieldType + "]")
	return "", errors.New(msg.String())
}

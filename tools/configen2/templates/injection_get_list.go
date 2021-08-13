package templates

import (
	"errors"
	"strings"

	"github.com/bitwormhole/starter-configen/tools/configen2"
)

type injectionGetListTemplate struct {
	baseGroupTemplate
}

func (inst *injectionGetListTemplate) _impl() innerTemplate {
	return inst
}

func (inst *injectionGetListTemplate) InitWithGroup(group *innerTemplateGroup) error {

	err := inst.InitTemplate(group, "configen2/templates/injection_get_list.template")
	if err != nil {
		return err
	}

	inst.AddComputed("list.item.type", func(ctx *configen2.BuildingContext) (string, error) { return inst.propertyListItemType(ctx) })
	//	inst.AddComputed("injection.getter.name", func(ctx *configen2.BuildingContext) (string, error) { return inst.propertyInjectionGetterName(ctx) })
	//	inst.AddComputed("injection.field.type", func(ctx *configen2.BuildingContext) (string, error) { return inst.propertyInjectionFieldType(ctx) })

	return nil
}

func (inst *injectionGetListTemplate) Build(ctx *configen2.BuildingContext) (string, error) {

	return inst.BuildFromTemplate(ctx)
}

func (inst *injectionGetListTemplate) propertyListItemType(ctx *configen2.BuildingContext) (string, error) {
	const prefix = "[]"
	field := ctx.Injection
	fieldType := strings.TrimSpace(field.FieldType)
	if !strings.HasPrefix(fieldType, prefix) {
		err := errors.New("the field is not a list, field: " + field.FieldName)
		return "", err
	}
	offset := len(prefix)
	return fieldType[offset:], nil
}

package templates

import (
	"errors"
	"strings"

	"github.com/bitwormhole/starter-configen/tools/configen2"
)

type injectionGetterTemplate struct {
	baseGroupTemplate
}

func (inst *injectionGetterTemplate) _impl() innerTemplate {
	return inst
}

func (inst *injectionGetterTemplate) InitWithGroup(group *innerTemplateGroup) error {

	err := inst.InitTemplate(group, "configen2/templates/injection_getter.template")
	if err != nil {
		return err
	}

	inst.AddComputed("component.adapter.type", func(ctx *configen2.BuildingContext) (string, error) { return inst.propertyComAdapterType(ctx) })
	inst.AddComputed("injection.field.type", func(ctx *configen2.BuildingContext) (string, error) { return inst.propertyInjectionFieldType(ctx) })
	inst.AddComputed("injection.getter.name", func(ctx *configen2.BuildingContext) (string, error) { return inst.propertyInjectionGetterName(ctx) })
	inst.AddComputed("injection.getter.inner", func(ctx *configen2.BuildingContext) (string, error) { return inst.propertyInjectionGetterInner(ctx) })

	return nil
}

func (inst *injectionGetterTemplate) Build(ctx *configen2.BuildingContext) (string, error) {

	selector := strings.TrimSpace(ctx.Injection.Selector)
	if selector == "" {
		return "", nil
	}

	return inst.BuildFromTemplate(ctx)
}

func (inst *injectionGetterTemplate) propertyComAdapterType(ctx *configen2.BuildingContext) (string, error) {
	return ctx.Component.StructName, nil
}

func (inst *injectionGetterTemplate) propertyInjectionGetterName(ctx *configen2.BuildingContext) (string, error) {
	return ctx.Injection.InjectionGetterMethod, nil
}

func (inst *injectionGetterTemplate) propertyInjectionFieldType(ctx *configen2.BuildingContext) (string, error) {
	return ctx.Injection.FieldType, nil
}

func (inst *injectionGetterTemplate) propertyInjectionGetterInner(ctx *configen2.BuildingContext) (string, error) {

	selector := strings.TrimSpace(ctx.Injection.Selector)
	fieldType := ctx.Injection.FieldType

	if inst.isSelectorForProperty(selector, fieldType) {
		return inst.group.injectionGetProperty.Build(ctx)

	} else if inst.isSelectorForList(selector, fieldType) {
		return inst.group.injectionGetList.Build(ctx)

	} else if inst.isSelectorForMap(selector, fieldType) {
		return inst.group.injectionGetMap.Build(ctx)

	} else if inst.isSelectorForContext(selector, fieldType) {
		return inst.group.injectionGetContext.Build(ctx)

	} else if inst.isSelectorForSimpleValue(selector, fieldType) {
		return inst.group.injectionGetSimpleValue.Build(ctx)

	} else if inst.isSelectorForObject(selector, fieldType) {
		return inst.group.injectionGetObject.Build(ctx)

	} else {
		// unsupported
	}

	return "", errors.New("unsupported selector: " + selector)
}

func (inst *injectionGetterTemplate) isSelectorForMap(selector string, fieldType string) bool {
	// class selector (map)
	index := strings.Index(fieldType, "]")
	b1 := (0 < index)
	b2 := strings.HasPrefix(fieldType, "map[")
	return b1 && b2
}

func (inst *injectionGetterTemplate) isSelectorForList(selector string, fieldType string) bool {
	// class selector (list)
	return strings.HasPrefix(fieldType, "[]")
}

func (inst *injectionGetterTemplate) isSelectorForContext(selector string, fieldType string) bool {
	return selector == "context"
}

func (inst *injectionGetterTemplate) isSelectorForProperty(selector string, fieldType string) bool {
	const prefix = "${"
	const suffix = "}"
	return strings.HasPrefix(selector, prefix) && strings.HasSuffix(selector, suffix)
}

func (inst *injectionGetterTemplate) isSelectorForObject(selector string, fieldType string) bool {
	// id selector
	return true
}

func (inst *injectionGetterTemplate) isSelectorForSimpleValue(selector string, fieldType string) bool {
	// id selector
	const prefix = "@value("
	const suffix = ")"
	return strings.HasPrefix(selector, prefix) && strings.HasSuffix(selector, suffix)
}

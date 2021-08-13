package templates

import "github.com/bitwormhole/starter-configen/tools/configen2"

type injectionGetMapTemplate struct {
	baseGroupTemplate
}

func (inst *injectionGetMapTemplate) _impl() innerTemplate {
	return inst
}

func (inst *injectionGetMapTemplate) InitWithGroup(group *innerTemplateGroup) error {

	err := inst.InitTemplate(group, "configen2/templates/injection_get_map.template")
	if err != nil {
		return err
	}

	//	inst.AddComputed("component.adapter.type", func(ctx *configen2.BuildingContext) (string, error) { return inst.propertyComAdapterType(ctx) })
	//	inst.AddComputed("injection.getter.name", func(ctx *configen2.BuildingContext) (string, error) { return inst.propertyInjectionGetterName(ctx) })
	//	inst.AddComputed("injection.field.type", func(ctx *configen2.BuildingContext) (string, error) { return inst.propertyInjectionFieldType(ctx) })

	return nil
}

func (inst *injectionGetMapTemplate) Build(ctx *configen2.BuildingContext) (string, error) {

	return inst.BuildFromTemplate(ctx)
}

func (inst *injectionGetMapTemplate) propertyComAdapterType(ctx *configen2.BuildingContext) (string, error) {
	return ctx.Component.StructName, nil
}

package templates

import "github.com/bitwormhole/starter-configen/tools/configen2"

type injectionGetObjectTemplate struct {
	baseGroupTemplate
}

func (inst *injectionGetObjectTemplate) _impl() innerTemplate {
	return inst
}

func (inst *injectionGetObjectTemplate) InitWithGroup(group *innerTemplateGroup) error {

	err := inst.InitTemplate(group, "configen2/templates/injection_get_object.template")
	if err != nil {
		return err
	}

	inst.AddComputed("injection.field.type", func(ctx *configen2.BuildingContext) (string, error) { return inst.propertyInjectionFieldType(ctx) })
	//	inst.AddComputed("injection.getter.name", func(ctx *configen2.BuildingContext) (string, error) { return inst.propertyInjectionGetterName(ctx) })
	//	inst.AddComputed("injection.field.type", func(ctx *configen2.BuildingContext) (string, error) { return inst.propertyInjectionFieldType(ctx) })

	return nil
}

func (inst *injectionGetObjectTemplate) Build(ctx *configen2.BuildingContext) (string, error) {

	return inst.BuildFromTemplate(ctx)
}

func (inst *injectionGetObjectTemplate) propertyInjectionFieldType(ctx *configen2.BuildingContext) (string, error) {
	return ctx.Injection.FieldType, nil
}

package templates

import "github.com/bitwormhole/starter-configen/tools/configen2"

type ComInfoTemplate struct {
	baseGroupTemplate
}

func (inst *ComInfoTemplate) _impl() innerTemplate {
	return inst
}

func (inst *ComInfoTemplate) InitWithGroup(group *innerTemplateGroup) error {

	err := inst.InitTemplate(group, "configen2/templates/com_info.template")
	if err != nil {
		return err
	}

	inst.AddComputed("component.name", func(ctx *configen2.BuildingContext) (string, error) { return inst.getComponentName(ctx) })
	inst.AddComputed("component.type", func(ctx *configen2.BuildingContext) (string, error) { return inst.getComponentType(ctx) })
	inst.AddComputed("component.id", func(ctx *configen2.BuildingContext) (string, error) { return inst.getComponentId(ctx) })
	inst.AddComputed("component.class", func(ctx *configen2.BuildingContext) (string, error) { return inst.getComponentClass(ctx) })
	inst.AddComputed("component.scope", func(ctx *configen2.BuildingContext) (string, error) { return inst.getComponentScope(ctx) })
	inst.AddComputed("component.aliases", func(ctx *configen2.BuildingContext) (string, error) { return inst.getComponentAliases(ctx) })

	inst.AddComputed("component.infunc.oninit", func(ctx *configen2.BuildingContext) (string, error) { return inst.getCodeInFuncOnInit(ctx) })
	inst.AddComputed("component.infunc.ondestroy", func(ctx *configen2.BuildingContext) (string, error) { return inst.getCodeInFuncOnDestroy(ctx) })
	inst.AddComputed("component.adapter.type", func(ctx *configen2.BuildingContext) (string, error) { return inst.getComAdapterType(ctx) })

	return nil
}

func (inst *ComInfoTemplate) Build(ctx *configen2.BuildingContext) (string, error) {
	return inst.BuildFromTemplate(ctx)
}

func (inst *ComInfoTemplate) getComponentName(ctx *configen2.BuildingContext) (string, error) {
	return ctx.Component.StructName, nil
}

func (inst *ComInfoTemplate) getComponentType(ctx *configen2.BuildingContext) (string, error) {
	return ctx.Component.ComType, nil
}

func (inst *ComInfoTemplate) getComponentId(ctx *configen2.BuildingContext) (string, error) {
	return ctx.Component.ID, nil
}

func (inst *ComInfoTemplate) getComponentScope(ctx *configen2.BuildingContext) (string, error) {
	return ctx.Component.Scope, nil
}

func (inst *ComInfoTemplate) getComponentClass(ctx *configen2.BuildingContext) (string, error) {
	return ctx.Component.Class, nil
}

func (inst *ComInfoTemplate) getComponentAliases(ctx *configen2.BuildingContext) (string, error) {
	return ctx.Component.Aliases, nil
}

func (inst *ComInfoTemplate) getCodeInFuncOnInit(ctx *configen2.BuildingContext) (string, error) {
	method := ctx.Component.InitMethod
	comType := ctx.Component.ComType
	if method == "" {
		return "return nil", nil
	}
	// return o.(*car.Body).Start()
	code := "return o.(*" + comType + ")." + method + "()"
	return code, nil
}

func (inst *ComInfoTemplate) getCodeInFuncOnDestroy(ctx *configen2.BuildingContext) (string, error) {
	method := ctx.Component.DestroyMethod
	comType := ctx.Component.ComType
	if method == "" {
		return "return nil", nil
	}
	// return o.(*car.Body).Start()
	code := "return o.(*" + comType + ")." + method + "()"
	return code, nil
}

func (inst *ComInfoTemplate) getComAdapterType(ctx *configen2.BuildingContext) (string, error) {
	return ctx.Component.StructName, nil
}

package templates

import (
	"sort"
	"strings"

	"github.com/bitwormhole/starter-configen/tools/configen2"
)

type injectionAdapterTemplate struct {
	baseGroupTemplate
}

func (inst *injectionAdapterTemplate) _impl() innerTemplate {
	return inst
}

func (inst *injectionAdapterTemplate) InitWithGroup(group *innerTemplateGroup) error {

	err := inst.InitTemplate(group, "configen2/templates/injection_adapter.template")
	if err != nil {
		return err
	}

	inst.AddComputed("component.adapter.type", func(ctx *configen2.BuildingContext) (string, error) { return inst.propComAdapterType(ctx) })
	inst.AddComputed("injection.from.list", func(ctx *configen2.BuildingContext) (string, error) { return inst.propInjectFromList(ctx) })
	inst.AddComputed("injection.to.list", func(ctx *configen2.BuildingContext) (string, error) { return inst.propInjectToList(ctx) })
	inst.AddComputed("invoke.custom.injectMethod", func(ctx *configen2.BuildingContext) (string, error) {
		return inst.propertyInvokeCustomInjectMethod(ctx)
	})

	return nil
}

func (inst *injectionAdapterTemplate) keysOf(table map[string]*configen2.Dom2injection) []string {
	keys := []string{}
	for key := range table {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (inst *injectionAdapterTemplate) buildInjectMethod(ctx *configen2.BuildingContext) (string, error) {
	code, err := inst.BuildFromTemplate(ctx)
	if err != nil {
		return "", err
	}
	return code, nil
}

func (inst *injectionAdapterTemplate) buildGetterMethod(ctx *configen2.BuildingContext) (string, error) {
	template := inst.group.injectionGetterTemplate
	return template.Build(ctx)
}

func (inst *injectionAdapterTemplate) buildGetterMethods(ctx *configen2.BuildingContext) (string, error) {

	builder := &strings.Builder{}
	com := ctx.Component
	injectionMap := com.InjectionMap
	keys := inst.keysOf(injectionMap)

	for index := range keys {
		key := keys[index]
		injection := injectionMap[key]
		child := ctx.NewChild()
		child.Injection = injection
		code, err := inst.buildGetterMethod(child)
		if err != nil {
			return "", err
		}
		builder.WriteString(code)
	}

	return builder.String(), nil
}

func (inst *injectionAdapterTemplate) Build(ctx *configen2.BuildingContext) (string, error) {

	// inject method
	code1, err := inst.buildInjectMethod(ctx)
	if err != nil {
		return "", err
	}

	//  getters methods
	code2, err := inst.buildGetterMethods(ctx)
	if err != nil {
		return "", err
	}

	return code1 + code2, nil
}

func (inst *injectionAdapterTemplate) propInjectFromList(ctx *configen2.BuildingContext) (string, error) {

	builder := &strings.Builder{}
	table := ctx.Component.InjectionMap
	keys := inst.keysOf(table)

	for index := range keys {
		key := keys[index]
		item := table[key]
		if item.Selector == "" {
			continue
		}
		fieldName := item.FieldName
		getterName := item.InjectionGetterMethod
		builder.WriteString(TAB)
		builder.WriteString("inst.")
		builder.WriteString(fieldName)
		builder.WriteString("=")
		builder.WriteString("inst.")
		builder.WriteString(getterName)
		builder.WriteString("(injection, \"")
		builder.WriteString(item.Selector)
		builder.WriteString("\")")
		builder.WriteString(NL)
	}

	return builder.String(), nil
}

func (inst *injectionAdapterTemplate) propertyInvokeCustomInjectMethod(ctx *configen2.BuildingContext) (string, error) {

	//	err = inst.__inject2__(injection)
	//	if err != nil {
	//		return err
	//	}

	method := strings.TrimSpace(ctx.Component.InjectMethod)
	if method == "" {
		return "", nil
	}

	builder := &strings.Builder{}
	builder.WriteString(TAB + "err = inst." + method + "(injection)" + NL)
	builder.WriteString(TAB + "if err !=nil {" + NL)
	builder.WriteString(TAB + "    return err" + NL)
	builder.WriteString(TAB + "}" + NL)
	return builder.String(), nil
}

func (inst *injectionAdapterTemplate) propInjectToList(ctx *configen2.BuildingContext) (string, error) {

	builder := &strings.Builder{}
	table := ctx.Component.InjectionMap
	keys := inst.keysOf(table)

	for index := range keys {
		key := keys[index]
		item := table[key]
		if item.Selector == "" {
			continue
		}
		if !item.Auto {
			continue
		}
		fieldName := item.FieldName
		builder.WriteString(TAB)
		builder.WriteString("instance.")
		builder.WriteString(fieldName)
		builder.WriteString("=")
		builder.WriteString("inst.")
		builder.WriteString(fieldName)
		builder.WriteString(NL)
	}

	return builder.String(), nil
}

func (inst *injectionAdapterTemplate) propComAdapterType(ctx *configen2.BuildingContext) (string, error) {
	return ctx.Component.StructName, nil
}

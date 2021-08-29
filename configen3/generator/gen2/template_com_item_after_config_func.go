package gen2

import (
	"strings"
	"text/template"

	"github.com/bitwormhole/starter-configen/configen3/data/model"
	"github.com/bitwormhole/starter-configen/configen3/mark"
)

type comItemAfterConfigFuncTemplate interface {
	BuildComponent(info *ComInfo) (string, error)
}

type comItemAfterConfigFuncTemplateData struct {
	ComID                string
	ComStructName        string
	ComProxyStructName   string
	ComFactoryStructName string

	InvokeInitMethod    string // return nil | ...
	InvokeDestroyMethod string // return nil | ...
	InvokeInjectMethod  string // return nil | ...
	ComSelectorsInit    string
	ComSelectors        string
	ComGetters          string
	// ComSetters          string
}

type comItemAfterConfigFuncTemplateImpl struct {
	context *gen2Context
	templ   *template.Template
	// data    mainTemplateData
}

func (inst *comItemAfterConfigFuncTemplateImpl) init(ctx *gen2Context) comItemAfterConfigFuncTemplate {

	tmp := ctx.loadTemplate("com-item-after-config-func.template", nil)

	inst.context = ctx
	inst.templ = tmp
	return inst
}

func (inst *comItemAfterConfigFuncTemplateImpl) makeGetterFuncName(f *model.StructField) string {
	name := f.Name
	return "getterForField" + name + "Selector"
}

func (inst *comItemAfterConfigFuncTemplateImpl) makeSelectorName(f *model.StructField) string {
	name := f.Name
	return "m" + name + "Selector"
}

func (inst *comItemAfterConfigFuncTemplateImpl) isInjectedField(f *model.StructField) bool {
	if f == nil {
		return false
	}
	inj := f.TagItems["inject"]
	if inj == "" {
		return false
	}
	name := f.Name
	if len(name) > 0 {
		ch0 := rune(name[0])
		if ('A' < ch0) && (ch0 < 'Z') {
			return true
		}
	}
	return false
}

func (inst *comItemAfterConfigFuncTemplateImpl) registerComAndFields(list []*FieldInfo, com *ComInfo) {
	ib := &inst.context.importsBuilder
	const noHash = false
	for _, f := range list {
		ib.AddComplexType(&f.Field.Type, noHash)
	}
	ib.AddComplexType(&com.InstanceField.Type, noHash)
}

func (inst *comItemAfterConfigFuncTemplateImpl) getNormalFields(info *ComInfo) []*FieldInfo {
	dst := make([]*FieldInfo, 0)
	src := info.Struct.Fields
	for _, field := range src {
		if !inst.isInjectedField(field) {
			continue
		}
		item2 := &FieldInfo{}
		item2.Field = field
		item2.Owner = info
		item2.SelectorName = inst.makeSelectorName(field)
		item2.SetterName = field.Name
		item2.GetterName = inst.makeGetterFuncName(field)
		item2.SelectorExpr = field.TagItems["inject"]
		dst = append(dst, item2)
	}
	return dst
}

func (inst *comItemAfterConfigFuncTemplateImpl) makeGetters(fields []*FieldInfo) (string, error) {
	ctx := inst.context
	getterTempl := ctx.templateFactory.createComFieldGetterTemplate(ctx)
	builder := strings.Builder{}
	for _, item := range fields {
		str, err := getterTempl.BuildGetter(item)
		if err != nil {
			return "", err
		}
		builder.WriteString("//" + item.GetterName)
		builder.WriteString(str)
		builder.WriteString(mark.NL)
	}
	return builder.String(), nil
}

func (inst *comItemAfterConfigFuncTemplateImpl) makeSelectors(fields []*FieldInfo) string {
	builder := strings.Builder{}
	builder.WriteString(mark.NL)
	for _, item := range fields {
		builder.WriteString(mark.TAB)
		builder.WriteString(item.SelectorName)
		builder.WriteString(" config.InjectionSelector")
		builder.WriteString(mark.NL)
	}
	return builder.String()
}

func (inst *comItemAfterConfigFuncTemplateImpl) makeFilterExprForField(field *FieldInfo) string {

	// func(name string, holder application.ComponentHolder) bool {
	// 	pt := holder.GetPrototype()
	// 	_, ok := pt.(*strings.Builder)
	// 	return ok
	// }

	elementTypeName := field.Field.Type.ValueType.String()

	builder := strings.Builder{}
	builder.WriteString("func(name string, holder application.ComponentHolder) bool {" + mark.NL)
	builder.WriteString("            pt := holder.GetPrototype()" + mark.NL)
	builder.WriteString("            _, ok := pt.(" + elementTypeName + ")" + mark.NL)
	builder.WriteString("            return ok" + mark.NL)
	builder.WriteString("        }")
	return builder.String()
}

func (inst *comItemAfterConfigFuncTemplateImpl) makeSelectorsInit(fields []*FieldInfo) string {
	builder := strings.Builder{}
	builder.WriteString(mark.NL)
	for _, item := range fields {
		// example: inst.F1 = config.NewInjectionSelector("#abc", nil)
		selector := strings.TrimSpace(item.SelectorExpr)
		filterExpr := "nil"
		if selector == "*" {
			filterExpr = inst.makeFilterExprForField(item)
		}
		builder.WriteString(mark.TAB)
		builder.WriteString("inst.")
		builder.WriteString(item.SelectorName)
		builder.WriteString(" = config.NewInjectionSelector(")
		builder.WriteString(mark.QUOTE2 + selector + mark.QUOTE2)
		builder.WriteString("," + filterExpr + ")")
		builder.WriteString(mark.NL)
	}
	return builder.String()
}

func (inst *comItemAfterConfigFuncTemplateImpl) BuildComponent(info *ComInfo) (string, error) {

	fields := inst.getNormalFields(info)
	inst.registerComAndFields(fields, info)

	strGetters, err := inst.makeGetters(fields)
	if err != nil {
		return "", err
	}

	data := &comItemAfterConfigFuncTemplateData{}
	data.ComID = info.ComID
	data.ComStructName = inst.makeComStructName(info)
	data.ComFactoryStructName = info.FactoryStructName
	data.ComProxyStructName = info.ProxyStructName
	data.InvokeInitMethod = inst.makeInvokeInitMethod(info)
	data.InvokeDestroyMethod = inst.makeInvokeDestroyMethod(info)
	data.InvokeInjectMethod = inst.makeInvokeInjectMethod(info, fields)

	data.ComSelectorsInit = inst.makeSelectorsInit(fields)
	data.ComSelectors = inst.makeSelectors(fields)
	data.ComGetters = strGetters

	return executeTemplate(inst.templ, data)
}

func (inst *comItemAfterConfigFuncTemplateImpl) makeComStructName(info *ComInfo) string {
	field := info.InstanceField
	ct := field.Type.Clone(true)
	st := ct.ValueType
	st.HasStar = false
	str := st.String()
	return str
}

func (inst *comItemAfterConfigFuncTemplateImpl) makeInvokeInitMethod(info *ComInfo) string {
	method := strings.TrimSpace(info.ComInitMethod)
	if method == "" {
		return "return nil"
	}
	return "return inst.castObject(instance)." + method + "()"
}

func (inst *comItemAfterConfigFuncTemplateImpl) makeInvokeDestroyMethod(info *ComInfo) string {
	method := strings.TrimSpace(info.ComDestroyMethod)
	if method == "" {
		return "return nil"
	}
	return "return inst.castObject(instance)." + method + "()"
}

func (inst *comItemAfterConfigFuncTemplateImpl) makeInvokeInjectMethod(info *ComInfo, fields []*FieldInfo) string {

	const nop = "return nil"
	if info == nil || fields == nil {
		return nop
	}
	if len(fields) < 1 {
		return nop
	}

	// 如果有字段需要注入，执行下列过程

	builder := strings.Builder{}
	builder.WriteString(mark.NL + mark.TAB)
	builder.WriteString("obj := inst.castObject(instance)" + mark.NL)

	for _, item := range fields {
		builder.WriteString(mark.TAB + "obj." + item.SetterName + " = inst.")
		builder.WriteString(item.GetterName)
		builder.WriteString("(context)" + mark.NL)
	}

	builder.WriteString(mark.TAB + "return context.LastError()")
	return builder.String()
}

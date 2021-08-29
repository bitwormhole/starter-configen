package gen2

import (
	"errors"
	"text/template"

	"github.com/bitwormhole/starter-configen/configen3/data/model"
	"github.com/bitwormhole/starter/util"
)

type comFieldGetterTemplate interface {
	BuildGetter(info *FieldInfo) (string, error)
}

type comFieldGetterTemplateData struct {

	// ComID      string
	// ComScope   string
	// ComClass   string
	// ComAliases string
	// ComFactoryStructName string

	// for all
	TemplatePath string

	StructName   string
	MethodName   string
	ReturnExpr   string
	SelectorName string
	ComID        string
	FieldName    string

	// for simple value
	SimpleValueGetterMethod string

	// for object list
	ElementType string
}

type comFieldGetterTemplateImpl struct {
	context *gen2Context
	templ   *template.Template
	// data    mainTemplateData
}

func (inst *comFieldGetterTemplateImpl) init(ctx *gen2Context) comFieldGetterTemplate {

	inst.context = ctx
	return inst
}

func (inst *comFieldGetterTemplateImpl) BuildGetter(info *FieldInfo) (string, error) {

	ft := &info.Field.Type
	vt := ft.ValueType
	selector := info.Field.TagItems["inject"]
	data := &comFieldGetterTemplateData{}
	var err error = nil
	comID := info.Owner.ComID

	if vt == nil {
		panic("vt==nil, com-id:" + comID)
	}

	err = inst.prepareCommon(info, data)
	if err != nil {
		return "", err
	}

	if ft.IsMap {
		// as map[k]v
		return "", errors.New("inject map[key]value is unsupported, com-id:" + comID)

	} else if ft.IsArray {
		//  as []element
		err = inst.prepareForObjectList(info, data)

	} else if selector == "context" {
		// as context
		err = inst.prepareForContext(info, data)

	} else if selector == "pool" {
		// as pool
		err = inst.prepareForPool(info, data)

	} else if vt.PackageName == "" {
		// simple value
		err = inst.prepareForSimpleValue(info, data)

	} else {
		// simple object
		err = inst.prepareForSimpleObject(info, data)
	}

	if err != nil {
		return "", err
	}

	return inst.doBuild(info, data)
}

func (inst *comFieldGetterTemplateImpl) doBuild(info *FieldInfo, data *comFieldGetterTemplateData) (string, error) {

	templ := inst.context.loadTemplate(data.TemplatePath, nil)
	return executeTemplate(templ, data)
}

func (inst *comFieldGetterTemplateImpl) prepareCommon(info *FieldInfo, data *comFieldGetterTemplateData) error {

	com := info.Owner
	field := info.Field

	data.StructName = com.FactoryStructName
	data.MethodName = info.GetterName
	data.ReturnExpr = inst.makeGetterReturnExpr(field)
	data.ComID = com.ComID
	data.FieldName = field.Name
	data.SelectorName = info.SelectorName

	return nil
}

func (inst *comFieldGetterTemplateImpl) makeGetterReturnExpr(field *model.StructField) string {
	ft := field.Type
	return ft.String()
}

var theSimpleValueGetterTypeToMethodTable map[string]string

func (inst *comFieldGetterTemplateImpl) getSimpleValueGetterTypeToMethodTable() map[string]string {

	table := theSimpleValueGetterTypeToMethodTable
	if table == nil {
		table = make(map[string]string)

		table["int"] = "GetInt"
		table["int8"] = "GetInt8"
		table["int16"] = "GetInt16"
		table["int32"] = "GetInt32"
		table["int64"] = "GetInt64"

		table["float32"] = "GetFloat32"
		table["float64"] = "GetFloat64"

		table["bool"] = "GetBool"
		table["string"] = "GetString"
		theSimpleValueGetterTypeToMethodTable = table
	}
	return table
}

func (inst *comFieldGetterTemplateImpl) prepareForSimpleValue(info *FieldInfo, data *comFieldGetterTemplateData) error {

	tName := info.Field.Type.ValueType.TypeName
	table := inst.getSimpleValueGetterTypeToMethodTable()
	method := table[tName]

	if method == "" {
		eb := &util.ErrorBuilder{}
		eb.Message("unsupported type")
		eb.Set("type", tName)
		eb.Set("comID", info.Owner.ComID)
		return eb.Create()
	}

	data.SimpleValueGetterMethod = method
	data.SelectorName = info.SelectorName
	data.TemplatePath = "getters/simple-value-getter.template"
	return nil
}

func (inst *comFieldGetterTemplateImpl) prepareForSimpleObject(info *FieldInfo, data *comFieldGetterTemplateData) error {

	data.TemplatePath = "getters/simple-object-getter.template"
	return nil
}

func (inst *comFieldGetterTemplateImpl) prepareForObjectList(info *FieldInfo, data *comFieldGetterTemplateData) error {

	vt := info.Field.Type.ValueType

	data.ElementType = vt.String()
	data.TemplatePath = "getters/object-list-getter.template"
	return nil
}

func (inst *comFieldGetterTemplateImpl) prepareForContext(info *FieldInfo, data *comFieldGetterTemplateData) error {
	data.TemplatePath = "getters/context-getter.template"
	return nil
}

func (inst *comFieldGetterTemplateImpl) prepareForPool(info *FieldInfo, data *comFieldGetterTemplateData) error {
	data.TemplatePath = "getters/pool-getter.template"
	return nil
}

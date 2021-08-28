package gen1

import (
	"errors"
	"strings"
	"text/template"

	"github.com/bitwormhole/starter-configen/configen3/data/model"
	"github.com/bitwormhole/starter-configen/configen3/mark"
)

type structListTemplate interface {
	BuildStructs(s *model.GoSource) (string, error)
}

type structListTemplateImpl struct {
	context  *gen1Context
	template *template.Template
}

func (inst *structListTemplateImpl) init(ctx *gen1Context) structListTemplate {

	inst.context = ctx
	return inst
}

func (inst *structListTemplateImpl) BuildStructs(s *model.GoSource) (string, error) {

	list := s.StructList
	builder := &strings.Builder{}

	if list == nil {
		return "", nil
	}

	for _, item := range list {
		builder.WriteString(mark.NL)
		err := inst.buildStruct(item, builder)
		builder.WriteString(mark.NL)
		if err != nil {
			return "", err
		}
	}

	return builder.String(), nil
}

func (inst *structListTemplateImpl) registerSimpleType(st *model.SimpleType) *model.SimpleType {
	if st == nil {
		return nil
	}
	imp := &model.GoImport{}
	imp.Path = st.PackageName
	imp.Alias = st.PackageAlias
	if imp.Path == "" {
		return st
	}
	result := inst.context.importsBuilder.Add(imp)
	st.PackageAlias = result.Alias
	st.PackageName = result.Path
	return st
}

func (inst *structListTemplateImpl) makeFieldTypeExpr(sf *model.StructField) (string, error) {

	src := sf.Type
	dst := &model.ComplexType{}

	dst.IsArray = src.IsArray
	dst.IsMap = src.IsMap
	dst.KeyType = inst.registerSimpleType(src.KeyType)
	dst.ValueType = inst.registerSimpleType(src.ValueType)

	return dst.String(), nil
}

func (inst *structListTemplateImpl) makeStructTypeExpr(aStruct *model.Struct) (string, error) {
	aType := inst.registerSimpleType(&aStruct.Type)
	element := aStruct.Name
	expr := "*" + aType.PackageAlias + "." + element
	return expr, nil
}

func (inst *structListTemplateImpl) buildStruct(st *model.Struct, builder *strings.Builder) error {

	name := "pCom" + st.Type.TypeName
	typeExpr, err := inst.makeStructTypeExpr(st)
	if err != nil {
		return err
	}

	builder.WriteString("type " + name + " struct {" + mark.NL)
	builder.WriteString(mark.TAB + "instance " + typeExpr + mark.NL)

	// 由于原生结构已经包含了markup.X，所以这里不再需要 markup.Component
	//	builder.WriteString(mark.TAB + "markup.Component" + mark.NL)

	fields := st.Fields
	if fields == nil {
		return errors.New("no filed")
	}

	for _, item := range fields {
		err := inst.buildField(item, builder)
		if err != nil {
			return err
		}
	}

	builder.WriteString("}" + mark.NL)
	return nil
}

func (inst *structListTemplateImpl) buildField(f *model.StructField, builder *strings.Builder) error {

	name := f.Name
	typeExpr, err := inst.makeFieldTypeExpr(f)
	tag := f.Tag

	if err != nil {
		return err
	}

	builder.WriteString(mark.TAB + name + mark.SPACE)
	builder.WriteString(typeExpr + mark.SPACE)
	builder.WriteString(mark.QUOTE0 + tag + mark.QUOTE0)
	builder.WriteString(mark.NL)

	return nil
}

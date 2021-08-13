package configen2

import (
	"errors"
	"fmt"
	"strings"
)

type Dom2Builder interface {
	Build(ctx *Context) error
}

type defaultDom2Builder struct {
	context *Context
}

func (inst *defaultDom2Builder) init(ctx *Context) Dom2Builder {
	inst.context = ctx
	return inst
}

func (inst *defaultDom2Builder) Build(ctx *Context) error {

	dom1 := &ctx.DOM1
	dom2 := &ctx.DOM2

	dom2.PackageName = dom1.PackageName

	err := inst.buildImports(dom1, dom2)
	if err != nil {
		return err
	}

	err = inst.buildComponents(dom1, dom2)
	if err != nil {
		return err
	}

	return nil
}

func (inst *defaultDom2Builder) buildImports(dom1 *Dom1root, dom2 *Dom2root) error {

	table1 := dom1.Imports
	im := inst.context.Imports

	im.AddImportWithoutHash("strings")
	im.AddImportWithoutHash("errors")
	im.AddImportWithoutHash("github.com/bitwormhole/starter/lang")
	im.AddImportWithoutHash("github.com/bitwormhole/starter/application")
	im.AddImportWithoutHash("github.com/bitwormhole/starter/application/config")

	for key := range table1 {
		path := key
		im.AddImport(path)
	}
	return nil
}

func (inst *defaultDom2Builder) buildComponents(dom1 *Dom1root, dom2 *Dom2root) error {

	docs := dom1.Documents
	if docs == nil {
		return nil
	}

	for index := range docs {
		doc := docs[index]
		err := inst.buildComponentsInDoc(doc, dom1, dom2)
		if err != nil {
			return err
		}
	}

	return nil
}

func (inst *defaultDom2Builder) buildComponentsInDoc(doc *Dom1doc, dom1 *Dom1root, dom2 *Dom2root) error {

	list1 := doc.ComList
	table2 := dom2.Components

	for index := range list1 {
		item := list1[index]
		if inst.isComponent(item) {
			com, err := inst.makeComponent(item, doc)
			if err != nil {
				return err
			}
			table2[com.StructName] = com
		}
	}

	return nil
}

func (inst *defaultDom2Builder) isComponent(st *Dom1struct) bool {

	if st == nil {
		return false
	}

	fields := st.Fields
	if fields == nil {
		return false
	}

	for index := range fields {
		f := fields[index]
		if f.Name == "" && f.Type == "markup.Component" {
			return true
		}
	}

	return false
}

func (inst *defaultDom2Builder) makeComponent(st *Dom1struct, doc *Dom1doc) (*Dom2component, error) {

	com := &Dom2component{}
	com.StructName = st.Name
	com.InjectionMap = make(map[string]*Dom2injection)
	com.Attributes = make(map[string]string)
	fields := st.Fields

	for index := range fields {
		f := fields[index]
		injection, err := inst.makeInjection(f, doc)
		if err != nil {
			return nil, err
		}
		com.InjectionMap[injection.FieldName] = injection
	}

	// 混合（markup.Component, instance, context）的属性
	err := inst.mixComponentAtts(com)
	if err != nil {
		return nil, err
	}

	// 加载 Component 的基本属性
	err = inst.loadComponentAtts(com)
	if err != nil {
		return nil, err
	}

	return com, nil
}

func (inst *defaultDom2Builder) mixComponentAtts(com *Dom2component) error {
	atts2 := map[string]string{}
	injectionTable := com.InjectionMap
	if injectionTable == nil {
		return errors.New("no injection")
	}
	for key := range injectionTable {
		injection := injectionTable[key]
		if injection == nil {
			continue
		}
		if key == "" || key == "instance" || key == "context" {
			// do loop
		} else {
			continue
		}
		atts1 := injection.Attributes
		if atts1 == nil {
			continue
		}
		for k2 := range atts1 {
			v2 := atts1[k2]
			older := atts2[k2]
			if older != "" {
				msg := fmt.Sprint("the attribute is duplicated. Attr:", k2, " ComStruct:", com.StructName)
				return errors.New(msg)
			}
			atts2[k2] = v2
		}
	}
	com.Attributes = atts2
	return nil
}

func (inst *defaultDom2Builder) loadComponentAtts(com *Dom2component) error {

	instance := com.InjectionMap["instance"]
	table := com.Attributes
	comId := table["id"]

	if comId == "" {
		comId = com.StructName
	}

	com.ComType = strings.ReplaceAll(instance.FieldType, "*", "")

	com.InjectMethod = table["injectMethod"]
	com.InitMethod = table["initMethod"]
	com.DestroyMethod = table["destroyMethod"]

	com.ID = comId
	com.Class = table["class"]
	com.Scope = table["scope"]
	com.Aliases = table["aliases"]

	return nil
}

func (inst *defaultDom2Builder) makeInjection(f *Dom1field, doc1 *Dom1doc) (*Dom2injection, error) {

	injection := &Dom2injection{}

	injection.FieldName = f.Name
	injection.FieldType = inst.convertComplexType(f.Type, doc1)
	injection.FieldTag = f.Tag
	injection.Auto = inst.isAutoInjection(f.Name)
	injection.Attributes = inst.parseTag(f.Tag)
	injection.Selector = injection.Attributes["inject"]

	return injection, nil
}

func (inst *defaultDom2Builder) convertComplexType(t string, doc1 *Dom1doc) string {
	ctbuilder := &ComplexTypeBuilder{}
	ctbuilder.Init(t)
	list := ctbuilder.FindPackageAliases()
	for i := range list {
		index := list[i]
		alias := ctbuilder.GetPart(index)
		alias = inst.convertImportPackageAlias(alias, doc1)
		ctbuilder.SetPart(index, alias)
	}
	return ctbuilder.String()
}

func (inst *defaultDom2Builder) convertImportPackageAlias(alias string, doc *Dom1doc) string {

	// find package full path
	list := doc.ImportList
	if list == nil {
		return alias
	}
	path := ""
	for i := range list {
		imp := list[i]
		if imp.Alias == alias {
			path = imp.Path
			break
		}
	}
	if path == "" {
		return alias
	}

	// find alias2 by path
	alias2, err := inst.context.Imports.FindAliasByPath(path)
	if err == nil {
		return alias2
	} else {
		return alias
	}
}

func (inst *defaultDom2Builder) isAutoInjection(name string) bool {
	// 如果字段名是大写字母开头，表示自动注入
	size := len(name)
	if size > 0 {
		ch := name[0]
		if 'A' <= ch && ch <= 'Z' {
			return true
		}
	}
	return false
}

func (inst *defaultDom2Builder) parseTag(tag string) map[string]string {
	table := make(map[string]string)
	reader := &FieldTagReader{}
	reader.init(tag)
	for {
		key, value, cnt := reader.read()
		if cnt < 0 {
			break
		}
		table[key] = value
	}
	return table
}

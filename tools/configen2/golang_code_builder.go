package configen2

import (
	"errors"

	"github.com/bitwormhole/starter/application"
)

type CodeBuilder interface {
	Build(dom *Dom2root) (string, error)
}

////////////////////////////////////////////////////////////////////////////////

type golangCodeBuilder struct {
	templateFactory CodeTemplateFactory
	context         application.Context
}

func (inst *golangCodeBuilder) _impl() CodeBuilder {
	return inst
}

func (inst *golangCodeBuilder) init(ctx *Context) error {

	ac := ctx.AppContext
	obj, err := ac.GetComponent("#configen2-main-template-factory")
	if err != nil {
		return err
	}

	factory, ok := obj.(CodeTemplateFactory)
	if !ok {
		return errors.New("com.(CodeTemplateFactory) return false")
	}

	inst.context = ctx.AppContext
	inst.templateFactory = factory
	return nil
}

func (inst *golangCodeBuilder) prepareInjectMethods1(dom *Dom2root) {
	table := dom.Components
	for key := range table {
		com := table[key]
		inst.prepareInjectMethods2(com)
	}
}

func (inst *golangCodeBuilder) prepareInjectMethods2(dom *Dom2component) {
	table := dom.InjectionMap
	for key := range table {
		item := table[key]
		name := item.FieldName
		item.InjectionGetterMethod = "__get_" + name + "__"
	}
}

func (inst *golangCodeBuilder) Build(dom *Dom2root) (string, error) {

	inst.prepareInjectMethods1(dom)

	// context
	bc := &BuildingContext{}
	bc.DOM = dom
	template, err := inst.templateFactory.Create(inst.context)
	if err != nil {
		return "", err
	}
	return template.Build(bc)
}

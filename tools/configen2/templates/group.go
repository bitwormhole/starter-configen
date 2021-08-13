package templates

import (
	"github.com/bitwormhole/starter-configen/tools/configen2"
	"github.com/bitwormhole/starter/application"
)

type innerTemplate interface {
	configen2.CodeTemplate
	InitWithGroup(group *innerTemplateGroup) error
}

////////////////////////////////////////////////////////////////////////////////
// innerTemplateGroup

type innerTemplateGroup struct {
	context application.Context

	mainTemplate             *MainTemplate
	comInfoTemplate          *ComInfoTemplate
	injectionAdapterTemplate *injectionAdapterTemplate

	injectionGetterTemplate *injectionGetterTemplate

	injectionGetContext     *injectionGetContextTemplate
	injectionGetProperty    *injectionGetPropertyTemplate
	injectionGetObject      *injectionGetObjectTemplate
	injectionGetList        *injectionGetListTemplate
	injectionGetMap         *injectionGetMapTemplate
	injectionGetSimpleValue *injectionGetObjectTemplate // todo ...
}

////////////////////////////////////////////////////////////////////////////////
// baseGroupTemplate

type baseGroupTemplate struct {
	configen2.TemplateCodeBuilderBase
	group *innerTemplateGroup
}

func (inst *baseGroupTemplate) InitTemplate(group *innerTemplateGroup, path string) error {

	inst.group = group

	err := inst.TemplateCodeBuilderBase.Init(group.context)
	if err != nil {
		return err
	}

	err = inst.TemplateCodeBuilderBase.LoadTemplate(group.context, path)
	if err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// MainTemplateFactory

type MainTemplateFactory struct {
}

func (inst *MainTemplateFactory) _impl() configen2.CodeTemplateFactory {
	return inst
}

func (inst *MainTemplateFactory) Create(ac application.Context) (configen2.CodeTemplate, error) {

	group := &innerTemplateGroup{context: ac}
	templist := make([]innerTemplate, 0)

	// create template
	group.mainTemplate = &MainTemplate{}
	group.comInfoTemplate = &ComInfoTemplate{}
	group.injectionAdapterTemplate = &injectionAdapterTemplate{}

	group.injectionGetterTemplate = &injectionGetterTemplate{}

	group.injectionGetContext = &injectionGetContextTemplate{}
	group.injectionGetList = &injectionGetListTemplate{}
	group.injectionGetMap = &injectionGetMapTemplate{}
	group.injectionGetObject = &injectionGetObjectTemplate{}
	group.injectionGetProperty = &injectionGetPropertyTemplate{}
	group.injectionGetSimpleValue = &injectionGetObjectTemplate{}

	// add to templist
	templist = append(templist, group.mainTemplate)
	templist = append(templist, group.comInfoTemplate)
	templist = append(templist, group.injectionAdapterTemplate)

	templist = append(templist, group.injectionGetterTemplate)

	templist = append(templist, group.injectionGetContext)
	templist = append(templist, group.injectionGetList)
	templist = append(templist, group.injectionGetMap)
	templist = append(templist, group.injectionGetObject)
	templist = append(templist, group.injectionGetProperty)
	templist = append(templist, group.injectionGetSimpleValue)

	// init
	for index := range templist {
		temp := templist[index]
		err := temp.InitWithGroup(group)
		if err != nil {
			return nil, err
		}
	}

	// result
	return group.mainTemplate, nil
}

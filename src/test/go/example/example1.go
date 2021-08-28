package example

import (
	"strings"

	"github.com/bitwormhole/starter/application"
	"github.com/bitwormhole/starter/application/config"
	"github.com/bitwormhole/starter/lang"
	"github.com/bitwormhole/starter/util"
)

func autoConfig(cb application.ConfigBuilder) error {

	var err error = nil
	cominfobuilder := config.ComInfo()

	// id:
	cominfobuilder.Next()
	cominfobuilder.ID("").Class("").Aliases("").Scope("")
	cominfobuilder.Factory(&myFactory{})
	err = cominfobuilder.CreateTo(cb)
	if err != nil {
		return err
	}

	return nil
}

type myFactory struct {
	mPrototype *strings.Builder
	F1         config.InjectionSelector
}

func (inst *myFactory) init() application.ComponentFactory {

	inst.F1 = config.NewInjectionSelector("#abc", func(name string, holder application.ComponentHolder) bool {
		pt := holder.GetPrototype()
		_, ok := pt.(*strings.Builder)
		return ok
	})

	return inst
}

func (inst *myFactory) newObject() *strings.Builder {
	return &strings.Builder{}
}

func (inst *myFactory) cast(instance application.ComponentInstance) *strings.Builder {
	o := instance.Get()
	return o.(*strings.Builder)
}

func (inst *myFactory) GetPrototype() lang.Object {
	obj := inst.mPrototype
	if obj == nil {
		obj = inst.newObject()
		inst.mPrototype = obj
	}
	return obj
}

func (inst *myFactory) NewInstance() application.ComponentInstance {
	return config.SimpleInstance(inst, inst.newObject())
}

func (inst *myFactory) AfterService() application.ComponentAfterService {
	return inst
}

func (inst *myFactory) Inject(instance application.ComponentInstance, context application.InstanceContext) error {
	o := inst.cast(instance)
	o.Len() // =
	inst.__getter__(context)
	return nil
}

func (inst *myFactory) Init(instance application.ComponentInstance) error {
	inst.cast(instance).Len()
	return nil
}

func (inst *myFactory) Destroy(instance application.ComponentInstance) error {
	inst.cast(instance).Len()
	return nil
}

func (inst *myFactory) __getter__(context application.InstanceContext) int {
	return inst.F1.GetInt(context)
}

func (inst *myFactory) __getter2__(context application.InstanceContext) *strings.Builder {
	o1 := inst.F1.GetOne(context)
	o2, ok := o1.(*strings.Builder)
	if !ok {
		eb := &util.ErrorBuilder{}
		eb.Message("bad cast")
		eb.Set("com", "")
		eb.Set("field", "")
		eb.Set("type1", "")
		eb.Set("type2", "")
		context.HandleError(eb.Create())
		return nil
	}
	return o2
}

func (inst *myFactory) __getter3__(context application.InstanceContext) []*strings.Builder {
	list1 := inst.F1.GetList(context)
	list2 := make([]*strings.Builder, 0, len(list1))
	for _, item1 := range list1 {
		item2, ok := item1.(*strings.Builder)
		if ok {
			list2 = append(list2, item2)
		}
	}
	return list2
}

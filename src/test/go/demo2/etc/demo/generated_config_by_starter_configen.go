// (todo:gen2.template) 
// 这个配置文件是由 starter-configen 工具自动生成的。
// 任何时候，都不要手工修改这里面的内容！！！

package demo

import (
	deep0x1fa623 "github.com/bitwormhole/starter-configen/src/test/go/demo2/code/deep"
	application "github.com/bitwormhole/starter/application"
	config "github.com/bitwormhole/starter/application/config"
	lang "github.com/bitwormhole/starter/lang"
	util "github.com/bitwormhole/starter/util"
	strings0x1877f3 "strings"
    
)

func autoConfig (cb application.ConfigBuilder) error {

	var err error = nil
	cominfobuilder := config.ComInfo()

	// component: com0-deep0x1fa623.Example1
	cominfobuilder.Next()
	cominfobuilder.ID("com0-deep0x1fa623.Example1").Class("").Aliases("").Scope("")
	cominfobuilder.Factory((&comFactory4pComExample1{}).init())
	err = cominfobuilder.CreateTo(cb)
	if err != nil {
		return err
	}

	// component: Example2
	cominfobuilder.Next()
	cominfobuilder.ID("Example2").Class("Example").Aliases("").Scope("")
	cominfobuilder.Factory((&comFactory4pComExample2{}).init())
	err = cominfobuilder.CreateTo(cb)
	if err != nil {
		return err
	}

	// component: com2-deep0x1fa623.Example3
	cominfobuilder.Next()
	cominfobuilder.ID("com2-deep0x1fa623.Example3").Class("example demo element").Aliases("x y z").Scope("singleton")
	cominfobuilder.Factory((&comFactory4pComExample3{}).init())
	err = cominfobuilder.CreateTo(cb)
	if err != nil {
		return err
	}



    return nil
}

////////////////////////////////////////////////////////////////////////////////

// comFactory4pComExample1 : the factory of component: com0-deep0x1fa623.Example1
type comFactory4pComExample1 struct {

    mPrototype * deep0x1fa623.Example1

	
	mF1Selector config.InjectionSelector
	mF2Selector config.InjectionSelector
	mF3Selector config.InjectionSelector
	mF4Selector config.InjectionSelector
	mF5Selector config.InjectionSelector
	mF6i8Selector config.InjectionSelector
	mF6i16Selector config.InjectionSelector
	mF6i32Selector config.InjectionSelector
	mF6i64Selector config.InjectionSelector
	mF7Selector config.InjectionSelector
	mF8Selector config.InjectionSelector
	mF9Selector config.InjectionSelector
	mF10Selector config.InjectionSelector
	mF11Selector config.InjectionSelector

}

func (inst * comFactory4pComExample1) init() application.ComponentFactory {

	
	inst.mF1Selector = config.NewInjectionSelector("context",nil)
	inst.mF2Selector = config.NewInjectionSelector("pool",nil)
	inst.mF3Selector = config.NewInjectionSelector("${test.str.s1}",nil)
	inst.mF4Selector = config.NewInjectionSelector("hello,world",nil)
	inst.mF5Selector = config.NewInjectionSelector("1000",nil)
	inst.mF6i8Selector = config.NewInjectionSelector("${test.num.i64}",nil)
	inst.mF6i16Selector = config.NewInjectionSelector("${test.num.i64}",nil)
	inst.mF6i32Selector = config.NewInjectionSelector("${test.num.i64}",nil)
	inst.mF6i64Selector = config.NewInjectionSelector("${test.num.i64}",nil)
	inst.mF7Selector = config.NewInjectionSelector("false",nil)
	inst.mF8Selector = config.NewInjectionSelector("${test.num.f32}",nil)
	inst.mF9Selector = config.NewInjectionSelector("0.001",nil)
	inst.mF10Selector = config.NewInjectionSelector("*",func(name string, holder application.ComponentHolder) bool {
            pt := holder.GetPrototype()
            _, ok := pt.(*strings0x1877f3.Builder)
            return ok
        })
	inst.mF11Selector = config.NewInjectionSelector("*",func(name string, holder application.ComponentHolder) bool {
            pt := holder.GetPrototype()
            _, ok := pt.(*strings0x1877f3.Builder)
            return ok
        })


	inst.mPrototype = inst.newObject()
    return inst
}

func (inst * comFactory4pComExample1) newObject() * deep0x1fa623.Example1 {
	return & deep0x1fa623.Example1 {}
}

func (inst * comFactory4pComExample1) castObject(instance application.ComponentInstance) * deep0x1fa623.Example1 {
	return instance.Get().(*deep0x1fa623.Example1)
}

func (inst * comFactory4pComExample1) GetPrototype() lang.Object {
	return inst.mPrototype
}

func (inst * comFactory4pComExample1) NewInstance() application.ComponentInstance {
	return config.SimpleInstance(inst, inst.newObject())
}

func (inst * comFactory4pComExample1) AfterService() application.ComponentAfterService {
	return inst
}

func (inst * comFactory4pComExample1) Init(instance application.ComponentInstance) error {
	return inst.castObject(instance).Start()
}

func (inst * comFactory4pComExample1) Destroy(instance application.ComponentInstance) error {
	return nil
}

func (inst * comFactory4pComExample1) Inject(instance application.ComponentInstance, context application.InstanceContext) error {
	
	obj := inst.castObject(instance)
	obj.F1 = inst.getterForFieldF1Selector(context)
	obj.F2 = inst.getterForFieldF2Selector(context)
	obj.F3 = inst.getterForFieldF3Selector(context)
	obj.F4 = inst.getterForFieldF4Selector(context)
	obj.F5 = inst.getterForFieldF5Selector(context)
	obj.F6i8 = inst.getterForFieldF6i8Selector(context)
	obj.F6i16 = inst.getterForFieldF6i16Selector(context)
	obj.F6i32 = inst.getterForFieldF6i32Selector(context)
	obj.F6i64 = inst.getterForFieldF6i64Selector(context)
	obj.F7 = inst.getterForFieldF7Selector(context)
	obj.F8 = inst.getterForFieldF8Selector(context)
	obj.F9 = inst.getterForFieldF9Selector(context)
	obj.F10 = inst.getterForFieldF10Selector(context)
	obj.F11 = inst.getterForFieldF11Selector(context)
	return context.LastError()
}

//getterForFieldF1Selector
func (inst * comFactory4pComExample1) getterForFieldF1Selector (context application.InstanceContext) application.Context {
    return context.Context()
}

//getterForFieldF2Selector
func (inst * comFactory4pComExample1) getterForFieldF2Selector (context application.InstanceContext) lang.ReleasePool {
    return context.Pool()
}

//getterForFieldF3Selector
func (inst * comFactory4pComExample1) getterForFieldF3Selector (context application.InstanceContext) string {
    return inst.mF3Selector.GetString(context)
}

//getterForFieldF4Selector
func (inst * comFactory4pComExample1) getterForFieldF4Selector (context application.InstanceContext) string {
    return inst.mF4Selector.GetString(context)
}

//getterForFieldF5Selector
func (inst * comFactory4pComExample1) getterForFieldF5Selector (context application.InstanceContext) int {
    return inst.mF5Selector.GetInt(context)
}

//getterForFieldF6i8Selector
func (inst * comFactory4pComExample1) getterForFieldF6i8Selector (context application.InstanceContext) int8 {
    return inst.mF6i8Selector.GetInt8(context)
}

//getterForFieldF6i16Selector
func (inst * comFactory4pComExample1) getterForFieldF6i16Selector (context application.InstanceContext) int16 {
    return inst.mF6i16Selector.GetInt16(context)
}

//getterForFieldF6i32Selector
func (inst * comFactory4pComExample1) getterForFieldF6i32Selector (context application.InstanceContext) int32 {
    return inst.mF6i32Selector.GetInt32(context)
}

//getterForFieldF6i64Selector
func (inst * comFactory4pComExample1) getterForFieldF6i64Selector (context application.InstanceContext) int64 {
    return inst.mF6i64Selector.GetInt64(context)
}

//getterForFieldF7Selector
func (inst * comFactory4pComExample1) getterForFieldF7Selector (context application.InstanceContext) bool {
    return inst.mF7Selector.GetBool(context)
}

//getterForFieldF8Selector
func (inst * comFactory4pComExample1) getterForFieldF8Selector (context application.InstanceContext) float32 {
    return inst.mF8Selector.GetFloat32(context)
}

//getterForFieldF9Selector
func (inst * comFactory4pComExample1) getterForFieldF9Selector (context application.InstanceContext) float64 {
    return inst.mF9Selector.GetFloat64(context)
}

//getterForFieldF10Selector
func (inst * comFactory4pComExample1) getterForFieldF10Selector (context application.InstanceContext) *strings0x1877f3.Builder {

	o1 := inst.mF10Selector.GetOne(context)
	o2, ok := o1.(*strings0x1877f3.Builder)
	if !ok {
		eb := &util.ErrorBuilder{}
		eb.Message("bad cast")
		eb.Set("com", "com0-deep0x1fa623.Example1")
		eb.Set("field", "F10")
		eb.Set("type1", "?")
		eb.Set("type2", "*strings0x1877f3.Builder")
		context.HandleError(eb.Create())
		return nil
	}
	return o2
}

//getterForFieldF11Selector
func (inst * comFactory4pComExample1) getterForFieldF11Selector (context application.InstanceContext) []*strings0x1877f3.Builder {
	list1 := inst.mF11Selector.GetList(context)
	list2 := make([]*strings0x1877f3.Builder, 0, len(list1))
	for _, item1 := range list1 {
		item2, ok := item1.(*strings0x1877f3.Builder)
		if ok {
			list2 = append(list2, item2)
		}
	}
	return list2
}



////////////////////////////////////////////////////////////////////////////////

// comFactory4pComExample2 : the factory of component: Example2
type comFactory4pComExample2 struct {

    mPrototype * deep0x1fa623.Example2

	
	mContextSelector config.InjectionSelector
	mPoolSelector config.InjectionSelector

}

func (inst * comFactory4pComExample2) init() application.ComponentFactory {

	
	inst.mContextSelector = config.NewInjectionSelector("context",nil)
	inst.mPoolSelector = config.NewInjectionSelector("pool",nil)


	inst.mPrototype = inst.newObject()
    return inst
}

func (inst * comFactory4pComExample2) newObject() * deep0x1fa623.Example2 {
	return & deep0x1fa623.Example2 {}
}

func (inst * comFactory4pComExample2) castObject(instance application.ComponentInstance) * deep0x1fa623.Example2 {
	return instance.Get().(*deep0x1fa623.Example2)
}

func (inst * comFactory4pComExample2) GetPrototype() lang.Object {
	return inst.mPrototype
}

func (inst * comFactory4pComExample2) NewInstance() application.ComponentInstance {
	return config.SimpleInstance(inst, inst.newObject())
}

func (inst * comFactory4pComExample2) AfterService() application.ComponentAfterService {
	return inst
}

func (inst * comFactory4pComExample2) Init(instance application.ComponentInstance) error {
	return nil
}

func (inst * comFactory4pComExample2) Destroy(instance application.ComponentInstance) error {
	return nil
}

func (inst * comFactory4pComExample2) Inject(instance application.ComponentInstance, context application.InstanceContext) error {
	
	obj := inst.castObject(instance)
	obj.Context = inst.getterForFieldContextSelector(context)
	obj.Pool = inst.getterForFieldPoolSelector(context)
	return context.LastError()
}

//getterForFieldContextSelector
func (inst * comFactory4pComExample2) getterForFieldContextSelector (context application.InstanceContext) application.Context {
    return context.Context()
}

//getterForFieldPoolSelector
func (inst * comFactory4pComExample2) getterForFieldPoolSelector (context application.InstanceContext) lang.ReleasePool {
    return context.Pool()
}



////////////////////////////////////////////////////////////////////////////////

// comFactory4pComExample3 : the factory of component: com2-deep0x1fa623.Example3
type comFactory4pComExample3 struct {

    mPrototype * deep0x1fa623.Example3

	

}

func (inst * comFactory4pComExample3) init() application.ComponentFactory {

	


	inst.mPrototype = inst.newObject()
    return inst
}

func (inst * comFactory4pComExample3) newObject() * deep0x1fa623.Example3 {
	return & deep0x1fa623.Example3 {}
}

func (inst * comFactory4pComExample3) castObject(instance application.ComponentInstance) * deep0x1fa623.Example3 {
	return instance.Get().(*deep0x1fa623.Example3)
}

func (inst * comFactory4pComExample3) GetPrototype() lang.Object {
	return inst.mPrototype
}

func (inst * comFactory4pComExample3) NewInstance() application.ComponentInstance {
	return config.SimpleInstance(inst, inst.newObject())
}

func (inst * comFactory4pComExample3) AfterService() application.ComponentAfterService {
	return inst
}

func (inst * comFactory4pComExample3) Init(instance application.ComponentInstance) error {
	return inst.castObject(instance).Start()
}

func (inst * comFactory4pComExample3) Destroy(instance application.ComponentInstance) error {
	return inst.castObject(instance).Stop()
}

func (inst * comFactory4pComExample3) Inject(instance application.ComponentInstance, context application.InstanceContext) error {
	return nil
}





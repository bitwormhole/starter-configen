// 这个配置文件是由 starter-configen 工具自动生成的。
// 任何时候，都不要手工修改这里面的内容！！！

package demo

import (
	errors "errors"
	strings "strings"

	configen1_ff6423 "github.com/bitwormhole/starter-configen/tools/configen1"
	car_a924e5 "github.com/bitwormhole/starter-configen/tools/configen1/demo/car"
	application "github.com/bitwormhole/starter/application"
	config "github.com/bitwormhole/starter/application/config"
	lang "github.com/bitwormhole/starter/lang"
)

func Config(configbuilder application.ConfigBuilder) error {

	cominfobuilder := &config.ComInfoBuilder{}
	err := errors.New("OK")

	// builder
	cominfobuilder.Reset()
	cominfobuilder.ID("builder").Class("").Scope("").Aliases("")
	cominfobuilder.OnNew(func() lang.Object {
		return &strings.Builder{}
	})
	cominfobuilder.OnInit(func(o lang.Object) error {
		return nil
	})
	cominfobuilder.OnDestroy(func(o lang.Object) error {
		return nil
	})
	cominfobuilder.OnInject(func(o lang.Object, context application.Context) error {
		adapter := &builder{}
		adapter.instance = o.(*strings.Builder)
		// adapter.context = context
		err := adapter.__inject__(context)
		if err != nil {
			return err
		}
		return nil
	})
	err = cominfobuilder.CreateTo(configbuilder)
	if err != nil {
		return err
	}

	// car1
	cominfobuilder.Reset()
	cominfobuilder.ID("abc").Class("body").Scope("singleton").Aliases("")
	cominfobuilder.OnNew(func() lang.Object {
		return &car_a924e5.Body{}
	})
	cominfobuilder.OnInit(func(o lang.Object) error {
		return nil
	})
	cominfobuilder.OnDestroy(func(o lang.Object) error {
		return nil
	})
	cominfobuilder.OnInject(func(o lang.Object, context application.Context) error {
		adapter := &car1{}
		adapter.instance = o.(*car_a924e5.Body)
		// adapter.context = context
		err := adapter.__inject__(context)
		if err != nil {
			return err
		}
		return nil
	})
	err = cominfobuilder.CreateTo(configbuilder)
	if err != nil {
		return err
	}

	// car2
	cominfobuilder.Reset()
	cominfobuilder.ID("car2").Class("").Scope("").Aliases("x y z")
	cominfobuilder.OnNew(func() lang.Object {
		return &car_a924e5.Body{}
	})
	cominfobuilder.OnInit(func(o lang.Object) error {
		return o.(*car_a924e5.Body).Start()
	})
	cominfobuilder.OnDestroy(func(o lang.Object) error {
		return o.(*car_a924e5.Body).Stop()
	})
	cominfobuilder.OnInject(func(o lang.Object, context application.Context) error {
		adapter := &car2{}
		adapter.instance = o.(*car_a924e5.Body)
		// adapter.context = context
		err := adapter.__inject__(context)
		if err != nil {
			return err
		}
		return nil
	})
	err = cominfobuilder.CreateTo(configbuilder)
	if err != nil {
		return err
	}

	// door
	cominfobuilder.Reset()
	cominfobuilder.ID("door199").Class("").Scope("").Aliases("")
	cominfobuilder.OnNew(func() lang.Object {
		return &car_a924e5.Door{}
	})
	cominfobuilder.OnInit(func(o lang.Object) error {
		return nil
	})
	cominfobuilder.OnDestroy(func(o lang.Object) error {
		return nil
	})
	cominfobuilder.OnInject(func(o lang.Object, context application.Context) error {
		adapter := &door{}
		adapter.instance = o.(*car_a924e5.Door)
		// adapter.context = context
		err := adapter.__inject__(context)
		if err != nil {
			return err
		}
		return nil
	})
	err = cominfobuilder.CreateTo(configbuilder)
	if err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// type builder struct

func (inst *builder) __inject__(context application.Context) error {

	// prepare
	instance := inst.instance
	injection, err := context.Injector().OpenInjection(context)
	if err != nil {
		return err
	}
	defer injection.Close()
	if instance == nil {
		return nil
	}

	// from getters

	// to instance

	// invoke custom inject method

	return injection.Close()
}

////////////////////////////////////////////////////////////////////////////////
// type car1 struct

func (inst *car1) __inject__(context application.Context) error {

	// prepare
	instance := inst.instance
	injection, err := context.Injector().OpenInjection(context)
	if err != nil {
		return err
	}
	defer injection.Close()
	if instance == nil {
		return nil
	}

	// from getters
	inst.BackDoor = inst.__get_BackDoor__(injection, "#door3")
	inst.LeftDoor = inst.__get_LeftDoor__(injection, "#door1")
	inst.RightDoor = inst.__get_RightDoor__(injection, "#door2")

	// to instance
	instance.BackDoor = inst.BackDoor
	instance.LeftDoor = inst.LeftDoor
	instance.RightDoor = inst.RightDoor

	// invoke custom inject method

	return injection.Close()
}

func (inst *car1) __get_BackDoor__(injection application.Injection, selector string) *car_a924e5.Door {

	reader := injection.Select(selector)
	defer reader.Close()

	cnt := reader.Count()
	if cnt != 1 {
		err := errors.New("select.result.count != 1")
		injection.OnError(err)
		return nil
	}

	o1, err := reader.Read()
	if err != nil {
		injection.OnError(err)
		return nil
	}

	o2, ok := o1.(*car_a924e5.Door)
	if !ok {
		err := errors.New("cannot cast component instance to type: *car_a924e5.Door")
		injection.OnError(err)
		return nil
	}

	return o2

}

func (inst *car1) __get_LeftDoor__(injection application.Injection, selector string) *car_a924e5.Door {

	reader := injection.Select(selector)
	defer reader.Close()

	cnt := reader.Count()
	if cnt != 1 {
		err := errors.New("select.result.count != 1")
		injection.OnError(err)
		return nil
	}

	o1, err := reader.Read()
	if err != nil {
		injection.OnError(err)
		return nil
	}

	o2, ok := o1.(*car_a924e5.Door)
	if !ok {
		err := errors.New("cannot cast component instance to type: *car_a924e5.Door")
		injection.OnError(err)
		return nil
	}

	return o2

}

func (inst *car1) __get_RightDoor__(injection application.Injection, selector string) *car_a924e5.Door {

	reader := injection.Select(selector)
	defer reader.Close()

	cnt := reader.Count()
	if cnt != 1 {
		err := errors.New("select.result.count != 1")
		injection.OnError(err)
		return nil
	}

	o1, err := reader.Read()
	if err != nil {
		injection.OnError(err)
		return nil
	}

	o2, ok := o1.(*car_a924e5.Door)
	if !ok {
		err := errors.New("cannot cast component instance to type: *car_a924e5.Door")
		injection.OnError(err)
		return nil
	}

	return o2

}

////////////////////////////////////////////////////////////////////////////////
// type car2 struct

func (inst *car2) __inject__(context application.Context) error {

	// prepare
	instance := inst.instance
	injection, err := context.Injector().OpenInjection(context)
	if err != nil {
		return err
	}
	defer injection.Close()
	if instance == nil {
		return nil
	}

	// from getters
	inst.context2 = inst.__get_context2__(injection, "context")
	inst.list2 = inst.__get_list2__(injection, ".info-list")
	inst.numBool = inst.__get_numBool__(injection, "${demo.num.bool}")
	inst.numFloat32 = inst.__get_numFloat32__(injection, "${demo.num.float32}")
	inst.numFloat64 = inst.__get_numFloat64__(injection, "${demo.num.float64}")
	inst.numInt32 = inst.__get_numInt32__(injection, "${demo.num.int32}")
	inst.numInt64 = inst.__get_numInt64__(injection, "${demo.num.int64}")
	inst.numString = inst.__get_numString__(injection, "${demo.num.string}")

	// to instance

	// invoke custom inject method
	err = inst.gao1xiao(injection)
	if err != nil {
		return err
	}

	return injection.Close()
}

func (inst *car2) __get_context2__(injection application.Injection, selector string) application.Context {
	return injection.Context()
}

func (inst *car2) __get_list2__(injection application.Injection, selector string) []*configen1_ff6423.ComConfigInfo {
	list := make([]*configen1_ff6423.ComConfigInfo, 0)
	reader := injection.Select(selector)
	defer reader.Close()
	for reader.HasMore() {
		o1, err := reader.Read()
		if err != nil {
			injection.OnError(err)
			return list
		}
		o2, ok := o1.(*configen1_ff6423.ComConfigInfo)
		if !ok {
			// err = errors.New("bad cast, selector:" + selector)
			// injection.OnError(err)
			// return list
			// warning ...
			continue
		}
		list = append(list, o2)
	}
	return list

}

func (inst *car2) __get_numBool__(injection application.Injection, selector string) bool {
	reader := injection.Select(selector)
	defer reader.Close()
	value, err := reader.ReadBool()
	if err != nil {
		injection.OnError(err)
	}
	return value
}

func (inst *car2) __get_numFloat32__(injection application.Injection, selector string) float32 {
	reader := injection.Select(selector)
	defer reader.Close()
	value, err := reader.ReadFloat32()
	if err != nil {
		injection.OnError(err)
	}
	return value
}

func (inst *car2) __get_numFloat64__(injection application.Injection, selector string) float64 {
	reader := injection.Select(selector)
	defer reader.Close()
	value, err := reader.ReadFloat64()
	if err != nil {
		injection.OnError(err)
	}
	return value
}

func (inst *car2) __get_numInt32__(injection application.Injection, selector string) int32 {
	reader := injection.Select(selector)
	defer reader.Close()
	value, err := reader.ReadInt32()
	if err != nil {
		injection.OnError(err)
	}
	return value
}

func (inst *car2) __get_numInt64__(injection application.Injection, selector string) int64 {
	reader := injection.Select(selector)
	defer reader.Close()
	value, err := reader.ReadInt64()
	if err != nil {
		injection.OnError(err)
	}
	return value
}

func (inst *car2) __get_numString__(injection application.Injection, selector string) string {
	reader := injection.Select(selector)
	defer reader.Close()
	value, err := reader.ReadString()
	if err != nil {
		injection.OnError(err)
	}
	return value
}

////////////////////////////////////////////////////////////////////////////////
// type door struct

func (inst *door) __inject__(context application.Context) error {

	// prepare
	instance := inst.instance
	injection, err := context.Injector().OpenInjection(context)
	if err != nil {
		return err
	}
	defer injection.Close()
	if instance == nil {
		return nil
	}

	// from getters

	// to instance

	// invoke custom inject method

	return injection.Close()
}

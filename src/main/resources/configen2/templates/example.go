package templates

import (
	"errors"

	"github.com/bitwormhole/starter-tools/tools/configen1/demo/car"
	"github.com/bitwormhole/starter/application"
	"github.com/bitwormhole/starter/application/config"
	"github.com/bitwormhole/starter/lang"
)

type example1car struct {
	instance *car.Body
	context  application.Context
}

func (inst *example1car) __inject__(context application.Context) error {

	// prepare
	instance := inst.instance
	injection, err := context.Injector().OpenInjection(context)
	if err != nil {
		return err
	}
	defer injection.Close()

	// b=getter
	// inst.xxx = inst._get__()

	// a=b
	instance.Name = inst.instance.Name

	// xxx
	err = inst.__inject2__(injection)
	if err != nil {
		return err
	}

	return injection.Close()
}

func (inst *example1car) __inject2__(injection application.Injection) error {
	return nil
}

func (inst *example1car) __get_one__(injection application.Injection, selector string) *car.Body {

	reader := injection.Select(selector)
	defer reader.Close()

	cnt := reader.Count()
	if cnt != 1 {
		err := errors.New("select.result.count != 1")
		injection.OnError(err)
		return nil
	}

	o, err := reader.Read()
	if err != nil {
		injection.OnError(err)
		return nil
	}

	o2, ok := o.(*car.Body)
	if !ok {
		err := errors.New("cannot cast component instance to type: {{}}")
		injection.OnError(err)
		return nil
	}

	return o2
}

func (inst *example1car) __get_list__(injection application.Injection, selector string) []*car.Body {
	// for list
	list := make([]*car.Body, 0)
	reader := injection.Select(selector)
	defer reader.Close()
	for reader.HasMore() {
		o1, err := reader.Read()
		if err != nil {
			injection.OnError(err)
			return list
		}
		o2, ok := o1.(*car.Body)
		if !ok {
			err = errors.New("bad cast, selector:" + selector)
			injection.OnError(err)
			return list
		}
		list = append(list, o2)
	}
	return list
}

func (inst *example1car) __get_map__(injection application.Injection, selector string) map[string]*car.Body {

	return nil
}

func (inst *example1car) __get_property__(injection application.Injection, selector string) string {
	reader := injection.Select(selector)
	defer reader.Close()
	value, err := reader.ReadString()
	if err != nil {
		injection.OnError(err)
	}
	return value
}

func (inst *example1car) __get_context__(injection application.Injection, selector string) application.Context {
	return injection.Context()
}

////////////////////////////////////

func Config(cb application.ConfigBuilder) error {

	cominfobuilder := &config.ComInfoBuilder{}

	// example1
	cominfobuilder.Reset()
	cominfobuilder.ID("").Class("").Scope("").Aliases("")
	cominfobuilder.OnNew(func() lang.Object {
		return &car.Body{}
	})
	cominfobuilder.OnInit(func(o lang.Object) error {
		return o.(*car.Body).Start()
	})
	cominfobuilder.OnDestroy(func(o lang.Object) error {
		return nil
	})
	cominfobuilder.OnInject(func(obj lang.Object, context application.Context) error {
		adapter := &example1car{}
		adapter.instance = obj.(*car.Body)
		adapter.context = context
		err := adapter.__inject__(context)
		if err != nil {
			return err
		}
		return nil
	})
	cominfobuilder.CreateTo(cb)

	return nil
}

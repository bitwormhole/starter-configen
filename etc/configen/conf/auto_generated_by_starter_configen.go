// 这个配置文件是由 starter-configen 工具自动生成的。
// 任何时候，都不要手工修改这里面的内容！！！

package conf

import(
	errors "errors"
	templates_588f8b "github.com/bitwormhole/starter-configen/tools/configen2/templates"
	application "github.com/bitwormhole/starter/application"
	config "github.com/bitwormhole/starter/application/config"
	lang "github.com/bitwormhole/starter/lang"
)


func autoGenConfig(configbuilder application.ConfigBuilder) error {

	cominfobuilder := &config.ComInfoBuilder{}
	err := errors.New("OK")

    
	// theConfigenBootstrap
	cominfobuilder.Reset()
	cominfobuilder.ID("theConfigenBootstrap").Class("").Scope("").Aliases("")
	cominfobuilder.OnNew(func() lang.Object {
		return &ConfigenBootstrap{}
	})
	cominfobuilder.OnInit(func(o lang.Object) error {
		return o.(*ConfigenBootstrap).Start()
	})
	cominfobuilder.OnDestroy(func(o lang.Object) error {
		return nil
	})
	cominfobuilder.OnInject(func(o lang.Object, context application.Context) error {
		adapter := &theConfigenBootstrap{}
		adapter.instance = o.(*ConfigenBootstrap)
		// adapter.context = context
		err := adapter.__inject__(context)
		if err != nil {
			return err
		}
		return nil
	})
	err = cominfobuilder.CreateTo(configbuilder)
    if err !=nil{
        return err
    }

	// theMainTemplateFactory
	cominfobuilder.Reset()
	cominfobuilder.ID("configen2-main-template-factory").Class("").Scope("").Aliases("")
	cominfobuilder.OnNew(func() lang.Object {
		return &templates_588f8b.MainTemplateFactory{}
	})
	cominfobuilder.OnInit(func(o lang.Object) error {
		return nil
	})
	cominfobuilder.OnDestroy(func(o lang.Object) error {
		return nil
	})
	cominfobuilder.OnInject(func(o lang.Object, context application.Context) error {
		adapter := &theMainTemplateFactory{}
		adapter.instance = o.(*templates_588f8b.MainTemplateFactory)
		// adapter.context = context
		err := adapter.__inject__(context)
		if err != nil {
			return err
		}
		return nil
	})
	err = cominfobuilder.CreateTo(configbuilder)
    if err !=nil{
        return err
    }


	return nil
}


////////////////////////////////////////////////////////////////////////////////
// type theConfigenBootstrap struct

func (inst *theConfigenBootstrap) __inject__(context application.Context) error {

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
	inst.AppContext=inst.__get_AppContext__(injection, "context")


	// to instance
	instance.AppContext=inst.AppContext


	// invoke custom inject method


	return injection.Close()
}

func (inst * theConfigenBootstrap) __get_AppContext__(injection application.Injection,selector string) application.Context {
	return injection.Context()
}

////////////////////////////////////////////////////////////////////////////////
// type theMainTemplateFactory struct

func (inst *theMainTemplateFactory) __inject__(context application.Context) error {

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

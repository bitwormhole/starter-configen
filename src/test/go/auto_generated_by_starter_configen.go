// 这个配置文件是由 starter-configen 工具自动生成的。
// 任何时候，都不要手工修改这里面的内容！！！

package srctestgo

import(
	errors "errors"
	application "github.com/bitwormhole/starter/application"
	config "github.com/bitwormhole/starter/application/config"
	lang "github.com/bitwormhole/starter/lang"
)


func autoGenConfig(configbuilder application.ConfigBuilder) error {

	cominfobuilder := &config.ComInfoBuilder{}
	err := errors.New("OK")

    
	// theExample
	cominfobuilder.Reset()
	cominfobuilder.ID("theExample").Class("").Scope("").Aliases("")
	cominfobuilder.OnNew(func() lang.Object {
		return &example1{}
	})
	cominfobuilder.OnInit(func(o lang.Object) error {
		return nil
	})
	cominfobuilder.OnDestroy(func(o lang.Object) error {
		return nil
	})
	cominfobuilder.OnInject(func(o lang.Object, context application.Context) error {
		adapter := &theExample{}
		adapter.instance = o.(*example1)
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
// type theExample struct

func (inst *theExample) __inject__(context application.Context) error {

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
	inst.Ctx=inst.__get_Ctx__(injection, "context")
	inst.Num=inst.__get_Num__(injection, "${test.num}")


	// to instance
	instance.Ctx=inst.Ctx
	instance.Num=inst.Num


	// invoke custom inject method


	return injection.Close()
}

func (inst * theExample) __get_Ctx__(injection application.Injection,selector string) application.Context {
	return injection.Context()
}

func (inst * theExample) __get_Num__(injection application.Injection,selector string) int64 {
	reader := injection.Select(selector)
	defer reader.Close()
	value, err := reader.ReadInt64()
	if err != nil {
		injection.OnError(err)
	}
	return value
}


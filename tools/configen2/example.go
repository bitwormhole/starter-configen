package configen2

import (
	"github.com/bitwormhole/starter/application"
	"github.com/bitwormhole/starter/lang"
	"github.com/bitwormhole/starter/markup"
)

type ComExample struct {
	Name    string
	Pool    lang.Disposable ` a:"x" b:"y" `
	Context application.Context
}

type Com1 struct {
	markup.Component `id:"foo" class:"bar"`

	// 约定：'instance' 固定表示组件对象
	instance *ComExample `injectMethod:"Com1injector"`

	// 约定：'context'固定表示上下文对象
	context application.Context ``

	// 约定：小写开头的字段为间接注入字段
	foo float64 `inject:"${foo}"`

	// 约定：大写开头的字段为直接（自动）注入字段
	Pool lang.Disposable `inject:"#pool"`
	Name string          `inject:"${abc}"`
}

func (inst *Com1) Com1injector(ctx application.Context) error {
	// 这个方法是可选的，执行一些用户定义的注入行为
	return nil
}

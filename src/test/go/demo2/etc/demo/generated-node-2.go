// 这个配置文件是由 starter-configen 工具自动生成的。
// 任何时候，都不要手工修改这里面的内容！！！

package demo

import (
	deep0x1fa623 "github.com/bitwormhole/starter-configen/src/test/go/demo2/code/deep"
	application0x67f6c5 "github.com/bitwormhole/starter/application"
	collection0xee69f0 "github.com/bitwormhole/starter/collection"
	lang0xbf4f1f "github.com/bitwormhole/starter/lang"
	markup0x23084a "github.com/bitwormhole/starter/markup"
	strings0x1877f3 "strings"
)

type pComExample1 struct {
	instance *deep0x1fa623.Example1
	 markup0x23084a.Component `initMethod:"Start"`
	F0 *collection0xee69f0.Properties ``
	F1 application0x67f6c5.Context `inject:"context"`
	F2 lang0xbf4f1f.ReleasePool `inject:"pool"`
	F3 string `inject:"${test.str.s1}"`
	F4 string `inject:"hello,world"`
	F5 int `inject:"1000"`
	F6i8 int8 `inject:"${test.num.i64}"`
	F6i16 int16 `inject:"${test.num.i64}"`
	F6i32 int32 `inject:"${test.num.i64}"`
	F6i64 int64 `inject:"${test.num.i64}"`
	F7 bool `inject:"false"`
	F8 float32 `inject:"${test.num.f32}"`
	F9 float64 `inject:"0.001"`
	F10 *strings0x1877f3.Builder `inject:"*"`
	F11 []*strings0x1877f3.Builder `inject:"*"`
}


type pComExample2 struct {
	instance *deep0x1fa623.Example2
	 markup0x23084a.Controller `id:"Example2" class:"Example"`
	Context application0x67f6c5.Context `inject:"context"`
	Pool lang0xbf4f1f.ReleasePool `inject:"pool"`
}


type pComExample3 struct {
	instance *deep0x1fa623.Example3
	 markup0x23084a.Controller `class:"example demo element" scope:"singleton" aliases:"x y z" initMethod:"Start" destroyMethod:"Stop"`
}


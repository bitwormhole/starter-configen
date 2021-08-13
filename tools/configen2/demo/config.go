package demo

import (
	xstr "strings"

	"github.com/bitwormhole/starter-configen/tools/configen1"
	"github.com/bitwormhole/starter-configen/tools/configen1/demo/car"
	"github.com/bitwormhole/starter/application"
	"github.com/bitwormhole/starter/collection"
	"github.com/bitwormhole/starter/markup"
)

type car1 struct {
	markup.Component

	instance *car.Body `as:"instance" id:"abc"    class:"body" scope:"singleton" `

	LeftDoor  *car.Door `inject:"#door1" to:"LeftDoor"`
	RightDoor *car.Door `inject:"#door2"`
	BackDoor  *car.Door `inject:"#door3"`
}

type car2 struct {
	markup.Component `aliases:"x y z"  injectMethod:"gao1xiao" `

	instance *car.Body `aliases2:"i j k"  initMethod:"Start" destroyMethod:"Stop"`

	context2 application.Context `inject:"context"`

	Table1 map[string]int
	Table2 map[int]*configen1.DomInjection

	List1 []application.Looper
	list2 []*configen1.ComConfigInfo `inject:".info-list"`

	numString  string  `inject:"${demo.num.string}"`
	numInt32   int32   `inject:"${demo.num.int32}"`
	numInt64   int64   `inject:"${demo.num.int64}"`
	numFloat32 float32 `inject:"${demo.num.float32}"`
	numFloat64 float64 `inject:"${demo.num.float64}"`
	numBool    bool    `inject:"${demo.num.bool}"`
}

func (inst *car2) gao1xiao(injection application.Injection) error {

	return nil
}

type door struct {
	markup.Component `id:"door199"`

	instance *car.Door
}

type builder struct {
	markup.Component `id:"" class:""`

	// 小写开头的是特殊字段
	instance *xstr.Builder
	context  application.Context

	num int

	// 大写开头的是直接注入字段
	Env collection.Environment
}

func (inst *builder) inject() error {

	return nil
}

package deep

import (
	"strings"

	"github.com/bitwormhole/starter/application"
	"github.com/bitwormhole/starter/application/config"
	coll "github.com/bitwormhole/starter/collection"
	"github.com/bitwormhole/starter/lang"
	"github.com/bitwormhole/starter/markup"
)

type Example1 struct {
	markup.Component `initMethod:"Start"`

	F0 *coll.Properties
	F1 application.Context `inject:"context"`
	F2 lang.ReleasePool    `inject:"pool"`

	F3    string  `inject:"${test.str.s1}"`
	F4    string  `inject:"hello,world"`
	F5    int     `inject:"1000"`
	F6i8  int8    `inject:"${test.num.i64}"`
	F6i16 int16   `inject:"${test.num.i64}"`
	F6i32 int32   `inject:"${test.num.i64}"`
	F6i64 int64   `inject:"${test.num.i64}"`
	F7    bool    `inject:"false"`
	F8    float32 `inject:"${test.num.f32}"`
	F9    float64 `inject:"0.001"`

	F10 *strings.Builder   `inject:"*"`
	F11 []*strings.Builder `inject:"*"`
	//	F8 map[string]*strings.Builder `inject:"*"`
}

func (inst *Example1) Start() error {
	return nil
}

type Example2 struct {
	markup.Controller `id:"Example2" class:"Example"`
	Context           application.Context `inject:"context"`
	Pool              lang.ReleasePool    `inject:"pool"`
}

type Example3 struct {
	markup.Controller `class:"example demo element" scope:"singleton" aliases:"x y z" initMethod:"Start" destroyMethod:"Stop"`
}

func (inst *Example3) Start() error {
	return nil
}

func (inst *Example3) Stop() error {
	return nil
}

////////////////////////////////////////////////////////////////////////////////

type comFactory struct {
	mPrototype lang.Object
}

func (inst *comFactory) init() application.ComponentFactory {
	return inst
}

func (inst *comFactory) newInst() *Example3 {
	return &Example3{}
}

func (inst *comFactory) cast(instance application.ComponentInstance) *Example3 {
	return instance.Get().(*Example3)
}

func (inst *comFactory) GetPrototype() lang.Object {
	pt := inst.mPrototype
	if pt == nil {
		pt = inst.newInst()
		inst.mPrototype = pt
	}
	return pt
}

func (inst *comFactory) NewInstance() application.ComponentInstance {
	return config.SimpleInstance(inst, inst.newInst())
}

func (inst *comFactory) AfterService() application.ComponentAfterService {
	return inst
}

func (inst *comFactory) Init(instance application.ComponentInstance) error {
	return inst.cast(instance).Start()
}

func (inst *comFactory) Destroy(instance application.ComponentInstance) error {
	return inst.cast(instance).Stop()
}

func (inst *comFactory) Inject(instance application.ComponentInstance, context application.InstanceContext) error {

	return nil
}

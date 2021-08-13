package srctestgo

import (
	"github.com/bitwormhole/starter/application"
	"github.com/bitwormhole/starter/markup"
)

func Config(cb application.ConfigBuilder) error {
	return autoGenConfig(cb)
}

type theExample struct {
	markup.Component
	instance *example1
	Ctx      application.Context `inject:"context"`
	Num      int64               `inject:"${test.num}"`
}

type example1 struct {
	Ctx application.Context
	Num int64
}

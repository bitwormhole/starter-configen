package gen1

import (
	"github.com/bitwormhole/starter-configen/configen3/data/model"
	"github.com/bitwormhole/starter/collection"
)

// Template Gen1-模板
type Template interface {
	Build(source *model.GoSource) (string, error)
}

// NewTemplate 创建 Gen1-模板
func NewTemplate(res collection.Resources, basepath string) Template {
	ctx := &gen1Context{}
	ctx.resources = res
	ctx.resBasePath = basepath
	ctx.templateFactory = &templateFactoryImpl{}
	mt := ctx.templateFactory.createMainTemplate(ctx)
	return mt
}

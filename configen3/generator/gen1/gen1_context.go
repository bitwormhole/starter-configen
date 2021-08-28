package gen1

import (
	"text/template"

	"github.com/bitwormhole/starter-configen/configen3/generator/helper"
	"github.com/bitwormhole/starter/collection"
)

type gen1Context struct {
	resources       collection.Resources
	resBasePath     string
	templateFactory templateFactory
	importsBuilder  helper.ImportListBuilder
}

func (o *gen1Context) loadTemplate(path string) *template.Template {
	temp, err := o.innerLoadTemplate(path)
	if err != nil {
		panic(err)
	}
	return temp
}

func (o *gen1Context) innerLoadTemplate(path string) (*template.Template, error) {
	text, err := o.resources.GetText(o.resBasePath + "/" + path)
	if err != nil {
		return nil, err
	}
	return template.New(path).Parse(text)
}

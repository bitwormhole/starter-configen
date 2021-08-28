package gen2

import (
	"text/template"

	"github.com/bitwormhole/starter-configen/configen3/generator/helper"
	"github.com/bitwormhole/starter/collection"
)

type gen2Context struct {
	resources       collection.Resources
	resBasePath     string
	templateFactory templateFactory
	importsBuilder  helper.ImportListBuilder
}

func (o *gen2Context) loadTemplate(path string, fm template.FuncMap) *template.Template {
	temp, err := o.innerLoadTemplate(path, fm)
	if err != nil {
		panic(err)
	}
	return temp
}

func (o *gen2Context) innerLoadTemplate(path string, fm template.FuncMap) (*template.Template, error) {
	text, err := o.resources.GetText(o.resBasePath + "/" + path)
	if err != nil {
		return nil, err
	}
	tmp := template.New(path)
	if fm != nil {
		tmp.Funcs(fm)
	}
	return tmp.Parse(text)
}

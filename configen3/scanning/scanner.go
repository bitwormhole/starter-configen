package scanning

import (
	"github.com/bitwormhole/starter-configen/configen3/configen"
	"github.com/bitwormhole/starter-configen/configen3/data"
	"github.com/bitwormhole/starter-configen/configen3/data/model"
	"github.com/bitwormhole/starter-configen/configen3/generator/gen1"
	"github.com/bitwormhole/starter/io/fs"
	"github.com/bitwormhole/starter/vlog"
)

type Scanner struct {
	context *configen.Context
	mixdown model.GoSource
}

func (inst *Scanner) Init(ctx *configen.Context, target *model.Scan) data.Scanning {

	inst.mixdown.StructList = make([]*model.Struct, 0)
	inst.context = ctx

	return inst
}

func (inst *Scanner) AddSource(s *model.GoSource) {
	src := s.StructList
	dst := inst.mixdown.StructList
	if src == nil {
		return
	}
	for _, item := range src {
		if item == nil {
			continue
		}
		if !item.IsComponentDefine() {
			continue
		}
		dst = append(dst, item)
	}
	inst.mixdown.StructList = dst
}

func (inst *Scanner) build(output fs.Path) (string, error) {
	alias := output.Parent().Name()
	res := inst.context.Resources
	temp := gen1.NewTemplate(res, "/configen3/gen1/templates/")
	source := &inst.mixdown
	source.Package.Alias = alias
	return temp.Build(source)
}

func (inst *Scanner) WriteResultTo(file fs.Path) error {
	text, err := inst.build(file)
	if err != nil {
		return err
	}
	vlog.Info("write scanning result to file ", file)
	return file.GetIO().WriteText(text, nil, false)
}

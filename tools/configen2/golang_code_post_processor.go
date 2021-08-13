package configen2

import (
	"errors"
	"sort"
	"strings"
)

type golangCodePostProcessor struct {
	context *Context
}

func (inst *golangCodePostProcessor) Init(ctx *Context) {
	inst.context = ctx
}

func (inst *golangCodePostProcessor) isAliasInPart2(p2 string, alias string) bool {
	return strings.Contains(p2, alias)
}

func (inst *golangCodePostProcessor) buildImportItems(path2alias map[string]string, part2code string, buffer *strings.Builder) error {

	keys := make([]string, 0)
	for path := range path2alias {
		keys = append(keys, path)
	}
	sort.Strings(keys)

	buffer.WriteString("import(\n")
	for index := range keys {
		path := keys[index]
		alias := path2alias[path]
		if inst.isAliasInPart2(part2code, alias) {
			buffer.WriteString("\t")
			buffer.WriteString(alias)
			buffer.WriteString(" \"")
			buffer.WriteString(path)
			buffer.WriteString("\"\n")
		}
	}
	buffer.WriteString(")\n")

	return nil
}

func (inst *golangCodePostProcessor) ReImportPackages() error {

	const prefix = "<import>"
	const suffix = "</import>"
	dom := &inst.context.DOM2
	code := inst.context.OutputSourceCode
	builder := &strings.Builder{}

	i1 := strings.Index(code, prefix)
	i2 := strings.Index(code, suffix)
	len2 := len(suffix)
	if 0 < i1 && i1 < i2 {
		// continue
	} else {
		// skip
		return errors.New("no tag-pair: " + prefix + suffix)
	}

	part1 := code[0:i1]
	part2 := code[i2+len2:]

	// build import items
	builder.WriteString(part1)
	inst.buildImportItems(dom.Imports, part2, builder)
	builder.WriteString(part2)
	code = builder.String()

	inst.context.OutputSourceCode = code
	return nil
}

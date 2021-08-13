package configen2

import (
	"errors"
	"strings"

	"github.com/bitwormhole/starter/application"
)

type CodeTemplate interface {
	Build(ctx *BuildingContext) (string, error)
}

type CodeTemplateFactory interface {
	Create(ctx application.Context) (CodeTemplate, error)
}

////////////////////////////////////////////////////////////////////////////////

type BuildingContext struct {
	DOM       *Dom2root
	Component *Dom2component
	Injection *Dom2injection
}

func (inst *BuildingContext) NewChild() *BuildingContext {
	child := &BuildingContext{}
	child.DOM = inst.DOM
	child.Component = inst.Component
	child.Injection = inst.Injection
	return child
}

////////////////////////////////////////////////////////////////////////////////

type TemplateCodeBuilderBase struct {
	Context application.Context

	templateUri  string
	templateText string

	computed map[string]func(ctx *BuildingContext) (string, error)
	children map[string]CodeTemplate
}

func (inst *TemplateCodeBuilderBase) Init(ctx application.Context) error {
	inst.Context = ctx
	inst.children = make(map[string]CodeTemplate)
	inst.computed = make(map[string]func(ctx *BuildingContext) (string, error))
	return nil
}

func (inst *TemplateCodeBuilderBase) AddChild(name string, child CodeTemplate) {
	inst.children[name] = child
}

func (inst *TemplateCodeBuilderBase) AddComputed(name string, fn func(ctx *BuildingContext) (string, error)) {
	inst.computed[name] = fn
}

func (inst *TemplateCodeBuilderBase) LoadTemplate(ctx application.Context, path string) error {
	res := ctx.GetResources()
	text, err := res.GetText(path)
	if err != nil {
		return err
	}
	inst.templateText = text
	inst.templateUri = path
	return nil
}

func (inst *TemplateCodeBuilderBase) BuildFromTemplate(ctx *BuildingContext) (string, error) {

	const beginning = "{{"
	const ending = "}}"
	builder := &strings.Builder{}
	text := inst.templateText

	for {
		// part1
		token := beginning
		i1 := strings.Index(text, token)
		if i1 < 0 {
			builder.WriteString(text)
			break
		} else {
			p1 := text[0:i1]
			p2 := text[i1+len(token):]
			builder.WriteString(p1)
			text = p2
		}
		// part2
		token = ending
		i2 := strings.Index(text, token)
		if i2 < 0 {
			return "", errors.New("no ending token: " + token)
		} else {
			p1 := text[0:i2]
			p2 := text[i2+len(token):]
			p1, err := inst.resolve(p1, ctx)
			if err != nil {
				return "", err
			}
			builder.WriteString(p1)
			text = p2
		}
	}

	text = builder.String()
	return text, nil
}

func (inst *TemplateCodeBuilderBase) resolve(name string, ctx *BuildingContext) (string, error) {

	name = strings.TrimSpace(name)
	fn := inst.computed[name]
	child := inst.children[name]

	if child != nil && fn != nil {
		return "", errors.New("the token.name is duplicated:" + name)
	}

	if child != nil {
		return child.Build(ctx)
	} else if fn != nil {
		return fn(ctx)
	}

	return "", errors.New("no handler for token named: " + name)
}

////////////////////////////////////////////////////////////////////////////////

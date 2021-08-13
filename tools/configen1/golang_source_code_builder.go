package configen1

import (
	"crypto/sha1"
	"errors"
	"sort"
	"strings"

	"github.com/bitwormhole/starter/application"
	"github.com/bitwormhole/starter/application/config"
	"github.com/bitwormhole/starter/lang"
)

type SourceCodeBuilder interface {
	Begin(packageName string)
	AppendComponent(com *ComConfigInfo)
	End()
	Create() (string, error)
}

func NewSourceCodeBuilder() SourceCodeBuilder {
	return &golangSourceCodeBuilder{
		importers: make(map[string]string),
	}
}

type golangSourceCodeBuilder struct {
	buffer        strings.Builder
	packageName   string
	importers     map[string]string // map[path] tag
	scopeMappings map[string]string // map[input]output
}

func (inst *golangSourceCodeBuilder) Config(cb application.ConfigBuilder) error {

	// 这只是一个示例，仅供参考，不要调用

	cb.AddComponent(&config.ComInfo{
		ID:      "",
		Class:   "",
		Scope:   application.ScopeSingleton,
		Aliases: []string{},

		OnNew: func() lang.Object {
			return &lang.TryChain{}
		},

		OnInject: func(obj lang.Object, context application.Context) error {
			target := obj.(*lang.TryChain)
			// inject(target , context)
			return target.Result()
		},

		OnInit: func(obj lang.Object) error {
			target := obj.(*lang.TryChain)
			return target.Result()
		},

		OnDestroy: func(obj lang.Object) error {
			target := obj.(*lang.TryChain)
			return target.Result()
		},
	})

	return nil
}

func (inst *golangSourceCodeBuilder) makePackageNamePrefix(packagePath string, hash bool) string {

	index := strings.LastIndex(packagePath, "/")
	prefix := strings.TrimSpace(packagePath[index+1:])

	if hash {
		builder := &strings.Builder{}
		data := []byte(packagePath)
		sum := sha1.Sum(data)
		chs := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
		for i := 0; i < 8; i++ {
			n := sum[i]
			builder.WriteRune(chs[n&0x0f])
		}
		prefix += "_" + builder.String()
	}

	return prefix
}

func (inst *golangSourceCodeBuilder) getPackageNamePrefix(packagePath string, create bool, hash bool) string {
	prefix := inst.importers[packagePath]
	if create && prefix == "" {
		prefix = inst.makePackageNamePrefix(packagePath, hash)
		inst.importers[packagePath] = prefix
	}
	return prefix
}

func (inst *golangSourceCodeBuilder) Begin(packageName string) {

	inst.packageName = packageName

	// 这里是默认导入的包
	const create = true
	const hash = false

	inst.getPackageNamePrefix("github.com/bitwormhole/starter/application", create, hash)
	inst.getPackageNamePrefix("github.com/bitwormhole/starter/application/config", create, hash)
	inst.getPackageNamePrefix("github.com/bitwormhole/starter/lang", create, hash)

	// func Config
	builder := &inst.buffer
	builder.WriteString("func Config(cb application.ConfigBuilder) error {")
	builder.WriteString("\n\n")
}

func (inst *golangSourceCodeBuilder) AppendComponent(com *ComConfigInfo) {

	if com.Enable == "" {
		return
	}

	builder := &inst.buffer
	builder.WriteString("    // ")
	builder.WriteString(com.InjectionFuncName)
	builder.WriteString("\n")
	builder.WriteString("    cb.AddComponent(&config.ComInfo{\n")

	inst.buildID(com, builder)
	inst.buildClass(com, builder)
	inst.buildScope(com, builder)
	inst.buildAliases(com, builder)

	inst.buildOnNew(com, builder)
	inst.buildOnInit(com, builder)
	inst.buildOnDestroy(com, builder)
	inst.buildOnInject(com, builder)

	builder.WriteString("    })\n")
	builder.WriteString("\n")
}

func (inst *golangSourceCodeBuilder) End() {

	// func Config
	builder := &inst.buffer
	builder.WriteString("    return nil\n")
	builder.WriteString("}\n")
	builder.WriteString("\n")
}

func (inst *golangSourceCodeBuilder) buildID(com *ComConfigInfo, builder *strings.Builder) error {
	id := com.ID
	if id == "" {
		id = "$component#" + com.InjectionFuncName
	}
	builder.WriteString("\t\tID: ")
	builder.WriteString(inst.wrapString(id))
	builder.WriteString(",\n")
	return nil
}

func (inst *golangSourceCodeBuilder) buildClass(com *ComConfigInfo, builder *strings.Builder) error {
	class := com.Class
	builder.WriteString("\t\tClass: ")
	builder.WriteString(inst.wrapString(class))
	builder.WriteString(",\n")
	return nil
}

func (inst *golangSourceCodeBuilder) buildAliases(com *ComConfigInfo, builder *strings.Builder) error {

	const sep1 = ","
	const space = " "

	text := com.Aliases
	text = strings.ReplaceAll(text, "\t", sep1)
	text = strings.ReplaceAll(text, space, sep1)
	items := strings.Split(text, sep1)
	sep := ""

	builder.WriteString("\t\tAliases: []string{")

	for index := range items {
		item := items[index]
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		builder.WriteString(sep)
		builder.WriteString(inst.wrapString(item))
		sep = sep1
	}

	builder.WriteString("},\n")
	return nil
}

func (inst *golangSourceCodeBuilder) buildScope(com *ComConfigInfo, builder *strings.Builder) error {

	mappings := inst.getScopeMappings()
	scope1 := strings.ToLower(com.Scope)
	scope2 := mappings[scope1]

	if scope2 == "" {
		scope2 = ".ScopeSingleton"
	}

	builder.WriteString("\t\tScope: ")
	builder.WriteString("application")
	builder.WriteString(scope2)
	builder.WriteString(",\n")
	return nil
}

func (inst *golangSourceCodeBuilder) buildOnNew(com *ComConfigInfo, builder *strings.Builder) error {

	const tab = "\t\t"
	const nl = "\n"
	const create = true
	const hash = true

	packagePath := com.TargetTypePackagePath
	packageToken := inst.getPackageNamePrefix(packagePath, create, hash)
	simpleName := com.TargetTypeSimpleName

	builder.WriteString(tab + "OnNew: func() lang.Object {" + nl)
	builder.WriteString(tab + "    return &" + packageToken + "." + simpleName + "{}" + nl)
	builder.WriteString(tab + "}," + nl)
	return nil
}

func (inst *golangSourceCodeBuilder) buildOnInit(com *ComConfigInfo, builder *strings.Builder) error {

	const tab = "\t\t"
	const nl = "\n"
	const create = false
	const hash = true

	packagePath := com.TargetTypePackagePath
	packageToken := inst.getPackageNamePrefix(packagePath, create, hash)
	simpleName := com.TargetTypeSimpleName
	funcName := com.InitMethod

	if funcName == "" {
		return nil
	}

	builder.WriteString(tab + "OnInit: func(obj lang.Object) error {" + nl)
	builder.WriteString(tab + "    target := obj.(*" + packageToken + "." + simpleName + ")" + nl)
	builder.WriteString(tab + "    return target." + funcName + "()" + nl)
	builder.WriteString(tab + "}," + nl)
	return nil
}

func (inst *golangSourceCodeBuilder) buildOnDestroy(com *ComConfigInfo, builder *strings.Builder) error {

	const tab = "\t\t"
	const nl = "\n"
	const create = false
	const hash = true

	packagePath := com.TargetTypePackagePath
	packageToken := inst.getPackageNamePrefix(packagePath, create, hash)
	simpleName := com.TargetTypeSimpleName
	funcName := com.DestroyMethod

	if funcName == "" {
		return nil
	}

	builder.WriteString(tab + "OnDestroy: func(obj lang.Object) error {" + nl)
	builder.WriteString(tab + "    target := obj.(*" + packageToken + "." + simpleName + ")" + nl)
	builder.WriteString(tab + "    return target." + funcName + "()" + nl)
	builder.WriteString(tab + "}," + nl)
	return nil
}

func (inst *golangSourceCodeBuilder) buildOnInject(com *ComConfigInfo, builder *strings.Builder) error {

	// todo ...

	const tab = "\t\t"
	const nl = "\n"
	const create = false
	const hash = true

	packagePath := com.TargetTypePackagePath
	packageToken := inst.getPackageNamePrefix(packagePath, create, hash)
	simpleName := com.TargetTypeSimpleName
	funcName := com.InjectionFuncName

	if funcName == "" {
		return nil
	}

	builder.WriteString(tab + "OnInject: func(obj lang.Object,context application.Context) error {" + nl)
	builder.WriteString(tab + "    target := obj.(*" + packageToken + "." + simpleName + ")" + nl)
	builder.WriteString(tab + "    return " + funcName + "(target,context)" + nl)
	builder.WriteString(tab + "}," + nl)
	return nil
}

func (inst *golangSourceCodeBuilder) buildImports(builder *strings.Builder) error {

	src := inst.importers
	if src == nil {
		return errors.New("src==nil")
	}
	builder.WriteString("import(\n")

	keys := make([]string, 0)
	for key := range src {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for index := range keys {
		key := keys[index]
		path := key
		token := src[key]
		builder.WriteString("\t")
		builder.WriteString(token)
		builder.WriteString(" ")
		builder.WriteString(inst.wrapString(path))
		builder.WriteString("\n")
	}

	builder.WriteString(")")
	builder.WriteString("\n\n")
	return nil
}

func (inst *golangSourceCodeBuilder) wrapString(s string) string {
	const mark = "\""
	return mark + s + mark
}

func (inst *golangSourceCodeBuilder) getScopeMappings() map[string]string {
	mappings := inst.scopeMappings
	if mappings == nil {
		mappings = make(map[string]string)

		mappings["singleton"] = ".ScopeSingleton"
		mappings["prototype"] = ".ScopePrototype"

		inst.scopeMappings = mappings
	}
	return mappings
}

func (inst *golangSourceCodeBuilder) Create() (string, error) {

	builder := &inst.buffer
	text := builder.String()
	builder = &strings.Builder{}

	// package
	builder.WriteString("// 这个文件是由 starter-configen 工具生成的配置代码，千万不要手工修改里面的任何内容。\n")
	builder.WriteString("package ")
	builder.WriteString(inst.packageName)
	builder.WriteString("\n\n")

	// import items
	inst.buildImports(builder)

	// func Config(...)
	builder.WriteString(text)

	text = builder.String()
	return text, nil
}

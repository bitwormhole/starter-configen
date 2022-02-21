package configen3

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/bitwormhole/starter-configen/configen3/data/model"
	"github.com/bitwormhole/starter/io/fs"
	"github.com/bitwormhole/starter/markup"
	"github.com/bitwormhole/starter/vlog"
)

// GoSourceLoader 是go源码加载器
type GoSourceLoader struct{}

// Load 加载源码
func (inst *GoSourceLoader) Load(file fs.Path, packagePath string) (*model.GoSource, error) {

	src, err := file.GetIO().ReadText(nil)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, file.Path(), src, 0)
	if err != nil {
		return nil, err
	}

	// Print the AST.
	// ast.Print(fset, f)

	builder := &goSourceBuilder{}
	err = inst.scanAST(fset, f, builder)
	if err != nil {
		return nil, err
	}

	source, err := builder.create()
	if err != nil {
		return nil, err
	}
	source.Package.FullName = packagePath

	binder := &packageNameForTypeBinder{}
	err = binder.bind(source)
	if err != nil {
		return nil, err
	}

	return source, nil
}

func (inst *GoSourceLoader) scanAST(fset *token.FileSet, f *ast.File, builder *goSourceBuilder) error {

	decls := f.Decls
	for _, decl := range decls {
		vlog.Trace("decl: ", decl)
		inst.walkInto(decl, 0, builder)
	}
	return nil
}

func (inst *GoSourceLoader) walkInto(n ast.Node, depth int, builder *goSourceBuilder) error {

	depth++

	aGenDecl, ok := n.(*ast.GenDecl)
	if ok {
		return inst.walkIntoGenDecl(aGenDecl, depth, builder)
	}

	aStructType, ok := n.(*ast.StructType)
	if ok {
		return inst.walkIntoStructType(aStructType, depth, builder)
	}

	aImportSpec, ok := n.(*ast.ImportSpec)
	if ok {
		return inst.walkIntoImportSpec(aImportSpec, depth, builder)
	}

	aTypeSpec, ok := n.(*ast.TypeSpec)
	if ok {
		return inst.walkIntoTypeSpec(aTypeSpec, depth, builder)
	}

	vlog.Warn("unsupported node:", n)

	return nil
}

func (inst *GoSourceLoader) walkIntoTypeSpec(n *ast.TypeSpec, depth int, builder *goSourceBuilder) error {

	name := inst.stringifyIdent(n.Name, nil)
	t := n.Type

	vlog.Trace("TypeSpec name:", name)
	builder.onTypeName(name)

	return inst.walkInto(t, depth, builder)
	//	return nil
}

func (inst *GoSourceLoader) walkIntoImportSpec(n *ast.ImportSpec, depth int, builder *goSourceBuilder) error {
	name := inst.stringifyIdent(n.Name, nil)
	path := inst.stringifyBasicLit(n.Path)
	vlog.Trace("ImportSpec name:", name, "  path:", path)
	builder.addImport(name, path)
	return nil
}

func (inst *GoSourceLoader) walkIntoGenDecl(n *ast.GenDecl, depth int, builder *goSourceBuilder) error {

	for _, spec := range n.Specs {
		err := inst.walkInto(spec, depth, builder)
		if err != nil {
			return err
		}
	}
	return nil
}

func (inst *GoSourceLoader) walkIntoStructType(n *ast.StructType, depth int, builder *goSourceBuilder) error {

	fields := n.Fields
	list := fields.List

	builder.onStructBegin()

	for _, f := range list {
		names := f.Names
		typebuilder := &goTypeBuilder{}
		typeStr := inst.stringifyExpr(f.Type, typebuilder)
		typeModel := typebuilder.create()
		tag := inst.stringifyBasicLit(f.Tag)
		if names == nil {
			vlog.Trace("StructType.field type:", typeStr, " tag:", tag)
			builder.addStructField("", typeModel, tag)
			continue
		} else if len(names) == 0 {
			vlog.Trace("StructType.field type:", typeStr, " tag:", tag)
			builder.addStructField("", typeModel, tag)
			continue
		}
		for _, name := range names {
			nameStr := inst.stringifyIdent(name, typebuilder)
			vlog.Trace("StructType.field type:", typeStr, " tag:", tag, " name:", nameStr)
			builder.addStructField(nameStr, typeModel, tag)
		}
	}

	builder.onStructEnd()
	return nil
}

func (inst *GoSourceLoader) stringifyBasicLit(o *ast.BasicLit) string {
	if o == nil {
		return ""
	}
	return o.Value
}

func (inst *GoSourceLoader) stringifyIdent(o *ast.Ident, typebuilder *goTypeBuilder) string {
	if o == nil {
		return ""
	}
	if typebuilder != nil {
		st := &model.SimpleType{}
		st.TypeName = o.Name
		typebuilder.ValueType = st
	}
	return o.Name
}

func (inst *GoSourceLoader) stringifyExpr(o ast.Expr, typebuilder *goTypeBuilder) string {
	if o == nil {
		return ""
	}

	o2, ok := o.(*ast.SelectorExpr)
	if ok {
		return inst.stringifySelectorExpr(o2, typebuilder)
	}

	o3, ok := o.(*ast.StarExpr)
	if ok {
		return inst.stringifyStarExpr(o3, typebuilder)
	}

	o4, ok := o.(*ast.Ident)
	if ok {
		return inst.stringifyIdent(o4, typebuilder)
	}

	o5, ok := o.(*ast.ArrayType)
	if ok {
		return inst.stringifyArrayType(o5, typebuilder)
	}

	o6, ok := o.(*ast.MapType)
	if ok {
		return inst.stringifyMapType(o6, typebuilder)
	}

	return "???-Expr"
}

func (inst *GoSourceLoader) stringifySelectorExpr(o *ast.SelectorExpr, typebuilder *goTypeBuilder) string {
	if o == nil {
		return ""
	}
	sel := inst.stringifyIdent(o.Sel, typebuilder)
	x := inst.stringifyExpr(o.X, typebuilder)

	st := &model.SimpleType{}
	st.PackageAlias = x
	st.TypeName = sel
	typebuilder.ValueType = st

	return x + "." + sel
}

func (inst *GoSourceLoader) stringifyStarExpr(o *ast.StarExpr, typebuilder *goTypeBuilder) string {
	if o == nil {
		return ""
	}
	str := inst.stringifyExpr(o.X, typebuilder)
	typebuilder.setHasStar(true)
	return "*" + str
}

func (inst *GoSourceLoader) stringifyArrayType(o *ast.ArrayType, typebuilder *goTypeBuilder) string {
	if o == nil {
		return ""
	}
	typebuilder.reset()
	elt := inst.stringifyExpr(o.Elt, typebuilder)
	typebuilder.setIsArray(true)
	return "[]" + elt
}

func (inst *GoSourceLoader) stringifyMapType(o *ast.MapType, typebuilder *goTypeBuilder) string {
	if o == nil {
		return ""
	}

	typebuilder.reset()
	k := inst.stringifyExpr(o.Key, typebuilder)
	myKey := typebuilder.create()

	typebuilder.reset()
	v := inst.stringifyExpr(o.Value, typebuilder)
	myVal := typebuilder.create()

	typebuilder.reset()
	typebuilder.KeyType = myKey.ValueType
	typebuilder.ValueType = myVal.ValueType
	typebuilder.setIsMap(true)

	return "map[" + k + "]" + v
}

////////////////////////////////////////////////////////////////////////////////

type goSourceBuilder struct {
	source          model.GoSource
	currentTypeName string
	currentStruct   *model.Struct
}

func (inst *goSourceBuilder) create() (*model.GoSource, error) {
	gs := &model.GoSource{}
	gs.ImportList = inst.source.ImportList
	gs.StructList = inst.source.StructList
	return gs, nil
}

func (inst *goSourceBuilder) unwrapString(str string, prefix rune, suffix rune) string {
	if suffix == 0 {
		suffix = prefix
	}
	str = strings.TrimSpace(str)
	if strings.HasPrefix(str, string(prefix)) && strings.HasSuffix(str, string(suffix)) {
		return str[1 : len(str)-1]
	}
	return str
}

func (inst *goSourceBuilder) addImport(alias string, path string) {
	path = inst.unwrapString(path, '"', 0)
	o := &model.GoImport{Alias: alias, Path: path}
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash > 0 {
		o.DefaultAlias = path[lastSlash+1:]
	}
	inst.source.ImportList = append(inst.source.ImportList, o)
}

func (inst *goSourceBuilder) onTypeName(name string) {
	inst.currentTypeName = name
}

func (inst *goSourceBuilder) parseTag(tag string) map[string]string {

	const ch1 = '"'
	const ch2 = '`'

	tag = inst.unwrapString(tag, '`', 0)
	tag = strings.ReplaceAll(tag, string(ch1), string(ch2))
	array := strings.Split(tag, string(ch2))
	key := ""
	table := make(map[string]string)

	for _, str := range array {
		str = strings.TrimSpace(str)
		if key == "" {
			if strings.HasSuffix(str, ":") {
				key = str[0 : len(str)-1]
			} else {
				continue
			}
		} else {
			table[key] = str
			key = ""
		}
	}

	return table
}

func (inst *goSourceBuilder) addStructField(name string, xtype *model.ComplexType, tag string) {

	tag = inst.unwrapString(tag, '`', 0)

	f := &model.StructField{}
	f.Name = name
	f.Tag = tag
	f.TagItems = inst.parseTag(tag)

	if xtype != nil {
		f.Type = *xtype
		f.TypeExpr = xtype.String()
	}

	st := inst.currentStruct
	if st == nil {
		return
	}

	if inst.acceptField(f) {
		st.Fields = append(st.Fields, f)
	}
}

func (inst *goSourceBuilder) acceptField(f *model.StructField) bool {
	name := f.Name
	tag := f.TagItems
	if name == "instance" {
		return true
	}
	if name == "" {
		if markup.IsComponentMark(f.Type.ValueType.TypeName) {
			// 这是个例外
			return true
		}
		return false
	}
	if len(tag) < 1 {
		return false
	}
	ch0 := rune(name[0])
	return (('A' <= ch0) && (ch0 <= 'Z'))
}

func (inst *goSourceBuilder) onStructBegin() {
	st := &model.Struct{}
	st.Name = inst.currentTypeName
	st.Type.TypeName = inst.currentTypeName
	inst.currentStruct = st
}

func (inst *goSourceBuilder) onStructEnd() {
	st := inst.currentStruct
	inst.currentStruct = nil
	if st == nil {
		return
	}
	inst.source.StructList = append(inst.source.StructList, st)
}

////////////////////////////////////////////////////////////////////////////////

type goTypeBuilder struct {
	model.ComplexType
}

func (inst *goTypeBuilder) setIsMap(b bool) {
	inst.IsMap = b
}

func (inst *goTypeBuilder) setIsArray(b bool) {
	inst.IsArray = b
}

func (inst *goTypeBuilder) setHasStar(b bool) {
	val := inst.ValueType
	if val == nil {
		val = &model.SimpleType{}
		inst.ValueType = val
	}
	val.HasStar = b
}

func (inst *goTypeBuilder) reset() {
	inst.ComplexType.IsArray = false
	inst.ComplexType.IsMap = false
	inst.ComplexType.KeyType = nil
	inst.ComplexType.ValueType = nil
}

func (inst *goTypeBuilder) create() *model.ComplexType {
	res := &model.ComplexType{}
	*res = inst.ComplexType
	return res
}

////////////////////////////////////////////////////////////////////////////////

type packageNameForTypeBinder struct {
	keyForThisPkg   string
	thisPackagePath string
	packs           map[string]*model.GoImport
}

func (inst *packageNameForTypeBinder) bind(src *model.GoSource) error {

	if src == nil {
		return errors.New("bind:src==nil")
	}

	inst.keyForThisPkg = "[THIS]"
	inst.thisPackagePath = src.Package.FullName

	err := inst.initImports(src)
	if err != nil {
		return err
	}

	err = inst.seekAllTypesInSource(src)
	if err != nil {
		return err
	}

	return nil
}

func (inst *packageNameForTypeBinder) initImports(src *model.GoSource) error {

	list := src.ImportList
	table := make(map[string]*model.GoImport)
	inst.packs = table

	if list == nil {
		return nil
	}

	for _, item := range list {
		table[item.Alias] = item
		table[item.DefaultAlias] = item
		table[item.Path] = item
	}

	// the 'this' package
	thisPkg := &model.GoImport{}
	thisPkg.Alias = ""
	thisPkg.DefaultAlias = ""
	thisPkg.Path = src.Package.FullName
	table[""] = thisPkg
	table[thisPkg.Path] = thisPkg

	return nil
}

func (inst *packageNameForTypeBinder) seekAllTypesInSource(source *model.GoSource) error {

	list := source.StructList
	if list == nil {
		return nil
	}
	for _, item := range list {
		err := inst.handleSimpleType(&item.Type)
		if err != nil {
			return err
		}
		err = inst.seekAllTypesInFields(item.Fields)
		if err != nil {
			return err
		}
	}
	return nil
}

func (inst *packageNameForTypeBinder) seekAllTypesInFields(fields []*model.StructField) error {
	if fields == nil {
		return nil
	}
	for _, f := range fields {
		err := inst.handleComplexType(&f.Type)
		if err != nil {
			return err
		}
	}
	return nil
}

func (inst *packageNameForTypeBinder) handleComplexType(t *model.ComplexType) error {

	if t == nil {
		return nil
	}

	err := inst.handleSimpleType(t.KeyType)
	if err != nil {
		return err
	}

	err = inst.handleSimpleType(t.ValueType)
	if err != nil {
		return err
	}

	return nil
}

func (inst *packageNameForTypeBinder) handleSimpleType(t *model.SimpleType) error {
	if t == nil {
		return nil
	}
	alias := t.PackageAlias
	tName := t.TypeName
	if inst.isExtType(tName) {
		t.PackageName = inst.findPackageName(alias, true)
	}
	return nil
}

func (inst *packageNameForTypeBinder) findPackageName(alias string, includeThis bool) string {
	table := inst.packs
	if table != nil {
		item := table[alias]
		if item != nil {
			return item.Path
		}
	}
	if includeThis && alias == "" {
		return inst.thisPackagePath
	}
	return ""
}

func (inst *packageNameForTypeBinder) isExtType(name string) bool {
	name = strings.TrimSpace(name)
	if len(name) < 1 {
		return false
	}
	ch := name[0]
	return 'A' <= ch && ch <= 'Z'
}

package configen1

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"

	"github.com/bitwormhole/starter/collection"
)

type DomInjection struct {
	Name       string // the function name
	TargetType string // like 'package_path#struct_name'
	Properties map[string]string
}

type domHandler interface {
	Handle(n ast.Node) error
	Flush() error
}

type domWriter interface {
	WritePackageName(name string) error
	WriteInjectionBegin(funcName string) error
	WriteInjectionType(funcName string, targetTypeName string) error
	WriteInjectionProps(funcName string, props map[string]string) error
	WriteInjectionEnd(funcName string) error
}

type domBuilder interface {
	SetPackageName(pkgName string)
	GetPackageName() string
	Error() error
	OnError(err error)
	AppendInjectFunc(fn *DomInjection)
	Create() ([]*DomInjection, error)
}

////////////////////////////////////////////////////////////////////////////////
// struct  injectionPropertiesReader

type injectionPropertiesReader struct {
	groups []*ast.CommentGroup
	ptr    int
	size   int
}

func (inst *injectionPropertiesReader) init(file *ast.File) {
	dst := []*ast.CommentGroup{}
	inst.groups = dst
	if file == nil {
		return
	}
	src := file.Comments
	if src == nil {
		return
	}
	for index := range src {
		item := src[index]
		if item == nil {
			continue
		}
		dst = append(dst, item)
	}
	inst.groups = dst
	inst.size = len(dst)
	inst.ptr = 0
}

func (inst *injectionPropertiesReader) computePosDiff(inner ast.Node, outer ast.Node) (int, error) {

	// inner: moving
	// outer: static
	if inner == nil || outer == nil {
		return 0, errors.New("nil")
	}

	if inner.End() < outer.Pos() {
		return -1, nil
	} else if outer.End() < inner.Pos() {
		return 1, nil
	}
	return 0, nil
}

func (inst *injectionPropertiesReader) readPropertiesTextInBody(body ast.Node) (string, error) {
	var err error = nil
	var text string
	ptr := inst.ptr
	list := inst.groups
	size := inst.size
	for ; ptr < size; ptr++ {
		item := list[ptr]
		diff, _err := inst.computePosDiff(item, body)
		if _err != nil {
			err = _err
			break
		}
		if diff < 0 {
			continue
		} else if diff == 0 {
			text = item.Text() + "\n  enable=true"
			inst.ptr = ptr
			break
		} else {
			break
		}
	}
	return text, err
}

func (inst *injectionPropertiesReader) readPropertiesInBody(n ast.Node) (map[string]string, error) {

	text, err := inst.readPropertiesTextInBody(n)
	if err != nil {
		return nil, err
	}

	props, err := collection.ParseProperties(text, nil)
	if err != nil {
		return nil, err
	}

	table1 := props.Export(nil)
	table2 := make(map[string]string)

	for key1 := range table1 {
		key2 := strings.ToLower(key1)
		table2[key2] = table1[key1]
	}

	return table2, nil
}

////////////////////////////////////////////////////////////////////////////////
// struct defaultDomHandler

type defaultDomHandler struct {
	writer           domWriter
	file             *ast.File
	propertiesReader *injectionPropertiesReader
	importsTable     map[string]string // map[prefix] fullpath

	// state
	currentFuncName string
}

func createDomHandler(writer domWriter, file *ast.File) domHandler {

	reader := &injectionPropertiesReader{}
	reader.init(file)

	return &defaultDomHandler{
		writer:           writer,
		file:             file,
		propertiesReader: reader,
		importsTable:     make(map[string]string),
	}
}

func (inst *defaultDomHandler) Handle(n ast.Node) error {

	if n != nil {
		fmt.Println("-------------------------> pos:", n.Pos(), " end:", n.End())
	}

	t_func_decl, ok := n.(*ast.FuncDecl)
	if ok {
		return inst.handleFuncDecl(t_func_decl)
	}

	t_import_spec, ok := n.(*ast.ImportSpec)
	if ok {
		return inst.handleImportSpec(t_import_spec)
	}

	t_file, ok := n.(*ast.File)
	if ok {
		return inst.handleFile(t_file)
	}

	return nil
}

func (inst *defaultDomHandler) setCurrentFuncName(name string) {
	old := inst.currentFuncName
	inst.currentFuncName = name
	if old != "" {
		inst.writer.WriteInjectionEnd(old)
	}
	if name != "" {
		inst.writer.WriteInjectionBegin(name)
	}
}

func (inst *defaultDomHandler) tryGetTargetType(n *ast.FuncDecl) (string, error) {

	t := n.Type
	if t == nil {
		return "", errors.New("FuncType==nil")
	}

	params := t.Params
	if params == nil {
		return "", errors.New("FuncType.params==nil")
	}

	if params.NumFields() != 2 {
		return "", errors.New("FuncType.params.NumFields!=2")
	}

	field := params.List[0]
	if field == nil {
		return "", errors.New("FuncType.params.field[0]==nil")
	}

	expr, ok := field.Type.(*ast.StarExpr)
	if !ok {
		return "", errors.New("want a ast.StarExpr")
	}

	expr_x := expr.X
	expr_str := fmt.Sprint(expr_x) // like '${prefix Type}'
	expr_str = strings.TrimSpace(expr_str)

	if strings.HasPrefix(expr_str, "&{") && strings.HasSuffix(expr_str, "}") {
		// continue
	} else {
		return "", errors.New("bad expr:" + expr_str)
	}

	i1 := strings.IndexRune(expr_str, '{')
	i2 := strings.LastIndexByte(expr_str, '}')
	expr_str = expr_str[i1+1 : i2]
	i3 := strings.LastIndexByte(expr_str, ' ')
	prefix := strings.TrimSpace(expr_str[0:i3])
	simpleName := strings.TrimSpace(expr_str[i3+1:])

	path := inst.importsTable[prefix]
	result := path + "#" + simpleName
	return result, nil
}

func (inst *defaultDomHandler) handleFuncDecl(n *ast.FuncDecl) error {

	// name
	func_id := n.Name
	if func_id == nil {
		return errors.New("FuncDecl.Name==nil")
	}
	func_name := func_id.Name
	inst.setCurrentFuncName(func_name)

	// target-type
	target_type, err := inst.tryGetTargetType(n)
	if err == nil {
		inst.writer.WriteInjectionType(func_name, target_type)
	}

	// properties
	props, err := inst.propertiesReader.readPropertiesInBody(n.Body)
	if err == nil {
		inst.writer.WriteInjectionProps(func_name, props)
	}

	return nil
}

func (inst *defaultDomHandler) handleFile(n *ast.File) error {
	name := n.Name
	if name != nil {
		inst.writer.WritePackageName(name.Name)
	}
	return nil
}

func (inst *defaultDomHandler) handleImportSpec(n *ast.ImportSpec) error {

	_path := n.Path
	_prefix := n.Name
	if _path == nil {
		return errors.New("no path")
	}

	path := _path.Value
	path = strings.ReplaceAll(path, "\"", "\t")
	path = strings.TrimSpace(path)

	var prefix string
	if _prefix == nil {
		index := strings.LastIndexByte(path, '/')
		prefix = path[index+1:]
	} else {
		prefix = _prefix.String()
	}

	fmt.Println("[ImportSpec from:", path, " as:", prefix, "]")
	inst.importsTable[prefix] = path
	return nil
}

func (inst *defaultDomHandler) Flush() error {
	inst.setCurrentFuncName("")
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// struct defaultDomWriter

type defaultDomWriter struct {
	builder    domBuilder
	currentDOM *DomInjection
}

func createDomWriter(builder domBuilder) domWriter {
	return &defaultDomWriter{
		builder: builder,
	}
}

func (inst *defaultDomWriter) WritePackageName(name string) error {
	fmt.Println("[Package name:", name, "]")
	inst.builder.SetPackageName(name)
	return nil
}

func (inst *defaultDomWriter) WriteInjectionBegin(funcName string) error {

	fmt.Println("[injection.begin name:", funcName, "]")

	f := &DomInjection{}
	f.Name = funcName
	inst.currentDOM = f
	return nil
}

func (inst *defaultDomWriter) WriteInjectionType(funcName string, targetTypeName string) error {

	dom := inst.currentDOM
	if dom == nil {
		return nil
	}
	if funcName != dom.Name {
		return nil
	}

	fmt.Println("[injection.type  name:", funcName, " target:", targetTypeName, "]")
	dom.TargetType = targetTypeName
	return nil
}

func (inst *defaultDomWriter) WriteInjectionProps(funcName string, props map[string]string) error {

	dom := inst.currentDOM
	if dom == nil {
		return nil
	}
	if funcName != dom.Name {
		return nil
	}

	fmt.Println("[injection.props name:", funcName, " props:", props, "]")
	dom.Properties = props
	return nil
}

func (inst *defaultDomWriter) WriteInjectionEnd(funcName string) error {

	dom := inst.currentDOM
	inst.currentDOM = nil

	if dom == nil {
		return nil
	}
	if funcName != dom.Name {
		return nil
	}

	fmt.Println("[injection.end name:", dom.Name, "]")
	inst.builder.AppendInjectFunc(dom)
	return nil
}

func (inst *defaultDomWriter) Flush() error {
	dom := inst.currentDOM
	if dom == nil {
		return nil
	}
	inst.WriteInjectionEnd(dom.Name)
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// struct defaultDomBuilder

type defaultDomBuilder struct {
	injectionTable map[string]*DomInjection
	lastError      error
	packageName    string
}

func createDomBuilder() domBuilder {
	return &defaultDomBuilder{
		injectionTable: make(map[string]*DomInjection),
	}
}

func (inst *defaultDomBuilder) SetPackageName(name string) {
	inst.packageName = name
}

func (inst *defaultDomBuilder) GetPackageName() string {
	return inst.packageName
}

func (inst *defaultDomBuilder) Error() error {
	return inst.lastError
}

func (inst *defaultDomBuilder) OnError(err error) {
	if err == nil {
		return
	}
	inst.lastError = err
}

func (inst *defaultDomBuilder) AppendInjectFunc(fn *DomInjection) {
	if fn == nil {
		return
	}
	inst.injectionTable[fn.Name] = fn
}

func (inst *defaultDomBuilder) Create() ([]*DomInjection, error) {
	err := inst.lastError
	if err != nil {
		return nil, err
	}
	table := inst.injectionTable
	list := []*DomInjection{}
	for key := range table {
		item := table[key]
		list = append(list, item)
	}
	return list, nil
}

////////////////////////////////////////////////////////////////////////////////
// EOF

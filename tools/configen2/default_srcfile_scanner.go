package configen2

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/bitwormhole/starter/io/fs"
)

type defaultSourceFileScanner struct {
	Context *Context
}

func (inst *defaultSourceFileScanner) init(Context *Context) SourceFileScanner {
	return inst
}

func (inst *defaultSourceFileScanner) Scan(file fs.Path) (*Dom1doc, error) {

	fmt.Println("scan source file", file.Path())

	source, err := file.GetIO().ReadText()
	if err != nil {
		return nil, err
	}

	mode := parser.ParseComments
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file.Name(), source, mode)
	if err != nil {
		return nil, err
	}

	//	builder := createDomBuilder()
	//	writer := createDomWriter(builder)
	//	handler := createDomHandler(writer, f)
	context := &scanningContext{}
	handler := &scanningHandler{context: context}
	builder := &dom1builderImpl{context: context}
	unwrapper := &scanningUnwrapper{context: context}

	context.doc = &Dom1doc{}
	context.builder = builder
	context.handler = handler
	context.unwrapper = unwrapper
	context.source = source

	ast.Inspect(f, func(n ast.Node) bool {
		err := handler.HandleNode(n)
		handler.HandleError(err)
		return true
	})

	err = handler.Flush()
	if err != nil {
		return nil, err
	}

	return context.doc, nil
}

////////////////////////////////////////////////////////////////////////////////

type scanningHandler struct {
	context *scanningContext
}

func (inst *scanningHandler) HandleNode(n ast.Node) error {

	if n != nil {
		fmt.Println("-------------------------> pos:", n.Pos(), " end:", n.End())
	}

	t_struct_type, ok := n.(*ast.StructType)
	if ok {
		return inst.handleStructType(t_struct_type)
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

func (inst *scanningHandler) HandleError(err error) {
}

func (inst *scanningHandler) Flush() error {
	return nil
}

func (inst *scanningHandler) handleFile(n *ast.File) error {
	pkgName := n.Name.Name
	fmt.Println("  [package Name:", pkgName, "]")
	inst.context.builder.SetPackageName(pkgName)
	return nil
}

func (inst *scanningHandler) handleImportSpec(n *ast.ImportSpec) error {

	unwrapper := inst.context.unwrapper
	imp := &Dom1import{}
	imp.Alias = unwrapper.unwrapIdent(n.Name)
	imp.Path = unwrapper.unwrapBasicLit(n.Path)

	path := imp.Path
	path = strings.ReplaceAll(path, "\"", "\t")
	imp.Path = strings.TrimSpace(path)

	inst.context.builder.AddImport(imp)
	fmt.Println("  [import name:", imp.Alias, " path:", imp.Path, "]")
	return nil
}

func (inst *scanningHandler) handleStructType(n *ast.StructType) error {

	// fmt.Println("  [struct name:", n.Fields ,  "]")

	fields1 := n.Fields
	fields2 := make([]*Dom1field, 0)
	unwrapper := inst.context.unwrapper

	if fields1 != nil {
		list := fields1.List
		for index := range list {
			f := list[index]
			field := &Dom1field{}

			field.Name = unwrapper.unwrapIdentAt(0, f.Names)
			field.Type = unwrapper.unwrapExpr(f.Type)
			field.Tag = unwrapper.unwrapBasicLit(f.Tag)

			fmt.Println("      [field name:", field.Name, " type:", field.Type, " tag:", field.Tag, "]")
			fields2 = append(fields2, field)
		}
	}

	st := &Dom1struct{}
	st.Name = unwrapper.getStructName(n)
	st.Fields = fields2
	inst.context.builder.AddStruct(st)
	return nil
}

////////////////////////////////////////////////////////////////////////////////

type Dom1builder interface {
	SetPackageName(string)
	AddImport(*Dom1import)
	AddStruct(*Dom1struct)
}

////////////////////////////////////////////////////////////////////////////////

type dom1builderImpl struct {
	context *scanningContext
}

func (inst *dom1builderImpl) _impl_() Dom1builder {
	return inst
}

func (inst *dom1builderImpl) SetPackageName(name string) {
	inst.context.doc.PackageName = name
}

func (inst *dom1builderImpl) AddImport(item *Dom1import) {

	list := inst.context.doc.ImportList
	if list == nil {
		list = make([]*Dom1import, 0)
	}
	if item != nil {
		list = append(list, item)
	}
	inst.context.doc.ImportList = list
}

func (inst *dom1builderImpl) AddStruct(item *Dom1struct) {

	list := inst.context.doc.ComList
	if list == nil {
		list = make([]*Dom1struct, 0)
	}
	if item != nil {
		list = append(list, item)
	}
	inst.context.doc.ComList = list
}

////////////////////////////////////////////////////////////////////////////////

type scanningUnwrapper struct {
	context *scanningContext
}

func (inst *scanningUnwrapper) unwrap() string {
	return ""
}

func (inst *scanningUnwrapper) unwrapIdent(n *ast.Ident) string {
	if n == nil {
		return ""
	}
	return n.Name
}

func (inst *scanningUnwrapper) unwrapBasicLit(n *ast.BasicLit) string {
	if n == nil {
		return ""
	}
	return n.Value
}

func (inst *scanningUnwrapper) unwrapExpr(n ast.Expr) string {
	src := inst.context.source
	from := int(n.Pos()) - 1
	to := int(n.End()) - 1
	size := len(src)
	if 0 < from && from < to && to < size {
		return src[from:to]
	}
	return ""
}

func (inst *scanningUnwrapper) unwrapIdentAt(index int, list []*ast.Ident) string {
	if list == nil {
		return ""
	}
	size := len(list)
	if 0 <= index && index < size {
		return inst.unwrapIdent(list[index])
	}
	return ""
}

func (inst *scanningUnwrapper) getStructName(st *ast.StructType) string {

	const token1 = "type"
	line := ""
	pos := int(st.Pos()) - 1
	src := inst.context.source

	for i := pos; i > 0; i-- {
		b := src[i]
		if b == '\n' || b == '\r' {
			line = src[i:pos]
			break
		}
	}

	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, token1) {
		offset := len(token1)
		line = strings.TrimSpace(line[offset:])
	}
	return line
}

////////////////////////////////////////////////////////////////////////////////

type scanningContext struct {
	unwrapper *scanningUnwrapper
	handler   *scanningHandler
	doc       *Dom1doc
	builder   Dom1builder
	source    string
}

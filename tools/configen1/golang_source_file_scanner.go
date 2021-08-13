package configen1

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/bitwormhole/starter/io/fs"
)

type golangSourceFileScanner struct {
}

type golangSourceFileScanningResult struct {
	injectionList []*DomInjection
	packageName   string
}

func (inst *golangSourceFileScanner) scan(file fs.Path) (*golangSourceFileScanningResult, error) {

	fmt.Println("scan file ", file.Path())

	result := &golangSourceFileScanningResult{}
	source, err := inst.loadSourceCode(file)
	if err != nil {
		return nil, err
	}

	mode := parser.ParseComments
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file.Name(), source, mode)
	if err != nil {
		return nil, err
	}

	builder := createDomBuilder()
	writer := createDomWriter(builder)
	handler := createDomHandler(writer, f)

	ast.Inspect(f, func(n ast.Node) bool {
		err := handler.Handle(n)
		builder.OnError(err)
		return true
	})

	err = handler.Flush()
	if err != nil {
		return nil, err
	}

	dom, err := builder.Create()
	if err != nil {
		return nil, err
	}

	result.packageName = builder.GetPackageName()
	result.injectionList = dom
	return result, nil
}

func (inst *golangSourceFileScanner) loadSourceCode(file fs.Path) (string, error) {
	return file.GetIO().ReadText()
}

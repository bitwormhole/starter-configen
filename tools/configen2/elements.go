package configen2

import "github.com/bitwormhole/starter/io/fs"

type ComponentInfo struct {
	ID      string
	Classes []string
}

type ComponentInfoTable interface {
	Add(com *ComponentInfo)
	All() []*ComponentInfo
}

type DirectoryScanner interface {
	Scan(dir fs.Path) error
}

type SourceFileScanner interface {
	Scan(file fs.Path) (*Dom1doc, error)
}

type SourceCodeBuilder interface {
	Build() (string, error)
}

type Runner interface {
	Run() error
}

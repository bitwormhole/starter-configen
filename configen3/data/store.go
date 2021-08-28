package data

import (
	"github.com/bitwormhole/starter-configen/configen3/data/model"
	"github.com/bitwormhole/starter/io/fs"
)

type Nodes struct {
	Name string
	List []*model.Node
}

type SourceFiles struct {
	Name string
}

type TargetFiles struct {
	Name string
}

////////////////////////////////////////////////////////////////////////////////

type Store struct {
	PWD             fs.Path
	Nodes           Nodes
	RootNode        *model.Node
	SourceFiles     SourceFiles
	TargetFiles     TargetFiles
	CurrentScanning Scanning
}

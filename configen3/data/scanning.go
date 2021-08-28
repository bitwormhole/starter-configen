package data

import (
	"github.com/bitwormhole/starter-configen/configen3/data/model"
	"github.com/bitwormhole/starter/io/fs"
)

type Scanning interface {
	AddSource(item *model.GoSource)
	WriteResultTo(file fs.Path) error
}

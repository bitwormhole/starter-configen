package configen

import (
	"github.com/bitwormhole/starter-configen/configen3/data"
	"github.com/bitwormhole/starter/collection"
)

type Context struct {
	Resources collection.Resources
	Store     data.Store
}

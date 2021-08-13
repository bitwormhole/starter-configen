package configen

import (
	"github.com/bitwormhole/starter"
	"github.com/bitwormhole/starter-configen/etc/configen/conf"
	srcmain "github.com/bitwormhole/starter-configen/src/main"
	"github.com/bitwormhole/starter/application"
)

const (
	myName     = "github.com/bitwormhole/starter-configen"
	myVersion  = "v0.0.3"
	myRevision = 3
)

func Module() application.Module {

	mod := &application.DefineModule{
		Name:     myName,
		Version:  myVersion,
		Revision: myRevision,
	}

	mod.OnMount = func(cb application.ConfigBuilder) error { return conf.ExportConfig(cb, mod) }
	mod.Resources = srcmain.ExportResources()
	mod.AddDependency(starter.Module())

	return mod
}

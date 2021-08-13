package configen2

import (
	"github.com/bitwormhole/starter/application"
	"github.com/bitwormhole/starter/io/fs"
)

type Context struct {
	PWD             fs.Path
	InputFile       fs.Path // the 'configen.properties'
	OutputFile      fs.Path // the 'auto_generated_by_starter_configen.go'
	ConfigenVersion string  // 'configen.version'='v2'
	DryRun          bool
	AppContext      application.Context

	Dom2Builder       Dom2Builder
	CodeBuilder       CodeBuilder
	DirectoryScanner  DirectoryScanner
	SourceFileScanner SourceFileScanner
	ComTable          ComponentInfoTable
	Imports           ImportManager

	DOM1             Dom1root
	DOM2             Dom2root
	OutputSourceCode string
}

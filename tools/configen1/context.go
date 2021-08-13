package configen1

import "github.com/bitwormhole/starter/io/fs"

type configenContext struct {
	PWD               fs.Path
	GoSourceFiles     []fs.Path
	DomInjectionTable map[string]*DomInjection
	ComConfigList     []*ComConfigInfo
	PackageName       string
	OutputCode        string
	OutputFileName    string
	ConfigFileName    string
}

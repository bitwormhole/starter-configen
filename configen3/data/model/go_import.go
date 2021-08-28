package model

// GoImport 代表一个导入语句
type GoImport struct {
	DefaultAlias string
	Alias        string
	Path         string
	NoHash       bool
}

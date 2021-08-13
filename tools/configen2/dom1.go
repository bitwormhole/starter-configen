package configen2

type Dom1root struct {
	PackageName string
	Documents   map[string]*Dom1doc
	Imports     map[string]bool // map[path] yes
}

type Dom1doc struct {
	PackageName string
	ImportList  []*Dom1import
	ComList     []*Dom1struct
}

type Dom1import struct {
	Alias string
	Path  string
}

type Dom1struct struct {
	Name   string
	Fields []*Dom1field
}

type Dom1field struct {
	Name string
	Type string
	Tag  string
}

package model

import (
	"github.com/bitwormhole/starter/collection"
	"github.com/bitwormhole/starter/io/fs"
)

type Child struct {
	Name string
	Path fs.Path
	Href string
}

type Scan struct {
	Name         string
	Path         fs.Path
	WriteTo      fs.Path
	Href         string
	WriteToChild string
	WriteToFile  string
	R            bool
}

type Node struct {
	Version string
	Path    fs.Path
	// Text        string
	Properties  collection.Properties
	Parent      *Node
	Children    []*Child
	ScanTargets []*Scan
}

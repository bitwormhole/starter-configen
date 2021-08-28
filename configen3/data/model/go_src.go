package model

// GoSource 代表一个Go源文件
type GoSource struct {
	Package    PackageRef
	ImportList []*GoImport
	StructList []*Struct
}

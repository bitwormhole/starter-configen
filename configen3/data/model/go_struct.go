package model

import "github.com/bitwormhole/starter/markup"

// Struct 代表一个Go结构体
type Struct struct {
	// OwnerPackage *PackageRef
	Name   string
	Type   SimpleType
	Fields []*StructField
}

// StructRef 是对Struct 的引用
type StructRef struct {
	// markup.Component
}

func (o *Struct) isTypeOfComponentDefine(t *ComplexType) bool {
	if t == nil {
		return false
	}
	kt := t.KeyType
	vt := t.ValueType
	if t.IsMap || t.IsArray || (kt != nil) || (vt == nil) {
		return false
	}
	if vt.HasStar {
		return false
	}
	return markup.IsComponentMarkWithPackage(vt.PackageName, vt.TypeName)
}

// IsComponentDefine 判断结构是否是一个组件定义
func (o *Struct) IsComponentDefine() bool {
	fields := o.Fields
	if fields == nil {
		return false
	}
	for _, f := range fields {
		if f == nil {
			continue
		}
		// 2个必须的条件：
		con1 := (f.Name == "")
		con2 := o.isTypeOfComponentDefine(&f.Type)
		if con1 && con2 {
			return true
		}
	}
	return false
}

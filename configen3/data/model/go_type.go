package model

import "strings"

////////////////////////////////////////////////////////////////////////////////

// SimpleType 表示一个简单的类型（number|struct|interface）
type SimpleType struct {
	TypeName     string // the 'b' of 'a.b'
	PackageName  string // aka. PackagePath
	PackageAlias string // the 'a' of 'a.b'
	HasStar      bool   // the '*'
}

func (inst *SimpleType) Clone() *SimpleType {
	if inst == nil {
		return inst
	}
	other := &SimpleType{}
	other.TypeName = inst.TypeName
	other.PackageName = inst.PackageName
	other.PackageAlias = inst.PackageAlias
	other.HasStar = inst.HasStar
	return other
}

func (inst *SimpleType) isPublicObjectType() bool {
	name := inst.TypeName
	if name == "" {
		return false
	}
	ch0 := rune(name[0])
	return (('A' <= ch0) && (ch0 <= 'Z'))
}

func (inst *SimpleType) String() string {
	builder := strings.Builder{}
	if inst.HasStar {
		builder.WriteString("*")
	}
	if inst.isPublicObjectType() {
		alias := inst.PackageAlias
		builder.WriteString(alias)
		builder.WriteRune('.')
	}
	builder.WriteString(inst.TypeName)
	return builder.String()
}

////////////////////////////////////////////////////////////////////////////////

// ComplexType 表示一个复杂的类型（array|map）
type ComplexType struct {
	KeyType   *SimpleType
	ValueType *SimpleType
	IsMap     bool
	IsArray   bool
}

func (inst *ComplexType) Clone(deep bool) *ComplexType {
	if inst == nil {
		return inst
	}
	other := &ComplexType{}
	other.IsArray = inst.IsArray
	other.IsMap = inst.IsMap
	if deep {
		other.KeyType = inst.KeyType.Clone()
		other.ValueType = inst.ValueType.Clone()
	} else {
		other.KeyType = inst.KeyType
		other.ValueType = inst.ValueType
	}
	return other
}

func (inst *ComplexType) stringify(st *SimpleType) string {
	if st == nil {
		return ""
	}
	return st.String()
}

func (inst *ComplexType) String() string {
	if inst.IsArray {
		str := inst.stringify(inst.ValueType)
		return "[]" + str
	} else if inst.IsMap {
		str1 := inst.stringify(inst.KeyType)
		str2 := inst.stringify(inst.ValueType)
		return "map[" + str1 + "]" + str2
	}
	return inst.stringify(inst.ValueType)
}

////////////////////////////////////////////////////////////////////////////////

package model

// StructField 表示结构的一个字段
type StructField struct {
	Name     string
	Tag      string
	TypeExpr string

	Type     ComplexType
	TagItems map[string]string
}

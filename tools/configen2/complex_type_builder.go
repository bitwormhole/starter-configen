package configen2

import "strings"

type ComplexTypeBuilder struct {
	parts []string
	size  int
}

func (inst *ComplexTypeBuilder) Init(typeStr string) {
	builder := &strings.Builder{}
	list := []string{}
	size := len(typeStr)
	for i := 0; i < size; i++ {
		b := typeStr[i]
		if b == '[' || b == ']' || b == '*' || b == '.' {
			list = inst.flush(builder, list)
			builder.WriteByte(b)
			list = inst.flush(builder, list)
		} else {
			builder.WriteByte(b)
		}
	}
	list = inst.flush(builder, list)
	inst.size = len(list)
	inst.parts = list
}

func (inst *ComplexTypeBuilder) flush(buffer *strings.Builder, list []string) []string {
	str := buffer.String()
	str = strings.TrimSpace(str)
	buffer.Reset()
	if str == "" {
		return list
	}
	list = append(list, str)
	return list
}

func (inst *ComplexTypeBuilder) FindPackageAliases() []int {
	list1 := inst.parts
	list2 := []int{}
	for index := range list1 {
		item := list1[index]
		if item == "." {
			if index > 0 {
				list2 = append(list2, index-1)
			}
		}
	}
	return list2
}

func (inst *ComplexTypeBuilder) GetPart(index int) string {
	if 0 <= index && index < inst.size {
		return inst.parts[index]
	}
	return ""
}

func (inst *ComplexTypeBuilder) SetPart(index int, value string) {
	if 0 <= index && index < inst.size {
		inst.parts[index] = value
	}
}

func (inst *ComplexTypeBuilder) String() string {
	builder := &strings.Builder{}
	list := inst.parts
	for index := range list {
		p := strings.TrimSpace(list[index])
		if p == "" {
			continue
		}
		builder.WriteString(p)
	}
	return builder.String()
}

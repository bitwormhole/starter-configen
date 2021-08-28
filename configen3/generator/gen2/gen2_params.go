package gen2

import (
	"errors"
	"strings"

	"github.com/bitwormhole/starter-configen/configen3/data/model"
)

type ComInfo struct {
	Struct        *model.Struct
	InstanceField *model.StructField

	ProxyStructName    string
	InstanceStructName string
	FactoryStructName  string // 'comFactoryForXXX'

	ComID            string
	ComClass         string
	ComAliases       string
	ComScope         string
	ComInitMethod    string
	ComDestroyMethod string
}

func (inst *ComInfo) init(item *model.Struct) error {

	tags := inst.mixTags(item)
	instStructName, err := inst.getInstanceStructName(item)

	if err != nil {
		return err
	}

	inst.ComID = tags["id"]
	inst.ComClass = tags["class"]
	inst.ComScope = tags["scope"]
	inst.ComAliases = tags["aliases"]
	inst.ComInitMethod = tags["initMethod"]
	inst.ComDestroyMethod = tags["destroyMethod"]

	inst.Struct = item
	inst.ProxyStructName = item.Name
	inst.InstanceStructName = instStructName
	inst.FactoryStructName = "comFactory4" + item.Name

	return nil
}

func (inst *ComInfo) getInstanceStructName(item *model.Struct) (string, error) {
	f := inst.InstanceField
	if f == nil {
		return "", errors.New("no field:'instance' in struct:" + item.Name)
	}
	vt := f.Type.ValueType
	name := vt.String()
	index := strings.Index(name, "*")
	if index >= 0 {
		name = name[index+1:]
	}
	return name, nil
}

func (inst *ComInfo) mixTags(item *model.Struct) map[string]string {
	dst := make(map[string]string)
	src := item.Fields
	for _, f := range src {
		kvs := f.TagItems
		for k, v := range kvs {
			dst[k] = v
		}
		if f.Name == "instance" {
			inst.InstanceField = f
		}
	}
	return dst
}

////////////////////////////////////////////////////////////////////////////////

type FieldInfo struct {
	Field        *model.StructField
	Owner        *ComInfo
	SelectorExpr string
	SelectorName string
	SetterName   string
	GetterName   string
}

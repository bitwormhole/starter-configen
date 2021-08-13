package configen1

import (
	"errors"
	"strings"
)

// ComConfigInfo (Component Config Info)
type ComConfigInfo struct {
	Enable string

	ID                    string
	Class                 string
	Aliases               string
	InjectionFuncName     string
	TargetTypeSimpleName  string
	TargetTypePackagePath string
	Scope                 string
	InitMethod            string
	DestroyMethod         string
}

func (inst *ComConfigInfo) init(dom *DomInjection) error {

	// target-type
	targetTypeStr := dom.TargetType
	index := strings.IndexRune(targetTypeStr, '#')
	if index < 1 {
		return errors.New("bad target type: " + targetTypeStr)
	}
	inst.TargetTypePackagePath = targetTypeStr[0:index]
	inst.TargetTypeSimpleName = targetTypeStr[index+1:]

	// injection
	inst.InjectionFuncName = dom.Name

	// propperties
	props := dom.Properties
	if props != nil {
		inst.ID = props["component.id"]
		inst.Class = props["component.class"]
		inst.Scope = props["component.scope"]
		inst.Enable = props["component.enable"]
		inst.Aliases = props["component.aliases"]
		inst.InitMethod = props["component.initmethod"]
		inst.DestroyMethod = props["component.destroymethod"]
	}

	if inst.ID == "" {
		inst.ID = inst.InjectionFuncName
	}

	return nil
}

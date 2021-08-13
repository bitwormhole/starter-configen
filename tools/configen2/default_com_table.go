package configen2

type defaultComInfoTable struct {
	Context *Context
}

func (inst *defaultComInfoTable) init(Context *Context) ComponentInfoTable {
	return inst
}

func (inst *defaultComInfoTable) Add(info *ComponentInfo) {

}

func (inst *defaultComInfoTable) All() []*ComponentInfo {
	return nil
}

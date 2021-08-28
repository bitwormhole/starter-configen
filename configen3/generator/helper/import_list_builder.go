package helper

import (
	"crypto/sha1"
	"sort"
	"strings"

	"github.com/bitwormhole/starter-configen/configen3/data/model"
	"github.com/bitwormhole/starter/util"
)

type importListItem *model.GoImport

// ImportListBuilder 是一个 go-import 列表构建器
type ImportListBuilder struct {
	items map[string]importListItem
}

func (inst *ImportListBuilder) getTable() map[string]importListItem {
	table := inst.items
	if table == nil {
		table = make(map[string]importListItem)
		inst.items = table
	}
	return table
}

func (inst *ImportListBuilder) clone(item importListItem) importListItem {
	if item == nil {
		return nil
	}
	o2 := &model.GoImport{}
	*o2 = *item
	return o2
}

// makeAlias 计算 model.GoImport 的别名
func (inst *ImportListBuilder) makeAliasWithHash(i *model.GoImport) string {
	path := i.Path
	sum := sha1.Sum([]byte(path))
	alias := inst.makeAliasWithoutHash(i)
	return alias + "0x" + util.StringifyBytes(sum[0:3])
}

// makeAliasDefault 计算 model.GoImport 的默认别名
func (inst *ImportListBuilder) makeAliasWithoutHash(i *model.GoImport) string {
	path := i.Path
	index := strings.LastIndex(path, "/")
	suffix := path
	builder := strings.Builder{}
	if index > 0 {
		suffix = path[index+1:]
	}
	array := []rune(suffix)
	for _, ch := range array {
		if ch == '-' {
			ch = '_'
		}
		builder.WriteRune(ch)
	}
	return builder.String()
}

// Create 创建go-import 列表
func (inst *ImportListBuilder) Create() []*model.GoImport {

	table := inst.getTable()
	keys := make([]string, 0)

	for key := range table {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	list := make([]*model.GoImport, 0)

	for _, key := range keys {
		item := table[key]
		list = append(list, item)
	}

	return list
}

// AddComplexType 添加一条go-import 记录
func (inst *ImportListBuilder) AddComplexType(t *model.ComplexType, disableHash bool) {
	if t == nil {
		return
	}
	inst.AddSimpleType(t.KeyType, disableHash)
	inst.AddSimpleType(t.ValueType, disableHash)
}

// AddSimpleType 添加一条go-import 记录
func (inst *ImportListBuilder) AddSimpleType(t *model.SimpleType, disableHash bool) {
	if t == nil {
		return
	}
	res := inst.AddPackagePath(t.PackageName, disableHash)
	if res == nil {
		return
	}
	t.PackageAlias = res.Alias
}

// AddPackagePath 添加一条go-import 记录
func (inst *ImportListBuilder) AddPackagePath(pkgPath string, disableHash bool) *model.GoImport {
	adding := &model.GoImport{}
	adding.Path = pkgPath
	adding.NoHash = disableHash
	return inst.Add(adding)
}

// Add 添加一条go-import 记录
func (inst *ImportListBuilder) Add(adding *model.GoImport) *model.GoImport {
	if adding == nil {
		return nil
	}
	adding = inst.clone(adding)
	table := inst.getTable()
	path := adding.Path
	older := table[path]
	if older == nil {
		aliasWithHash := inst.makeAliasWithHash(adding)
		alias := inst.makeAliasWithoutHash(adding)
		if adding.NoHash {
			adding.Alias = alias
			adding.DefaultAlias = alias
		} else {
			adding.Alias = aliasWithHash
			adding.DefaultAlias = alias
		}
		older = adding
		table[path] = adding
	}
	return inst.clone(older)
}

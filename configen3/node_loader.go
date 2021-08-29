package configen3

import (
	"errors"
	"strings"

	"github.com/bitwormhole/starter-configen/configen3/data/model"
	"github.com/bitwormhole/starter/collection"
	"github.com/bitwormhole/starter/io/fs"
	"github.com/bitwormhole/starter/vlog"
)

type nodeLoader struct {
}

func (inst *nodeLoader) loadProperties(file fs.Path) (collection.Properties, error) {
	text, err := file.GetIO().ReadText(nil)
	if err != nil {
		return nil, err
	}
	return collection.ParseProperties(text, nil)
}

func (inst *nodeLoader) Load(file fs.Path) (*model.Node, error) {

	vlog.Debug("load configen.properties, path=" + file.Path())

	props, err := inst.loadProperties(file)
	if err != nil {
		return nil, err
	}

	node := &model.Node{}
	node.Path = file
	node.Properties = props

	err = inst.parseBase(props, node)
	if err != nil {
		return nil, err
	}

	err = inst.parseChildren(props, node)
	if err != nil {
		return nil, err
	}

	err = inst.parseScans(props, node)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (inst *nodeLoader) parseBase(props collection.Properties, dst *model.Node) error {

	const wantVersion = "3"

	var err error = nil
	dst.Version, err = props.GetPropertyRequired("configen.version")
	if err != nil {
		return err
	}

	if dst.Version != wantVersion {
		return errors.New("unsupported configen.version:" + dst.Version)
	}

	return nil
}

func (inst *nodeLoader) findNamesByKeyPrefixAndSuffix(props collection.Properties, prefix string, suffix string) []string {
	list := make([]string, 0)
	table := props.Export(nil)
	for key := range table {
		if strings.HasPrefix(key, prefix) && strings.HasSuffix(key, suffix) {
			name := key[len(prefix) : len(key)-len(suffix)]
			list = append(list, name)
		}
	}
	return list
}

func (inst *nodeLoader) parseChildren(props collection.Properties, dst *model.Node) error {

	const prefix = "child."
	const suffix = ".path"
	names := inst.findNamesByKeyPrefixAndSuffix(props, prefix, suffix)

	for _, name := range names {
		kp := prefix + name + "." // the key prefix
		child := &model.Child{}
		child.Name = name
		child.Href = props.GetProperty(kp+"path", "(undefine)")
		dst.Children = append(dst.Children, child)
	}

	return nil
}

func (inst *nodeLoader) parseScans(props collection.Properties, dst *model.Node) error {

	const prefix = "scan."
	const suffix = ".path"
	names := inst.findNamesByKeyPrefixAndSuffix(props, prefix, suffix)
	dir := dst.Path

	if dir.IsFile() {
		dir = dir.Parent()
	}

	for _, name := range names {
		kp := prefix + name + "." // the key prefix
		scan := &model.Scan{}
		scan.Name = name
		scan.Href = props.GetProperty(kp+"path", "(undefine)")
		scan.WriteToChild = props.GetProperty(kp+"write-to-child", "(undefine)")
		scan.WriteToFile = props.GetProperty(kp+"write-to-file", "(undefine)")
		scan.R = props.GetProperty(kp+"r", "false") == "true"
		scan.Path = dir.GetChild(scan.Href)
		write2, err := inst.getFileWriteTo(scan, dst, dir)
		if err != nil {
			return err
		}
		scan.WriteTo = write2
		dst.ScanTargets = append(dst.ScanTargets, scan)
	}

	return nil
}

func (inst *nodeLoader) getFileWriteTo(scan *model.Scan, node *model.Node, dir fs.Path) (fs.Path, error) {

	var child *model.Child = nil
	childName := scan.WriteToChild
	children := node.Children

	for _, ch := range children {
		if ch.Name == childName {
			child = ch
			break
		}
	}

	if child == nil {
		return nil, errors.New("no child node named:" + childName)
	}

	path := child.Href
	file := "auto_generated_node-" + childName

	return dir.GetChild(path).GetChild(file + ".go"), nil
}

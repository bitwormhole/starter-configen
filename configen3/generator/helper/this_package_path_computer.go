package helper

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bitwormhole/starter-configen/configen3/configen"
	"github.com/bitwormhole/starter/io/fs"
	"github.com/bitwormhole/starter/util"
)

// ThisPackagePathComputer 工具：根据路径取完整的包名
type ThisPackagePathComputer struct{}

// Compute 计算包名
func (inst *ThisPackagePathComputer) Compute(ctx *configen.Context, path fs.Path) (string, error) {

	file := path
	if file.IsDir() {
		file = file.GetChild("index.tmp")
	} else if file.IsFile() {
		// NOP
	} else {
		// NOP
	}

	root := ctx.Store.RootNode.Path
	base := root.Parent().Path()
	full := file.Parent().Path()
	if !strings.HasPrefix(full, base) {
		msg := fmt.Sprint("the path[", full, "] is not start with prefix [", base, "]")
		return "", errors.New(msg)
	}

	part1, err := inst.loadGoModuleQName(root)
	if err != nil {
		return "", err
	}
	part2 := full[len(base):]

	pb := &util.PathBuilder{}
	pb.AppendPath(part1 + "/" + part2)
	return pb.String(), nil
}

func (inst *ThisPackagePathComputer) loadGoModuleQName(root fs.Path) (string, error) {

	file := root.Parent().GetChild("go.mod")
	text, err := file.GetIO().ReadText(nil)
	if err != nil {
		return "", err
	}
	text = strings.ReplaceAll(text, "\r", "\n")
	array := strings.Split(text, "\n")

	const space = " "
	const tab = "\t"

	for _, line := range array {
		line = strings.TrimSpace(line)
		line = strings.ReplaceAll(line, tab, space)
		a2 := strings.SplitN(line, space, 2)
		if len(a2) == 2 {
			str1 := strings.TrimSpace(a2[0])
			str2 := strings.TrimSpace(a2[1])
			if str1 == "module" {
				return str2, nil
			}
		}
	}

	return "", errors.New("no 'module' define in file: " + file.Path())
}

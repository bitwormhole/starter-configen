package help

import "fmt"

func RunHelpDetail(args []string) error {

	const tab = "\t"
	//	RunVersion(args)
	fmt.Println("usage: starter {command} {more...args}")

	fmt.Println("commands:")
	fmt.Println(tab, "about:    显示关于信息")
	fmt.Println(tab, "configen: 运行配置代码生成工具")
	fmt.Println(tab, "help:     显示帮助信息")
	fmt.Println(tab, "version:  显示版本信息")

	return nil
}

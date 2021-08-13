package help

import "fmt"

func RunVersion(args []string) error {

	about := &AboutInfo{}
	about.init()

	name := about.ApplicationName
	version := about.Version

	fmt.Println(name, version)
	return nil
}

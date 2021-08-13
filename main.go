package main

import (
	"github.com/bitwormhole/starter"
	"github.com/bitwormhole/starter-configen/etc/configen"
)

func main() {
	// fmt.Println("hello,starter-configen")
	i := starter.InitApp()
	i.Use(configen.Module())
	i.Run()
}

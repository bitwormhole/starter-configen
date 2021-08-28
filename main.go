package main

import (
	"embed"

	"github.com/bitwormhole/starter-configen/configen3"
	"github.com/bitwormhole/starter-configen/configen3/configen"
	"github.com/bitwormhole/starter/collection"
	"github.com/bitwormhole/starter/vlog"
)

//go:embed src/main/resources
var theRES embed.FS

func main() {

	vlog.UseSimpleLogger(vlog.DEBUG)

	res := collection.LoadEmbedResources(&theRES, "src/main/resources")
	ctx := &configen.Context{}
	mp := &configen3.MainProcess{}

	ctx.Resources = res
	err := mp.Run(ctx)
	if err != nil {
		panic(err)
	}
}

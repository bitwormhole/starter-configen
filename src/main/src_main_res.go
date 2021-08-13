package srcmain

import (
	"embed"

	"github.com/bitwormhole/starter/application/config"
	"github.com/bitwormhole/starter/collection"
)

//go:embed resources
var resources embed.FS

func ExportResources() collection.Resources {
	return config.LoadResourcesFromEmbedFS(&resources, "resources")
}

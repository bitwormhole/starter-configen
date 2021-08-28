package srcmain

import (
	"embed"

	"github.com/bitwormhole/starter/collection"
)

//go:embed resources
var resources embed.FS

// ExportResources 导出资源
func ExportResources() collection.Resources {
	return collection.LoadEmbedResources(&resources, "resources")
}

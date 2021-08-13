package conf

import (
	"github.com/bitwormhole/starter-configen/tools/configen2/templates"
	"github.com/bitwormhole/starter/application"
	"github.com/bitwormhole/starter/markup"
)

func ExportConfig(cb application.ConfigBuilder, mod application.Module) error {

	return autoGenConfig(cb)
}

////////////////////////////////////////////////////////////////////////////////

type theMainTemplateFactory struct {
	markup.Component
	instance *templates.MainTemplateFactory `id:"configen2-main-template-factory"`
}

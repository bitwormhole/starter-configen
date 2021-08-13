package demo

import (
	"testing"

	"github.com/bitwormhole/starter-configen/tools/configen1"
	"github.com/bitwormhole/starter/application"
)

func TestConfigen(t *testing.T) {

	args := []string{}
	context, _ := application.Run(nil, args)

	err := configen1.Run(context, args)
	if err != nil {
		t.Fatal(err)
	}
}

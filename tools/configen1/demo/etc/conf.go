package etc

import (
	xstr "strings"

	"github.com/bitwormhole/starter-configen/tools/configen1/demo/car"
	"github.com/bitwormhole/starter/application"
)

func car1(t *car.Body, context application.Context) error {

	// [component]
	// id=abc

	return nil
}

func car2(t *car.Body, context application.Context) error {

	// [component]
	// id=body2
	// class=body

	return nil
}

func door(t *car.Door, context application.Context) error {

	// [component]
	// id=door1
	// class=door

	return nil
}

func builder(t *xstr.Builder, context application.Context) error {

	// [component]
	// id=strbuilder1
	// class=strBuilder

	return nil
}

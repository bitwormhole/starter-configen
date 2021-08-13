package etc

import (
	xstr "strings"

	"github.com/bitwormhole/starter-configen/tools/configen1/demo/car"
	"github.com/bitwormhole/starter/application"
)

func car11(t *car.Body, context application.Context) error {

	// [component]
	// class=abc

	return nil
}

func car22(t *car.Body, context application.Context) error {

	// [component]
	// id=body2
	// class=body
	// aliases=			c22 c33 c44 		c50			c666
	// scope=singleton
	// initMethod=Start
	// destroyMethod=Stop

	return nil
}

func builder3(t *xstr.Builder, context application.Context) error {

	// [component]
	//    aliases= a1 a2 a3

	return nil
}

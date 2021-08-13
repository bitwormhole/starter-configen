package configen1

import (
	"errors"
	"os"

	"github.com/bitwormhole/starter/application"
	"github.com/bitwormhole/starter/io/fs"
)

func Run(ctx application.Context, args []string) error {

	context := &configenContext{}

	// PWD
	const ENV_PWD = "PWD"
	pwd := os.Getenv(ENV_PWD)
	if pwd == "" {
		return errors.New("no ENV: " + ENV_PWD)
	}
	fsys := fs.Default()
	context.PWD = fsys.GetPath(pwd)

	// runner
	runner := &configenRunner{context: context}
	err := runner.run()
	if err != nil {
		return err
	}

	return nil
}

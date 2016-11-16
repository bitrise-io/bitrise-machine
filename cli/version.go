package cli

import (
	"fmt"
	"runtime"

	"github.com/urfave/cli"
)

func version(c *cli.Context) {
	fmt.Println(c.App.Version)

	if c.Bool(FullFlagKey) {
		fmt.Println()
		fmt.Println("go: " + runtime.Version())
		fmt.Println("arch: " + runtime.GOARCH)
		fmt.Println("os: " + runtime.GOOS)
	}
}

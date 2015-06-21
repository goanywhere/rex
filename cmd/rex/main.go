package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/goanywhere/crypto"
)

var (
	cwd string
)

var commands = []cli.Command{
	// rex project template supports
	/*
	 *{
	 *    Name:   "new",
	 *    Usage:  "create a skeleton web application project",
	 *    Action: New,
	 *},
	 */
	// rex server (with livereload supports)
	{
		Name:   "run",
		Usage:  "start application server with livereload supports",
		Action: Run,
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "port",
				Value: 5000,
				Usage: "port to run the application server",
			},
		},
	},
	// helper to generate a secret key.
	{
		Name:  "secret",
		Usage: "generate a new application secret key",
		Action: func(ctx *cli.Context) {
			fmt.Println(crypto.Random(ctx.Int("length")))
		},
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "length",
				Value: 64,
				Usage: "length of the secret key",
			},
		},
	},
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd := cli.NewApp()
	cmd.Name = "rex"
	cmd.Usage = "manage rex application"
	cmd.Version = "0.9.0"
	cmd.Author = "GoAnywhere"
	cmd.Email = "code@goanywhere.io"
	cmd.Commands = commands
	cmd.Run(os.Args)
}

func init() {
	cwd, _ = os.Getwd()
}

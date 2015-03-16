/* ----------------------------------------------------------------------
 *       ______      ___                         __
 *      / ____/___  /   |  ____  __  ___      __/ /_  ___  ________
 *     / / __/ __ \/ /| | / __ \/ / / / | /| / / __ \/ _ \/ ___/ _ \
 *    / /_/ / /_/ / ___ |/ / / / /_/ /| |/ |/ / / / /  __/ /  /  __/
 *    \____/\____/_/  |_/_/ /_/\__. / |__/|__/_/ /_/\___/_/   \___/
 *                            /____/
 *
 * (C) Copyright 2015 GoAnywhere (http://goanywhere.io).
 * ----------------------------------------------------------------------
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 * ----------------------------------------------------------------------*/

package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/goanywhere/x/crypto"
)

var (
	cwd string
)

var commands = []cli.Command{
	// rex project template supports
	{
		Name:   "new",
		Usage:  "create a skeleton web application project",
		Action: New,
	},
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
			cli.StringFlag{
				Name:  "path",
				Value: ".",
				Usage: "path to run the application server",
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

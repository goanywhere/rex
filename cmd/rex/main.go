/* ----------------------------------------------------------------------
 *       ______      ___                         __
 *      / ____/___  /   |  ____  __  ___      __/ /_  ___  ________
 *     / / __/ __ \/ /| | / __ \/ / / / | /| / / __ \/ _ \/ ___/ _ \
 *    / /_/ / /_/ / ___ |/ / / / /_/ /| |/ |/ / / / /  __/ /  /  __/
 *    \____/\____/_/  |_/_/ /_/\__. / |__/|__/_/ /_/\___/_/   \___/
 *                            /____/
 *
 * (C) Copyright 2014 GoAnywhere (http://goanywhere.io).
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
	"os"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/goanywhere/env"
)

var commands = []cli.Command{
	{
		Name:   "new",
		Usage:  "create a skeleton web application project",
		Action: New,
	},
}

var flags = []cli.Flag{
	cli.IntFlag{
		Name:  "port",
		Value: 5000,
		Usage: "port to run application server",
	},
	cli.StringFlag{
		Name:  "npm",
		Value: "build",
		Usage: "script for npm (http://npmjs.com/)",
	},
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd := cli.NewApp()
	cmd.Name = "rex"
	cmd.Usage = "manage Rex application project"
	cmd.Version = "0.1.1"
	cmd.Author = "GoAnywhere"
	cmd.Email = "opensource@goanywhere.io"
	cmd.Commands = commands
	cmd.Flags = flags
	cmd.Action = func(ctx *cli.Context) {
		env.Set("Port", ctx.String("port"))
		app := NewApp()
		app.Script = ctx.String("npm")
		app.Start()
	}
	cmd.Run(os.Args)
}

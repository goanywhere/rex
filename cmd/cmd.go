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

package cmd

import (
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/codegangsta/cli"
)

var here string

var commands = []cli.Command{
	{
		Name:   "new",
		Usage:  "create a skeleton web application project",
		Action: Create,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "lang, l",
				Value: "english",
				Usage: "language for the greeting",
			},
		},
	},
	{
		Name:   "serve",
		Usage:  "start serving HTTP request with live reload supports",
		Action: Serve,
	},
}

func Execute() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	app := cli.NewApp()
	app.Name = "webapp"
	app.Usage = "manage web application project"
	app.Version = "0.1.1"
	app.Commands = commands
	app.Run(os.Args)
}

func init() {
	_, filename, _, _ := runtime.Caller(1)
	here, _ = filepath.Abs(path.Dir(filename))
}

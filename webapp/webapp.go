/**
 *  ------------------------------------------------------------
 *  @project	web.go
 *  @file       web.go
 *  @date       2014-10-10
 *  @author     Jim Zhan <jim.zhan@me.com>
 *
 *  Copyright © 2014 Jim Zhan.
 *  ------------------------------------------------------------
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
 *  ------------------------------------------------------------
 */
package main

import (
	"os"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/goanywhere/web/cmd"
)

var commands = []cli.Command{
	{
		Name:   "new",
		Usage:  "create a skeleton web application project",
		Action: cmd.Create,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "lang, l",
				Value: "english",
				Usage: "language for the greeting",
			},
		},
	},
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	app := cli.NewApp()
	app.Name = "webapp"
	app.Usage = "manage web application project"
	app.Version = "0.0.1"
	app.Commands = commands
	app.Run(os.Args)
}

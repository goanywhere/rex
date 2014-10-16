/**
 *  ------------------------------------------------------------
 *  @project
 *  @file       webapp.go
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
package webapp

import (
	"os"
)
import "github.com/codegangsta/cli"

func create(context *cli.Context) {
	args := context.Args()
	if len(args) != 1 {
		Error("Valid Project Name Missing")
	} else {
		// create skeleton here
	}
}

func Execute() {
	app := cli.NewApp()
	app.Name = "webapp"
	app.Usage = "manage web application project"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:   "create",
			Usage:  "create a skeleton web application project",
			Action: create,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "lang, l",
					Value: "english",
					Usage: "language for the greeting",
				},
			},
		},
	}
	app.Run(os.Args)
}

/**
 *  ------------------------------------------------------------
 *  @project
 *  @file       serve.go
 *  @date       2014-11-13
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
package cmd

import (
	"os"
	"path"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/pilu/fresh/runner"
)

func Serve(context *cli.Context) {
	_, filename, _, _ := runtime.Caller(1)
	config := path.Join(path.Dir(filename), "webapp.conf")
	os.Setenv("RUNNER_CONFIG_PATH", config)
	runner.Start()
}

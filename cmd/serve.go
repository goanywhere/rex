/**
 * ----------------------------------------------------------------------
 *       ______      ___                         __
 *      / ____/___  /   |  ____  __  ___      __/ /_  ___  ________
 *     / / __/ __ \/ /| | / __ \/ / / / | /| / / __ \/ _ \/ ___/ _ \
 *    / /_/ / /_/ / ___ |/ / / / /_/ /| |/ |/ / / / /  __/ /  /  __/
 *    \____/\____/_/  |_/_/ /_/\__. / |__/|__/_/ /_/\___/_/   \___/
 *                            /____/
 * ----------------------------------------------------------------------
 *  Copyright © 2014 GoAnywhere Ltd. All Rights Reserved.
 * ----------------------------------------------------------------------*/

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

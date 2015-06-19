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
package rex

import (
	"flag"
	"path/filepath"
	"runtime"

	"github.com/goanywhere/rex/internal"
	"github.com/goanywhere/x/env"
	"github.com/goanywhere/x/fs"
)

var (
	options = struct {
		debug    bool
		port     int
		maxprocs int
	}{
		debug:    true,
		port:     5000,
		maxprocs: runtime.NumCPU(),
	}
)

func init() {
	// ----------------------------------------
	// Project Root
	// ----------------------------------------
	var root = fs.Getcd(2)
	env.Set(internal.Root, root)
	env.Load(filepath.Join(root, ".env"))

	// cmd arguments
	flag.BoolVar(&options.debug, "debug", options.debug, "flag to toggle debug mode")
	flag.IntVar(&options.port, "port", options.port, "port to run the application server")
	flag.IntVar(&options.maxprocs, "maxprocs", options.maxprocs, "maximum cpu processes to run the server")

	flag.Parse()
}

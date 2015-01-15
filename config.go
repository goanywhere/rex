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
	"log"
	"os"
	"path/filepath"

	"github.com/goanywhere/x/env"
)

var (
	Root   string
	Secret string

	Port = 5000
	Mode = "debug"

	Dir = struct {
		Static    string
		Templates string
	}{
		Static:    "build",
		Templates: "templates",
	}

	URL = struct {
		Static string
	}{
		Static: "/static/",
	}
)

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to retrieve project root: %v", err)
	}
	Root, _ = filepath.Abs(cwd)

	env.Prefix = "rex"
	env.Set("root", Root)
	env.Set("port", "5000")
	env.Set("mode", "debug")
	env.Set("dir.static", "build")
	env.Set("dir.templates", "templates")
	env.Set("url.static", "/static/")

	env.Load(filepath.Join(cwd, ".env"))
	Secret = env.String("secret")
}

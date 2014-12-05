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

package web

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/goanywhere/env"
	"github.com/goanywhere/web/template"
)

var loader *template.Loader

// Shortcut to create map.
type H map[string]interface{}

// FIXME .env not loaded.
func New() *Mux {
	log.Printf("Application initializing...")
	loader = template.NewLoader(Settings.Templates)
	pages := loader.Load()
	log.Printf("Application loaded (%d templates)", pages)
	// Load custom settings.
	env.Load(Settings)
	mux := newMux()
	return mux
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if cwd, err := os.Getwd(); err == nil {
		if root, err := filepath.Abs(cwd); err == nil {
			Settings.Root = root
		} else {
			Panic("[web.go] could not initialize project root: %v", err)
		}
	} else {
		Panic("[web.go] could not retrieve current working directory: %v", err)
	}
}

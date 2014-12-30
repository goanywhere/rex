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

package rex

import (
	"log"
	"os"
	"path/filepath"

	"github.com/goanywhere/env"
	"github.com/goanywhere/rex/config"
	"github.com/goanywhere/rex/http"
	"github.com/goanywhere/rex/middleware"
)

var Settings = config.Settings()

// Shortcut to create map.
type H map[string]interface{}

// Defaults creates a web.Server with common middlewares enabled.
func Defaults() *http.Server {
	server := New()
	server.Use(middleware.LiveReload)
	return server
}

// New creates a plain web.Server.
func New() *http.Server {
	env.Load(filepath.Join(Settings.Root, ".env"))
	env.Dump(Settings)
	return http.NewServer()
}

func init() {
	if cwd, err := os.Getwd(); err == nil {
		Settings.Root, _ = filepath.Abs(cwd)
	} else {
		log.Fatalf("Failed to retrieve project root: %v", err)
	}
}

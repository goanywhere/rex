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

	"github.com/goanywhere/rex/core"
	"github.com/goanywhere/rex/web"
)

type H map[string]interface{}

var (
	Filters  = []web.Filter{}
	Settings = core.Settings()
)

// New creates a plain web.Server.
func New() *web.Server {
	server := web.NewServer()
	for _, filter := range Filters {
		server.Use(filter)
	}
	return server
}

func init() {
	if cwd, err := os.Getwd(); err == nil {
		web.Configure(cwd)
	} else {
		log.Fatalf("Failed to retrieve project root: %v", err)
	}
}

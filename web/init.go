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
package web

import (
	"log"
	"path/filepath"
	"sync"

	"github.com/goanywhere/rex/config"
	"github.com/goanywhere/rex/crypto"
	"github.com/goanywhere/rex/template"
	"github.com/goanywhere/x/env"
)

var (
	once      sync.Once
	signature *crypto.Signature
	settings  = config.Settings()

	loader = template.NewLoader(settings.Templates)
)

func Configure(cwd string) {
	once.Do(func() {
		settings.Root, _ = filepath.Abs(cwd)
		env.Load(filepath.Join(settings.Root, ".env"))
		env.Dump(settings)

		if settings.Secret == "" {
			log.Fatal("Secret key missing")
		}
		// creates a signature for accessing securecookie.
		signature = crypto.NewSignature(settings.Secret)
	})
}
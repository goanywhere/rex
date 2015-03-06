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
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/securecookie"

	"github.com/goanywhere/rex/internal"
	"github.com/goanywhere/rex/template"

	"github.com/goanywhere/x/fs"
)

var (
	settings = internal.Settings()

	// application secret keys
	secrets []securecookie.Codec

	// application page templates (Settings: dir.templates)
	templates *template.Loader
)

func createSecrets(keys ...string) {
	if len(keys) > 0 {
		var bytes [][]byte
		for _, key := range keys {
			bytes = append(bytes, []byte(key))
		}
		secrets = securecookie.CodecsFromPairs(bytes...)
	} else {
		log.Fatalf("Failed to setup application: secret key(s) missing")
	}
}

// configure initialize all application related settings before running.
// Server Secret Keys
func init() {
	settings.Load(".env")
	// ------------------------------------------------
	// if secret keys exists, create codecs.
	// ------------------------------------------------
	createSecrets(settings.Strings("SECRET_KEYS")...)

	// ------------------------------------------------
	// templates folder exists => load HTML templates.
	// ------------------------------------------------
	if dir := settings.String("DIR_TEMPLATES", "templates"); fs.Exists(dir) {
		templates = template.NewLoader(dir)
		templates.Load()
	}
}

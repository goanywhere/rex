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
	"html/template"
	"path"
	"path/filepath"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/securecookie"

	"github.com/goanywhere/rex/internal"

	"github.com/goanywhere/x/env"
	"github.com/goanywhere/x/fs"
)

var (
	FuncMap  = make(template.FuncMap)
	settings = internal.Settings()

	// application secret keys
	secrets []securecookie.Codec

	// application page templates (Settings: dir.templates)
	templates *TemplateLoader
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
	// ------------------------------------------------
	// setup fundamental project root.
	// ------------------------------------------------
	_, filename, _, _ := runtime.Caller(2)
	if root, err := filepath.Abs(path.Dir(filename)); err == nil {
		settings.Root = root
		// custom settings
		env.Load(filepath.Join(root, ".env"))
		env.Map(settings)
		env.Map(settings.Session)
		// ------------------------------------------------
		// templates folder exists => load HTML templates.
		// ------------------------------------------------
		if dir := filepath.Join(root, settings.Views); fs.Exists(dir) {
			templates = NewTemplateLoader(dir)
			templates.Load()
		}
	} else {
		log.Fatalf("Failed to retrieve project root: %v", err)
	}
	// ------------------------------------------------
	// if secret keys exists, create codecs.
	// ------------------------------------------------
	createSecrets(settings.SecretKeys...)
}

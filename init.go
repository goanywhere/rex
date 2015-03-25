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
	"html/template"

	"github.com/gorilla/securecookie"
)

var (
	FuncMap = make(template.FuncMap)

	// application secret keys
	secrets []securecookie.Codec

	// application page templates (Settings.Views)
)

// configure initialize all application related settings before running.
func init() {
	/*
		// ------------------------------------------------
		// templates folder exists => load HTML templates.
		// ------------------------------------------------
		if dir := filepath.Join(root, settings.View); fs.Exists(dir) {
			views = Load(dir)
		}

		// ------------------------------------------------
		// Create application secrets
		// ------------------------------------------------
		if len(settings.SecretKeys) > 0 {
			var bytes [][]byte
			for _, key := range settings.SecretKeys {
				bytes = append(bytes, []byte(key))
			}
			secrets = securecookie.CodecsFromPairs(bytes...)
		} else {
			log.Fatalf("Failed to setup application: secret key(s) missing")
		}
	*/
}

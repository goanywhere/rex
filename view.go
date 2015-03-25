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
	"os"
	"path/filepath"
	"regexp"

	pongo "github.com/flosch/pongo2"
	"github.com/goanywhere/x/fs"
)

var views map[string]*pongo.Template = make(map[string]*pongo.Template)

// loadViews load the html/xml documents from the pre-defined directory,
// rex will ignores directories named "layouts" & "include".
// TODO multiple paths supports.
func loadViews(root string) {
	var (
		files   = regexp.MustCompile(`\.(html|xml)$`)
		ignores = regexp.MustCompile(`(layouts|include|\.(\w+))`)
	)
	if fs.Exists(root) {
		filepath.Walk(root, func(path string, info os.FileInfo, e error) error {

			if info.IsDir() {
				if ignores.MatchString(info.Name()) {
					return filepath.SkipDir
				} else {
					return nil
				}
			}

			if files.MatchString(path) {
				key, _ := filepath.Rel(root, path)
				views[key] = pongo.Must(pongo.FromFile(path))
			}

			return e
		})
	}
}

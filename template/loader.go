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

package template

import (
	"bytes"
	"html/template"
	"path/filepath"
)

var templates *template.Template

func Load(path string) (err error) {
	if templates, err = template.ParseGlob(filepath.Join(path, "*.html")); err == nil {

	}
	return
}

func ExecuteTemplate(values interface{}, layouts ...string) (buffer *bytes.Buffer, err error) {
	for _, layout := range layouts {
		if templates.Lookup(layout) != nil {
			err = templates.ExecuteTemplate(buffer, layout, values)
		}
	}
	return
}

func init() {
	templates = template.New("templates")
}

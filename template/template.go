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
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/goanywhere/regex"
)

// chain finds all ancesters (via "extends") & combines them along
// with the filename iteself into correct order for parsing.
func chain(filename string) []string {
	var paths = []string{filename}
	cwd := filepath.Dir(filename)

	for {
		file, err := os.Open(filename)
		defer file.Close()
		if err != nil {
			break
		}

		buffer := make([]byte, 1024)
		reader := bufio.NewReader(file)

		var path string
		// Check if it contains "extends" tag.
		for {
			if buffer, _, err = reader.ReadLine(); err != nil {
				break
			} else {
				line := string(buffer)
				if match := regex.Tag.Extends.FindStringSubmatch(line); len(match) >= 2 {
					path = filepath.Join(cwd, match[1])
					if path == filename {
						panic(fmt.Errorf("web/template: template cannot extend itself (%s)", filename))
					}
					paths = append([]string{path}, paths...)
					break
				}
			}
		}
		// move to the ancester to check
		filename = path
	}
	return paths
}

// TODO modular filtering.
//func Filter(filename string, pattern *regexp.Regexp, filter func()()) {}

// Parse finds all extends chain & constructs the final page layout.
func Parse(filename string) *template.Template {
	var err error
	var page *template.Template

	filenames := chain(filename)

	for _, item := range filenames {
		if bits, err := ioutil.ReadFile(item); err == nil {
			// remove the custom tag "extends" for standard parsing.
			content := regex.Tag.Extends.ReplaceAllString(string(bits), "")

			var tmpl *template.Template
			if page == nil {
				page = template.New(item)
			}
			if item == page.Name() {
				tmpl = page
			} else {
				tmpl = page.New(item)
			}
			_, err = tmpl.Parse(content)
		}
	}

	return template.Must(page, err)
}

func Load(path string) (err error) {
	return
}

func ExecuteTemplate(values interface{}, layouts ...string) (buffer *bytes.Buffer, err error) {
	return
}

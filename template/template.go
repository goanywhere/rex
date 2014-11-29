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
	"fmt"
	"os"
	"path/filepath"

	"github.com/goanywhere/regex/tags"
)

type Page struct {
	Name   string
	loader *Loader
}

// Ancesters finds all ancesters using jinja's syntax & combines
// them along with the filename iteself into correct order for parsing.
// tag: {% extends "layout/base.html" %}
func (self *Page) Ancestors() (paths []string) {
	paths = append(paths, filepath.Join(self.loader.path, self.Name))

	var filename string = paths[0]

	for {
		file, err := os.Open(filename)
		defer file.Close()
		if err != nil {
			break
		}
		// find the very first "extends" tag.
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			matches := tags.Extends.FindStringSubmatch(scanner.Text())
			if len(matches) == 2 {
				if path := filepath.Join(self.loader.path, matches[1]); path == filename {
					panic(fmt.Errorf("web/template: template cannot extend itself (%s)", filename))
				} else {
					paths = append([]string{path}, paths...)
					filename = path // move to the ancester to check.
					break           // Only the very first one.
				}
			}
		}
	}
	return
}

// Include finds all included external file sources recursively
// & replace all the "include" tags with their actual sources.
func (self *Page) Include() (src string) {
	return
}

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
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/goanywhere/regex/tags"
)

type Page struct {
	Name   string
	loader *Loader
}

func (self *Page) path() string {
	return path.Join(self.loader.root, self.Name)
}

func (self *Page) source() (src string) {
	if bits, err := ioutil.ReadFile(self.path()); err == nil {
		src = string(bits)
	}
	return
}

// Ancesters finds all ancestors absolute path using jinja's syntax
// and combines them along with the page path iteself into correct order for parsing.
// tag: {% extends "layout/base.html" %}
func (self *Page) Ancestors() (paths []string) {
	var name = self.Name
	paths = append(paths, name)

	for {
		var orphan = false
		file, err := os.Open(path.Join(self.loader.root, name))
		defer file.Close()
		if err != nil {
			panic(fmt.Errorf("web/template: %v", err))
		}
		// find the very first "extends" tag.
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			result := tags.Extends.FindStringSubmatch(scanner.Text())
			if len(result) == 2 {
				if name == result[1] {
					panic(fmt.Errorf("web/template: template cannot extend itself (%s)", name))
				} else {
					paths = append([]string{result[1]}, paths...) // insert the ancester into the first place.
					name = result[1]
					break
				}
			} else {
				orphan = true
			}
		}

		if orphan {
			break
		}
	}

	return
}

// Include finds all included external file sources recursively
// & replace all the "include" tags with their actual sources.
// tag: {% include "partials/header.html" %}
func (self *Page) Include() (source string) {
	bits, err := ioutil.ReadFile(self.path())
	if err != nil {
		panic(fmt.Errorf("web/template: template cannot be opened (%s)", self.Name))
	}

	source = string(bits)
	for {
		result := tags.Include.FindAllStringSubmatch(source, -1)
		if result == nil {
			break
		}

		for _, match := range result {
			tag, name := match[0], match[1]
			if name == self.Name {
				panic(fmt.Errorf("web/template: template cannot include itself (%s)", name))
			}
			page := self.loader.Load(name)
			// reconstructs source to recursively find all included sources.
			source = strings.Replace(source, tag, page.source(), -1)
		}
	}
	return
}

func (self *Page) Parse() (output *template.Template) {
	var err error
	paths := self.Ancestors()

	for _, path := range paths {
		page := self.loader.Load(path)
		var tmpl *template.Template

		if output == nil {
			output = template.New(path)
		}
		if path == output.Name() {
			tmpl = output
		} else {
			tmpl = output.New(path)
		}
		_, err = tmpl.Parse(page.Include())
	}
	return template.Must(output, err)
}

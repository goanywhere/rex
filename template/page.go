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
	"html/template"
	"io/ioutil"
	"path"
	"strings"

	"github.com/goanywhere/regex/tags"
	"github.com/goanywhere/web"
)

type Page struct {
	Name   string  // name of the page under laoder's root path.
	loader *Loader // file loader.
}

// Ancesters finds all ancestors absolute path using jinja's syntax
// and combines them along with the page path iteself into correct order for parsing.
// tag: {% extends "layout/base.html" %}
func (self *Page) ancestors() (names []string) {
	var name = self.Name
	names = append(names, name)

	for {
		// find the very first "extends" tag.
		bits, err := ioutil.ReadFile(path.Join(self.loader.root, name))
		if err != nil {
			web.Panic("web/template: %v", err)
		}

		result := tags.Extends.FindSubmatch(bits)
		if result == nil {
			break
		}

		base := string(result[1])
		if base == name {
			web.Panic("web/template: template cannot extend itself (%s)", name)
		}

		names = append([]string{base}, names...) // insert the ancester into the first place.
		name = base
	}

	return
}

// Include finds all included external file sources recursively
// & replace all the "include" tags with their actual sources.
// tag: {% include "partials/header.html" %}
func (self *Page) include() (source string) {
	bits, err := ioutil.ReadFile(self.path())
	if err != nil {
		web.Panic("web/template: template cannot be opened (%s)", self.Name)
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
				web.Panic("web/template: template cannot include itself (%s)", name)
			}
			page := self.loader.Page(name)
			// reconstructs source to recursively find all included sources.
			source = strings.Replace(source, tag, page.source(), -1)
		}
	}
	return
}

// Path returns the abolute path of the page.
func (self *Page) path() string {
	return path.Join(self.loader.root, self.Name)
}

// Parse constructs `template.Template` object with additional
// "extends" & "include" like Jinja.
func (self *Page) parse() (output *template.Template) {
	var err error
	names := self.ancestors()

	for _, name := range names {
		page := self.loader.Page(name)
		var tmpl *template.Template

		if output == nil {
			output = template.New(name)
		}
		if name == output.Name() {
			tmpl = output
		} else {
			tmpl = output.New(name)
		}
		_, err = tmpl.Parse(page.include())
	}
	return template.Must(output, err)
}

// Source returns the plain raw source of the page.
func (self *Page) source() (src string) {
	if bits, err := ioutil.ReadFile(self.path()); err == nil {
		src = string(bits)
	} else {
		web.Panic("web/template: template cannot be opened (%s)", self.Name)
	}
	return src
}

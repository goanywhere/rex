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

package web

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/goanywhere/regex/tags"
)

type (
	Loader struct {
		root      string
		loaded    bool
		mutex     sync.RWMutex
		templates map[string]*template.Template
	}

	page struct {
		Name   string  // name of the page under laoder's root path.
		loader *Loader // file loader.
	}
)

func NewLoader(path string) *Loader {
	abspath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Failed to initialize templates path: %v", err)
	}
	loader := new(Loader)
	loader.root = abspath
	loader.templates = make(map[string]*template.Template)
	return loader
}

// Exists checks if the given filename exists under the root.
func (self *Loader) Exists(name string) bool {
	abspath := filepath.Join(self.root, name)
	if _, err := os.Stat(abspath); os.IsNotExist(err) {
		return false
	}
	return true
}

// Files lists all HTML files under the root.
func (self *Loader) Files() (names []string) {
	err := filepath.Walk(self.root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".html") {
			if name, e := filepath.Rel(self.root, path); e == nil {
				names = append(names, name)
			} else {
				err = e
			}
		}
		return err
	})

	if err != nil {
		log.Fatalf("rex/template: files list cannot be listed: %v", err)
	}

	return
}

// Get retrieves the parsed template from preloaded pool.
func (self *Loader) Get(name string) *template.Template {
	self.Load()
	return self.templates[name]
}

// Load loads & parses all templates under the root.
// This should be called ASAP since it will cache all
// parsed templates & cause panic if there's any error occured.
func (self *Loader) Load() (pages int) {
	if !self.loaded {
		self.mutex.Lock()
		defer self.mutex.Unlock()
		for _, name := range self.Files() {
			self.templates[name] = self.page(name).parse()
			pages++
		}
		self.loaded = true
	}
	return
}

// internal page helper.
func (self *Loader) page(name string) *page {
	page := new(page)
	page.Name = name
	page.loader = self
	return page
}

// Reset clears the cached pages.
func (self *Loader) Reset() {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	for k := range self.templates {
		delete(self.templates, k)
	}
}

// Ancesters finds all ancestors absolute path using jinja's syntax
// and combines them along with the page name iteself into correct order for parsing.
// tag: {% extends "layout/base.html" %}
func (self *page) ancestors() (names []string) {
	var name = self.Name
	names = append(names, name)

	for {
		// find the very first "extends" tag.
		var bits, err = ioutil.ReadFile(path.Join(self.loader.root, name))
		if err != nil {
			log.Fatalf("Failed to open template (%s): %v", name, err)
		}

		var result = tags.Extends.FindSubmatch(bits)
		if result == nil {
			break
		}

		var base = string(result[1])
		if base == name {
			log.Fatalf("Template cannot extend itself (%s)", name)
		}

		names = append([]string{base}, names...) // insert the ancester into the first place.
		name = base
	}

	return
}

// Include finds all included external file sources recursively
// & replace all the "include" tags with their actual sources.
// tag: {% include "partials/header.html" %}
func (self *page) include() (source string) {
	bits, err := ioutil.ReadFile(self.path())
	if err != nil {
		log.Fatalf("Failed to open template (%s): %v", self.Name, err)
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
				log.Fatalf("Template cannot include itself (%s)", name)
			}
			page := self.loader.page(name)
			// reconstructs source to recursively find all included sources.
			source = strings.Replace(source, tag, page.source(), -1)
		}
	}
	return
}

// Parse constructs `template.Template` object with additional // "extends" & "include" like Jinja.
func (self *page) parse() *template.Template {
	var (
		err    error
		output *template.Template
	)
	names := self.ancestors()

	for _, name := range names {
		page := self.loader.page(name)
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

// Path returns the abolute path of the page.
func (self *page) path() string {
	return path.Join(self.loader.root, self.Name)
}

// Source returns the plain raw source of the page.
func (self *page) source() (src string) {
	if bits, err := ioutil.ReadFile(self.path()); err == nil {
		src = string(bits)
	} else {
		log.Fatalf("Failed to open template (%s): %v", self.Name, err)
	}
	return src
}

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
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var mutex sync.RWMutex

type Loader struct {
	root      string
	loaded    bool
	templates map[string]*template.Template
}

func NewLoader(path string) *Loader {
	abspath, err := filepath.Abs(path)
	if os.IsNotExist(err) {
		log.Fatalf("web/template: %s does not exist", path)
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
		log.Fatalf("web/template: files list cannot be listed: %v", err)
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
func (self *Loader) Load() {
	if !self.loaded {
		mutex.Lock()
		defer mutex.Unlock()
		for _, name := range self.Files() {
			self.templates[name] = self.page(name).parse()
		}
		self.loaded = true
	}
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
	mutex.Lock()
	defer mutex.Unlock()
	for k := range self.templates {
		delete(self.templates, k)
	}
}

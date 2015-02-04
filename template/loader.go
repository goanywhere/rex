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
package template

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/goanywhere/x/fs"
)

var ignores = regexp.MustCompile(`(include|layout)s?`)

type Loader struct {
	sync.RWMutex

	root      string
	loaded    bool
	templates map[string]*template.Template
}

func NewLoader(path string) *Loader {
	loader := new(Loader)
	loader.root, _ = filepath.Abs(path)
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

// Get retrieves the parsed template from preloaded pool.
func (self *Loader) Get(name string) (*template.Template, bool) {
	self.Load()
	template, exists := self.templates[name]
	return template, exists
}

// Load loads & parses all templates under the root.
// This should be called ASAP since it will cache all
// parsed templates & cause panic if there's any error occured.
func (self *Loader) Load() (pages int) {
	if fs.Exists(self.root) && !self.loaded {
		self.Lock()
		defer self.Unlock()

		err := filepath.Walk(self.root, func(path string, info os.FileInfo, err error) error {
			// NOTE ignore folders for partial HTMLs: layout(s) & include(s).
			if info.IsDir() && ignores.MatchString(info.Name()) {
				return filepath.SkipDir
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".html") {
				if name, e := filepath.Rel(self.root, path); e == nil {
					self.templates[name] = self.page(name).parse()
					pages++
				} else {
					err = e
				}
			}
			return err
		})
		if err != nil {
			log.Fatalf("Failed to list HTML templates: %v", err)
		}

		self.loaded = true
	}
	return
}

// internal page helper.
func (self *Loader) page(name string) *page {
	page := new(page)
	page.name = name
	page.loader = self
	return page
}

// Reset clears the cached pages.
func (self *Loader) Reset() {
	self.Lock()
	defer self.Unlock()
	for k := range self.templates {
		delete(self.templates, k)
	}
}

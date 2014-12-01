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
	"os"
	"path/filepath"
	"sync"

	"github.com/goanywhere/web"
)

var mutex sync.RWMutex

type Loader struct {
	root      string
	templates map[string]*template.Template
}

func NewLoader(path string) *Loader {
	mutex.Lock()
	defer mutex.Unlock()
	abspath, err := filepath.Abs(path)
	if os.IsNotExist(err) {
		web.Panic("web/template: %s does not exist", path)
	}
	loader := new(Loader)
	loader.root = abspath
	loader.templates = make(map[string]*template.Template)
	return loader
}

// Reset clears the cached pages.
func (self *Loader) Reset() {
	mutex.Lock()
	defer mutex.Unlock()
	for k := range self.templates {
		delete(self.templates, k)
	}
}

// page creates a internal page helper.
func (self *Loader) Page(name string) *Page {
	// setup constructs a new Page object.
	page := new(Page)
	page.Name = name
	page.loader = self
	return page
}

// Load returns cached page template.
func (self *Loader) Load(name string) *template.Template {
	mutex.Lock()
	defer mutex.Unlock()

	abspath := filepath.Join(self.root, name)
	if _, err := os.Stat(abspath); os.IsNotExist(err) {
		web.Panic("web/template: template does not exist (%s)", name)
	}

	if page, exists := self.templates[name]; exists {
		return page
	}
	self.templates[name] = self.Page(name).parse()
	return self.templates[name]
}

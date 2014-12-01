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
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var mutex sync.RWMutex

type Loader struct {
	root  string
	pages map[string]*Page
}

func NewLoader(path string) *Loader {
	mutex.Lock()
	defer mutex.Unlock()
	abspath, err := filepath.Abs(path)
	if os.IsNotExist(err) {
		panic(fmt.Errorf("web/template: %s does not exist", path))
	}
	loader := new(Loader)
	loader.root = abspath
	loader.pages = make(map[string]*Page)
	return loader
}

// setup constructs a new Page object.
func (self *Loader) setup(path string) *Page {
	page := new(Page)
	page.Name = path
	page.loader = self
	return page
}

// Reset clears the cached pages.
func (self *Loader) Reset() {
	mutex.Lock()
	defer mutex.Unlock()
	for k := range self.pages {
		delete(self.pages, k)
	}
}

// Load returns cached page template.
func (self *Loader) Load(path string) *Page {
	mutex.Lock()
	defer mutex.Unlock()

	abspath := filepath.Join(self.root, path)
	if _, err := os.Stat(abspath); os.IsNotExist(err) {
		panic(fmt.Errorf("web/template: template does not exist (%s)", abspath))
	}

	if page, exists := self.pages[path]; exists {
		return page
	}
	self.pages[path] = self.setup(path)
	return self.pages[path]
}

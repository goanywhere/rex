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
package internal

import (
	"net/http"

	"github.com/Sirupsen/logrus"
)

type Module struct {
	cache http.Handler
	stack []func(http.Handler) http.Handler
}

func (self *Module) build() http.Handler {
	if self.cache == nil {
		next := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {}))
		// Activate modules in FIFO order.
		for index := len(self.stack) - 1; index >= 0; index-- {
			next = self.stack[index](next)
		}
		self.cache = next
	}
	return self.cache
}

// Use add the middleware module into the stack chain.
// Supported middleware modules:
//  - http.Handler
//  - func(http.Handler) http.Handler
func (self *Module) Use(mod interface{}) {
	switch H := mod.(type) {
	case http.Handler:
		handler := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				H.ServeHTTP(w, req)
				next.ServeHTTP(w, req)
			})
		}
		self.stack = append(self.stack, handler)

	case func(http.Handler) http.Handler:
		self.stack = append(self.stack, H)

	default:
		logrus.Fatal("Unsupported middleware module passed in.")
	}
}

// Implements the net/http Handler interface and calls the middleware stack.
func (self *Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self.build().ServeHTTP(w, r)
}

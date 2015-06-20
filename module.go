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
package rex

import "net/http"

type module struct {
	stack []func(http.Handler) http.Handler
}

// build sets up the whole middleware modules in a FIFO chain.
func (self *module) build() http.Handler {
	var next http.Handler = http.DefaultServeMux
	// Activate modules in FIFO order.
	for index := len(self.stack) - 1; index >= 0; index-- {
		next = self.stack[index](next)
	}
	return next
}

// Use add the middleware module into the stack chain.
func (self *module) Use(modules ...func(http.Handler) http.Handler) {
	self.stack = append(self.stack, modules...)
}

// Implements the net/http Handler interface and calls the middleware stack.
func (self *module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self.build().ServeHTTP(w, r)
}

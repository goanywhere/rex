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

type middleware struct {
	cache http.Handler
	stack []func(http.Handler) http.Handler
}

// Implements the net/http Handler interface and calls the middleware stack.
func (self *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if self.cache == nil {
		// setup the whole middleware modules in a FIFO chain.
		var next http.Handler = http.DefaultServeMux
		for index := len(self.stack) - 1; index >= 0; index-- {
			next = self.stack[index](next)
		}
		self.cache = next
	}
	self.cache.ServeHTTP(w, r)
}

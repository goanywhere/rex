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
package filters

import "net/http"

const (
	xFrameOptions       = "X-Frame-Options"
	xContentTypeOptions = "X-Content-Type-Options"
	xXSSProtection      = "X-XSS-Protection"
	xUACompatible       = "X-UA-Compatible"
)

type header struct {
	writer http.ResponseWriter
}

func (self *header) set(key string, value interface{}) {
	if v := self.writer.Header().Get(key); v == "" {
		if value != nil {
			self.writer.Header()[key] = []string{value.(string)}
		}
	}
}

// Header provides additional headers supports for response writer.
func Header(options Options) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var header = &header{w}
			header.set(xFrameOptions, options.Get(xFrameOptions, settings.X_Frame_Options))
			header.set(xContentTypeOptions, options.Get(xContentTypeOptions, settings.X_Content_Type_Options))
			header.set(xXSSProtection, options.Get(xXSSProtection, settings.X_XSS_Protection))
			header.set(xUACompatible, options.Get(xUACompatible, settings.X_UA_Compatible))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

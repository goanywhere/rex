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
package modules

import "net/http"

const (
	xFrameOptions       = "X-Frame-Options"
	xContentTypeOptions = "X-Content-Type-Options"
	xXSSProtection      = "X-XSS-Protection"
	xUACompatible       = "X-UA-Compatible"
)

func securelySetHeader(w http.ResponseWriter, key string, value interface{}) {
	if v := w.Header().Get(key); v == "" {
		if value != nil {
			w.Header().Set(key, value.(string))
		}
	}
}

// Security provides basic security supports by rendering common response headers.
func Security(options Options) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			securelySetHeader(w, xFrameOptions, options.Get("X-Frame-Options", settings.XFrameOptions))
			securelySetHeader(w, xContentTypeOptions, options.Get("X-Content-Type-Options", settings.XContentTypeOptions))
			securelySetHeader(w, xXSSProtection, options.Get("X-XSS-Protection", settings.XXSSProtection))
			securelySetHeader(w, xUACompatible, options.Get("X-UA-Compatible", settings.XUACompatible))

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

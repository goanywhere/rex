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

import (
	"net/http"
	"path"
	"strings"
)

// Static serves as file server for static assets,
// as convention, the given dir name will be used as the URL prefix.
func Static(dir string) func(http.Handler) http.Handler {
	var (
		fs     = http.Dir(dir)
		prefix = path.Join("/", dir)
	)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Accepts http GET | HEAD Only, ignores all requests not started with prefix.
			if r.Method != "GET" && r.Method != "HEAD" {
				next.ServeHTTP(w, r)
				return
			} else if !strings.HasPrefix(r.URL.Path, prefix) {
				next.ServeHTTP(w, r)
				return
			}

			filename := strings.TrimPrefix(r.URL.Path, prefix)
			if filename != "" && filename[0] != '/' {
				next.ServeHTTP(w, r)
				return
			}

			file, err := fs.Open(filename)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			defer file.Close()

			stat, err := file.Stat()
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			// try to serve index filename
			if stat.IsDir() {
				// redirect if missing trailing slash
				if !strings.HasSuffix(r.URL.Path, "/") {
					http.Redirect(w, r, r.URL.Path+"/", http.StatusFound)
					return
				}

				filename = path.Join(filename, "index.html")
				file, err = fs.Open(filename)
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}
				defer file.Close()

				stat, err = file.Stat()
				if err != nil || stat.IsDir() {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.ServeContent(w, r, filename, stat.ModTime(), file)
		})
	}
}

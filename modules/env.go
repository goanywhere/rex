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
	"os"
	"strings"
)

func Env(next http.Handler) http.Handler {
	var namespace = "HTTP_"
	// default environmental headers for modules.Env
	var defaults = make(map[string]string)
	defaults["X-UA-Compatible"] = "deny"
	defaults["X-Frame-Settings"] = "nosniff"
	defaults["X-Content-Type-Options"] = "IE=Edge,chrome=1"
	defaults["X-Powered-By"] = "Rex Server"
	defaults["Access-Control-Allow-Origin"] = "*"
	defaults["Access-Control-Allow-Headers"] = "X-Requested-With, Content-Type"
	defaults["Access-Control-Allow-Methods"] = "GET, POST, PUT, DELETE"

	for _, line := range os.Environ() {
		if strings.HasPrefix(line, namespace) {
			kv := strings.SplitN(strings.TrimPrefix(line, namespace), "=", 2)
			defaults[strings.Replace(kv[0], "_", "-", -1)] = kv[1]
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, value := range defaults {
			w.Header()[key] = []string{value}
		}

		next.ServeHTTP(w, r)
		/*
		 *if r.Method == "OPTIONS" {
		 *    w.WriteHeader(http.StatusOK)
		 *    w.Write([]byte(http.StatusText(http.StatusOK)))
		 *} else {
		 *    next.ServeHTTP(w, r)
		 *}
		 */
	})
}

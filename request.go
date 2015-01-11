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

import (
	"net/http"
	"regexp"
	"strings"
)

var (
	regexAcceptEncoding = regexp.MustCompile(`(\w+|\*)(;q=(1(\.0)?|0(\.[0-9])?))?`)
)

// AcceptEncodings fetches the requested encodings from client with priority.
func AcceptEncodings(request *http.Request) (encodings []string) {
	// find all encodings supported by backend server.
	matches := regexAcceptEncoding.FindAllString(request.Header.Get("Accept-Encoding"), -1)
	for _, item := range matches {
		units := strings.SplitN(item, ";", 2)
		// top priority with q=1|q=1.0|Not Specified.
		if len(units) == 1 {
			encodings = append(encodings, units[0])

		} else {
			if strings.HasPrefix(units[1], "q=1") {
				// insert the specified top priority to the first.
				encodings = append([]string{units[0]}, encodings...)

			} else if strings.HasSuffix(units[1], "0") {
				// not acceptable at client side.
				continue
			} else {
				// lower priority encoding
				encodings = append(encodings, units[0])
			}
		}
	}
	return
}

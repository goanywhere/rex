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
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/goanywhere/rex/modules/livereload"
)

type writer struct {
	http.ResponseWriter
	host string
}

func (self *writer) addJavaScript(data []byte) []byte {
	javascript := fmt.Sprintf(`<script src="//%s%s"></script>
</head>`, self.host, livereload.URL.JavaScript)
	return regexp.MustCompile(`</head>`).ReplaceAll(data, []byte(javascript))
}

func (self *writer) Write(data []byte) (size int, e error) {
	if strings.Contains(self.Header().Get("Content-Type"), "html") {
		var encoding = self.Header().Get("Content-Encoding")
		if encoding == "" {
			data = self.addJavaScript(data)
		} else {
			var reader io.ReadCloser
			var buffer *bytes.Buffer = new(bytes.Buffer)

			if encoding == "gzip" {
				// decode to add javascript reference.
				reader, _ = gzip.NewReader(bytes.NewReader(data))
				io.Copy(buffer, reader)
				output := self.addJavaScript(buffer.Bytes())
				reader.Close()
				buffer.Reset()
				// encode back to HTML with added javascript reference.
				writer := gzip.NewWriter(buffer)
				writer.Write(output)
				writer.Close()
				data = buffer.Bytes()

			} else if encoding == "deflate" {
				// decode to add javascript reference.
				reader, _ = zlib.NewReader(bytes.NewReader(data))
				io.Copy(buffer, reader)
				output := self.addJavaScript(buffer.Bytes())
				reader.Close()
				buffer.Reset()
				// encode back to HTML with added javascript reference.
				writer := zlib.NewWriter(buffer)
				writer.Write(output)
				writer.Close()
				data = buffer.Bytes()
			}
		}
	}
	return self.ResponseWriter.Write(data)
}

func LiveReload(next http.Handler) http.Handler {
	livereload.Start()
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == livereload.URL.WebSocket {
			livereload.ServeWebSocket(w, r)
		} else if r.URL.Path == livereload.URL.JavaScript {
			livereload.ServeJavaScript(w, r)
		} else {
			writer := &writer{w, r.Host}
			next.ServeHTTP(writer, r)
		}
	}
	return http.HandlerFunc(fn)
}

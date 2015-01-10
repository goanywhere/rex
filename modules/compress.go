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
	"compress/flate"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/goanywhere/rex/web"
)

var (
	regexContentType = regexp.MustCompile(`((message|text)\/.+)|((application\/).*(javascript|json|xml))`)
)

type compression interface {
	io.Writer
	io.Closer
}

type compressor struct {
	http.ResponseWriter
	encodings []string
}

func (self *compressor) filter(src []byte) ([]byte, string) {
	var mimetype = self.Header().Get("Content-Type")
	if mimetype == "" {
		mimetype = http.DetectContentType(src)
		self.Header().Set("Content-Type", mimetype)
	}

	if self.Header().Get("Content-Encoding") != "" {
		return src, ""
	}

	if !regexContentType.MatchString(strings.TrimSpace(strings.SplitN(mimetype, ";", 2)[0])) {
		return src, ""
	}

	// okay to start compressing.
	var e error
	var encoding string
	var writer compression
	var buffer *bytes.Buffer = new(bytes.Buffer)
	// try compress the data, if any error occrued, fallback to ResponseWriter.
	if self.encodings[0] == "deflate" {
		encoding = "deflate"
		writer, e = flate.NewWriter(buffer, flate.DefaultCompression)
	} else {
		encoding = "gzip"
		writer, e = gzip.NewWriterLevel(buffer, gzip.DefaultCompression)
	}
	_, e = writer.Write(src)
	writer.Close()
	if e == nil {
		return buffer.Bytes(), encoding
	}
	// fallback to standard http.ResponseWriter, nothing happened~ (~__~"")
	return src, ""
}

func (self *compressor) Write(data []byte) (size int, err error) {
	if bytes, encoding := self.filter(data); encoding != "" {
		log.Printf("Compress: %s", self.Header().Get("Content-Type"))
		self.Header().Set("Content-Encoding", encoding)
		self.Header().Add("Vary", "Accept-Encoding")
		self.Header().Del("Content-Length")
		return self.ResponseWriter.Write(bytes)
	}
	return self.ResponseWriter.Write(data)
}

func Compress(options Options) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Sec-WebSocket-Key") != "" || r.Method == "HEAD" {
				next.ServeHTTP(w, r)
			} else {
				encodings := web.AcceptEncodings(r)
				if len(encodings) == 0 {
					next.ServeHTTP(w, r)

				} else {
					compressor := new(compressor)
					compressor.ResponseWriter = w
					compressor.encodings = encodings
					next.ServeHTTP(compressor, r)
				}
			}
		}
		return http.HandlerFunc(fn)
	}
}

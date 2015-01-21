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
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type Context struct {
	server  *Server
	Request *http.Request
	Writer  http.ResponseWriter
}

// Cookie returns the cookie value previously set.
func (self *Context) Cookie(key string) (value string) {
	if cookie, err := self.Request.Cookie(key); err == nil {
		if cookie.Value != "" {
			value = cookie.Value
		}
	}
	return
}

// SetCookie writes cookie to ResponseWriter.
func (self *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(self.Writer, cookie)
}

// IsAjax checks if the incoming request is AJAX request.
func (self *Context) IsAjax() bool {
	return self.Request.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

// Query returns the URL query values.
func (self *Context) Query() url.Values {
	return self.Request.URL.Query()
}

// Error raises a HTTP error response according to the given status code.
func (self *Context) Error(status int, errors ...string) {
	if len(errors) == 1 {
		http.Error(self.Writer, errors[0], status)
	} else {
		http.Error(self.Writer, http.StatusText(status), status)
	}
}

// HTML renders cached HTML templates via `bytes.Buffer` to response.
func (self *Context) HTML(filename string, data ...map[string]interface{}) {
	var buffer = new(bytes.Buffer)
	self.Writer.Header()["Content-Type"] = []string{"text/html; charset=utf-8"}
	template := self.server.loader.Get(filename)

	var v map[string]interface{}
	if len(data) > 0 {
		v = data[0]
	}
	if err := template.Execute(buffer, v); err != nil {
		self.Error(http.StatusInternalServerError)
		return
	}
	self.Writer.Write(buffer.Bytes())
}

// JSON renders JSON data to response.
func (self *Context) JSON(v interface{}) {
	if data, e := json.Marshal(v); e == nil {
		self.Writer.Header()["Content-Type"] = []string{"application/json; charset=utf-8"}
		self.Writer.Write(data)
	} else {
		log.Printf("Failed to render JSON: %v", e)
		self.Error(http.StatusInternalServerError, e.Error())
	}
}

// String writes plain text back to the HTTP response.
func (self *Context) String(format string, values ...interface{}) {
	self.Writer.Header()["Content-Type"] = []string{"text/plain; charset=utf-8"}
	self.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// XML renders XML data to response.
func (self *Context) XML(v interface{}) {
	if data, e := xml.Marshal(v); e == nil {
		self.Writer.Header()["Content-Type"] = []string{"application/xml; charset=utf-8"}
		self.Writer.Write(data)
	} else {
		log.Printf("Failed to render XML: %v", e)
		self.Error(http.StatusInternalServerError, e.Error())
	}
}

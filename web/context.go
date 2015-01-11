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
package web

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
	Writer
	Request *http.Request

	data map[string]interface{}
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := new(Context)
	ctx.data = make(map[string]interface{})

	ctx.Request = r
	ctx.Writer = &writer{w, 0, 0}
	return ctx
}

// Get fetches context data under the given key.
func (self *Context) Get(key string) interface{} {
	if value, exists := self.data[key]; exists {
		return value
	}
	return nil
}

// Set write value under the given key to context data.
func (self *Context) Set(key string, value interface{}) {
	self.data[key] = value
}

// Clear wipes out all existing context data.
func (self *Context) Clear() {
	for key := range self.data {
		delete(self.data, key)
	}
}

// Delete removes context data under the given key.
func (self *Context) Delete(key string) {
	delete(self.data, key)
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

// SecureCookie decodes the signed value from cookie.
// Empty string value will be returned if the signature is invalide or expired.
func (self *Context) SecureCookie(key string) (value string) {
	if src := self.Cookie(key); src != "" {
		if bits, err := signature.Decode(key, src); err == nil {
			value = string(bits)
		}
	}
	return
}

// SetSecureCookie replaces the raw value with a signed one & write the cookie into Context.
func (self *Context) SetSecureCookie(cookie *http.Cookie) {
	if cookie.Value != "" {
		if value, err := signature.Encode(cookie.Name, []byte(cookie.Value)); err == nil {
			cookie.Value = value
		}
	}
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
func (self *Context) Error(status int) {
	http.Error(self.Writer, http.StatusText(status), status)
}

// HTML renders cached HTML templates via `bytes.Buffer` to response.
// Under Debug mode, livereload.js will be added to the end of <head>
// to provide browser-based LiveReload supports.
func (self *Context) HTML(filename string) {
	var buffer bytes.Buffer
	self.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := loader.Get(filename).Execute(&buffer, self.data); err != nil {
		self.Error(http.StatusInternalServerError)
		return
	}
	self.Writer.Write(buffer.Bytes())
}

// JSON renders JSON data to response.
func (self *Context) JSON(values map[string]interface{}) {
	var (
		data []byte
		err  error
	)
	data, err = json.Marshal(values)
	if err != nil {
		log.Printf("Failed to render JSON: %v", err)
		self.Error(http.StatusInternalServerError)
		return
	}
	self.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	self.Writer.Write(data)
}

// String writes plain text back to the HTTP response.
func (self *Context) String(format string, values ...interface{}) {
	self.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	self.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// XML renders XML data to response.
func (self *Context) XML(values interface{}) {
	self.Writer.Header().Set("Content-Type", "application/xml; charset=utf-8")
	xml.NewEncoder(self.Writer).Encode(values)
}

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
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/securecookie"

	. "github.com/goanywhere/rex/internal"
)

type Context struct {
	http.ResponseWriter
	Request *http.Request
	Layout  string

	size   int
	status int

	session Session
	buffer  *bytes.Buffer
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := new(Context)
	ctx.ResponseWriter = w
	ctx.Request = r
	ctx.buffer = new(bytes.Buffer)
	return ctx
}

// ----------------------------------------
// Context Values for Rendering
// ----------------------------------------
// Clear removes all values stored for the request.
func (self *Context) Clear() {
	context.Clear(self.Request)
}

// Del removes value stored with associated key for the request.
func (self *Context) Del(key string) {
	context.Delete(self.Request, key)
}

// Get retrieves value stored for the request with associated key.
func (self *Context) Get(key string) interface{} {
	return context.Get(self.Request, key)
}

// Set stores value with associated key for the request.
func (self *Context) Set(key string, value interface{}) {
	context.Set(self.Request, key, value)
}

// ----------------------------------------
// Session & (Secure) Cookie Supports
// ----------------------------------------
// Session fetches the securecookie based session from incoming request.
func (self *Context) Session() Session {
	if self.session == nil {
		name := settings.String("SESSION_COOKIE_NAME")
		session := &session{
			ctx:    self,
			values: make(map[string]interface{}),
		}
		self.SecureCookie(name, &session.values)
		self.session = session
	}
	return self.session
}

// Cookie fetches the value associated with the given name from request cookie.
func (self *Context) Cookie(name string) string {
	if cookie, err := self.Request.Cookie(name); err == nil {
		return cookie.Value
	}
	return ""
}

// SetCookie adds a Set-Cookie header to response with default options.
func (self *Context) SetCookie(name, value string, options ...*http.Cookie) {
	var cookie *http.Cookie
	if len(options) > 0 {
		cookie = options[0]
	} else {
		cookie = new(http.Cookie)
		cookie.Path = settings.String("SESSION_COOKIE_PATH")
		cookie.Domain = settings.String("SESSION_COOKIE_DOMAIN")
		cookie.MaxAge = settings.Int("SESSION_COOKIE_MAXAGE")
		cookie.Secure = settings.Bool("SESSION_COOKIE_SECURE")
		cookie.HttpOnly = settings.Bool("SESSION_COOKIE_HTTPONLY")
	}
	cookie.Name = name
	cookie.Value = value
	// IE 6/7/8 Compatible Mode.
	if cookie.MaxAge > 0 {
		cookie.Expires = time.Now().Add(time.Duration(cookie.MaxAge) * time.Second)
	} else if cookie.MaxAge < 0 {
		cookie.Expires = time.Unix(1, 0)
	}
	http.SetCookie(self, cookie)
}

// SecureCookie decodes the signed values associated with the given name from request cookie.
func (self *Context) SecureCookie(name string, ptr interface{}) error {
	if raw := self.Cookie(name); raw != "" {
		return securecookie.DecodeMulti(name, raw, ptr, secrets...)
	}
	return nil
}

// SetSecureCookie encode the raw value securely and adds a Set-Cookie header to response.
func (self *Context) SetSecureCookie(name string, value interface{}, options ...*http.Cookie) error {
	if raw, err := securecookie.EncodeMulti(name, value, secrets...); err == nil {
		self.SetCookie(name, raw, options...)
	} else {
		return err
	}
	return nil
}

// ----------------------------------------
// HTTP Utilities
// ----------------------------------------
// Query returns the URL query values.
func (self *Context) Query() url.Values {
	return self.Request.URL.Query()
}

// Error raises a HTTP error response according to the given status code.
func (self *Context) Error(status int, errors ...string) {
	if len(errors) > 0 {
		http.Error(self, errors[0], status)
	} else {
		http.Error(self, http.StatusText(status), status)
	}
}

// Flush sends any buffered data to the client & clear all context values.
func (self *Context) Flush() {
	if self.buffer.Len() > 0 {
		self.Layout = ""
		self.Write(self.buffer.Bytes())
		self.buffer.Reset()
	}
}

// RemoteAddr fetches the real remote address of incoming HTTP request.
func (self *Context) RemoteAddr() string {
	var address string

	if raw := self.Request.Header.Get("X-Forwarded-For"); raw != "" {
		index := strings.Index(raw, ", ")
		if index == -1 {
			index = len(raw)
		}
		address = raw[:index]

	} else if raw := self.Request.Header.Get("X-Real-IP"); raw != "" {
		address = raw
	}

	return address
}

// Render constructs final response data to the client side.
// If Context.Layout is given, HTML page will be constructed
// using built-in html/template package, data will be encoded
// using JSON (default) | XML (Specified by ContentType Header)
// encoder otherwise.
func (self *Context) Render(v interface{}) {
	var e error

	if self.Layout == "" {
		if ctype := self.Header().Get(ContentType.Name); strings.Contains(ctype, "xml") {
			e = xml.NewEncoder(self.buffer).Encode(v)

		} else {
			if ctype == "" {
				self.Header().Set(ContentType.Name, ContentType.JSON)
			}
			e = json.NewEncoder(self.buffer).Encode(v)
		}

	} else {
		if template, exists := templates.Get(self.Layout); exists {
			e = template.Execute(self.buffer, v)
		} else {
			e = fmt.Errorf("Template <%s> does not exists", self.Layout)
		}
	}

	if e == nil {
		self.Flush()
	} else {
		self.Error(http.StatusInternalServerError, e.Error())
	}
}

// Write writes the data to the connection as part of an HTTP reply.
// If WriteHeader has not yet been called, Write calls WriteHeader(http.StatusOK)
// before writing the data.  If the Header does not contain a
// Content-Type line, Write adds a Content-Type set to the result of passing
// the initial 512 bytes of written data to DetectContentType.
func (self *Context) Write(bytes []byte) (size int, err error) {
	size, err = self.ResponseWriter.Write(bytes)
	self.size += size
	return
}

// WriteHeader sends an HTTP response header with status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
func (self *Context) WriteHeader(status int) {
	if status >= 100 && status < 512 {
		self.status = status
		self.ResponseWriter.WriteHeader(status)
	}
}

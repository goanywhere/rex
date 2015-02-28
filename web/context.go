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
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/gorilla/securecookie"

	. "github.com/goanywhere/rex/internal"
)

type Context struct {
	http.ResponseWriter
	Request *http.Request

	size   int
	status int

	session Session
	buffer  *bytes.Buffer
	values  map[string]interface{}
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := new(Context)
	ctx.ResponseWriter = w
	ctx.Request = r
	ctx.buffer = new(bytes.Buffer)
	ctx.values = make(map[string]interface{})
	return ctx
}

// ----------------------------------------
// Context Values for Rendering
// ----------------------------------------
// Clear removes all context values.
func (self *Context) Clear() {
	for k, _ := range self.values {
		delete(self.values, k)
	}
}

// Delete removes a context value assoicated with the given key.
func (self *Context) Del(key string) {
	delete(self.values, key)
}

// Get retrieves the value with associated key from context values.
func (self *Context) Get(key string, ptr interface{}) {
	if value := self.values[key]; value != nil {
		if reflect.TypeOf(ptr).Kind() == reflect.Ptr {
			elem := reflect.ValueOf(ptr).Elem()
			elem.Set(reflect.ValueOf(value))
		}
	}
}

// Set adds the value with associated key into context to be rendered.
func (self *Context) Set(key string, value interface{}) {
	self.values[key] = value
}

// ----------------------------------------
// Session & (Secure) Cookie Supports
// ----------------------------------------
// Session fetches the securecookie based session from incoming request.
func (self *Context) Session() Session {
	if self.session == nil {
		name := settings.String("session.cookie.name")
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
		cookie.Path = settings.String("session.cookie.path")
		cookie.Domain = settings.String("session.cookie.domain")
		cookie.MaxAge = settings.Int("session.cookie.maxage")
		cookie.Secure = settings.Bool("session.cookie.secure")
		cookie.HttpOnly = settings.Bool("session.cookie.httponly")
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

// // Flush sends any buffered data to the client & clear all context values.
func (self *Context) Flush() {
	if self.buffer.Len() > 0 {
		self.Write(self.buffer.Bytes())
		self.Reset()
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

// Reset clears context's buffer & values.
func (self *Context) Reset() {
	self.buffer.Reset()
	self.Clear()
}

// Render constructs the final output using html/template & json.
// ContentType is determined by extension of the given object.
// if the object if string, context try parsing its extension
// to determined the content type (HTML|JSON|XML), otherwise, the basic
//	format:
//		{"status": <status code>, "data": <response data>}
// json dict will be used, this is to ensure the json output is always
// a dictionary due to security concern.
func (self *Context) Render(object interface{}) {
	var err error

	switch T := object.(type) {
	case string:
		switch filepath.Ext(T) {
		case ".html":
			self.Header().Set(ContentType.Name, ContentType.HTML)

		case ".json":
			self.Header().Set(ContentType.Name, ContentType.JSON)

		case ".xml":
			self.Header().Set(ContentType.Name, ContentType.XML)

		default:
			log.Fatalf("Unsupported file type: %s", T)
		}
		if document, exists := documents.Get(T); exists {
			err = document.Execute(self.buffer, self.values)
		} else {
			err = fmt.Errorf("Template <%s> does not exists", T)
		}

	default:
		self.Header().Set(ContentType.Name, ContentType.JSON)
		if self.status == 0 {
			self.values["status"] = http.StatusOK
		} else {
			self.values["status"] = self.status
		}
		self.values["data"] = object
		err = json.NewEncoder(self.buffer).Encode(self.values)
	}

	if err == nil {
		self.Flush()
	} else {
		self.Error(http.StatusInternalServerError, err.Error())
		self.Reset()
	}
}

func (self *Context) String(format string, values ...interface{}) {
	self.Header().Set(ContentType.Name, ContentType.Text)
	self.ResponseWriter.Write([]byte(fmt.Sprintf(format, values...)))
}

// Header returns the header map that will be sent by WriteHeader.
// Changing the header after a call to WriteHeader (or Write) has
// no effect.
func (self *Context) Header() http.Header {
	return self.ResponseWriter.Header()
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

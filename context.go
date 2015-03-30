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
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	pongo "github.com/flosch/pongo2"
	gschema "github.com/gorilla/schema"

	"github.com/goanywhere/x"
	"github.com/goanywhere/x/env"
)

var schema = gschema.NewDecoder()

type Validator interface {
	Validate() error
}

type Context struct {
	http.ResponseWriter
	Request *http.Request

	size   int
	status int

	values pongo.Context
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := new(Context)
	ctx.ResponseWriter = w
	ctx.Request = r
	ctx.values = pongo.Context{}
	return ctx
}

// ----------------------------------------
// Context Values for Rendering
// ----------------------------------------
// Get retrieves the value with associated key from the request context.
func (self *Context) Get(key string) interface{} {
	return self.values[key]
}

// Set stores value with associated key for the request context.
func (self *Context) Set(key string, value interface{}) {
	self.values[key] = value
}

// Del removes the value with associated key for request context.
func (self *Context) Del(key string) {
	delete(self.values, key)
}

// ----------------------------------------
// Secure Cookie Supports
// ----------------------------------------
// Cookie securely fetches the value associated with the given name from request cookie.
func (self *Context) Cookie(name string, v interface{}) error {
	if securecookie == nil {
		return errors.New("env:SECRET_KEYS is missing")
	}

	if cookie, err := self.Request.Cookie(name); err == nil {

		if cookie.Value != "" {
			return securecookie.Decode(name, cookie.Value, v)
		}

	} else {
		return err
	}
	return nil
}

// SetCookie securely encodes the value to response cookie via Set-Cookie header.
func (self *Context) SetCookie(name string, value interface{}, options ...*http.Cookie) error {
	if securecookie == nil {
		return errors.New("env:SECRET_KEYS is missing")
	}

	if str, err := securecookie.Encode(name, value); err == nil {

		var cookie *http.Cookie

		if len(options) > 0 {
			cookie = options[0]
		} else {
			cookie = &http.Cookie{
				Path:     env.String("COOKIE_PATH", "/"),
				Domain:   env.String("COOKIE_DOMAIN", ""),
				Secure:   env.Bool("COOKIE_SECURE", false),
				HttpOnly: env.Bool("COOKIE_HTTPONLY", true),
				MaxAge:   env.Int("COOKIE_MAGAGE", 3600*24*7),
			}
		}
		cookie.Name = name
		cookie.Value = str

		// IE 6/7/8 Compatible Mode.
		if cookie.MaxAge > 0 {
			cookie.Expires = time.Now().Add(time.Duration(cookie.MaxAge) * time.Second)

		} else if cookie.MaxAge < 0 {
			cookie.Expires = time.Unix(1, 0)
		}

		http.SetCookie(self, cookie)

	} else {
		return err
	}
	return nil
}

// ----------------------------------------
// HTTP Utilities
// ----------------------------------------
//Error writes the error message to http response.
func (self *Context) Error(status int, message ...string) {
	if len(message) == 0 {
		http.Error(self, http.StatusText(status), status)
	} else {
		http.Error(self, message[0], status)
	}
}

// ParseForm parsed the raw query from the URL and updates request.Form,
// decode the from to the given struct with Validator implemented.
func (self *Context) ParseForm(form Validator) (err error) {
	if err = self.Request.ParseForm(); err == nil {
		if err = schema.Decode(form, self.Request.Form); err == nil {
			err = form.Validate()
		}
	}
	return
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

// Render constructs response body using html/template.
func (self *Context) Render(filename string) {
	defer func() {
		self.values = pongo.Context{}
	}()

	if self.Header().Get("Content-Type") == "" {
		self.Header().Set("Content-Type", "text/html; charset=UTF-8")
	}

	if template, exists := views[filename]; exists {
		if out, err := template.ExecuteBytes(self.values); err == nil {
			self.Write(out)
		} else {
			http.Error(self, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(self, fmt.Sprintf("<Template: %s> does not exists", filename), http.StatusInternalServerError)
	}
}

// Send constructs response body using the given object, Content-Type
// will be automatically detected using StructTag for objects.
// Supported Objects:
//      string: text/plain
//      json:   application/json
//      xml:    application/xml
// TODO performance boost by determining the Response ContentType using request header.
func (self *Context) Send(v interface{}) {
	switch T := v.(type) {
	case string:
		self.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		self.Write([]byte(T))

	default:
		var (
			e      error
			buffer = new(bytes.Buffer)
			ctype  = self.Request.Header.Get("Content-Type")
		)

		// JSON/XML rendering.
		if strings.Contains(ctype, "application/json") || x.Has(v, "json") {
			self.Header().Set("Content-Type", "application/json; charset=UTF-8")
			e = json.NewEncoder(buffer).Encode(v)

		} else {
			self.Header().Set("Content-Type", "application/xml; charset=UTF-8")
			self.Write([]byte(xml.Header))
			e = xml.NewEncoder(buffer).Encode(v)
		}

		if e == nil {
			self.Write(buffer.Bytes())
		} else {
			http.Error(self, e.Error(), http.StatusInternalServerError)
		}
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

// Var returns the route variables for the current request.
func (self *Context) Var(key string) (value string) {
	if vars := mux.Vars(self.Request); vars != nil {
		if v, exists := vars[key]; exists {
			value = v
		}
	}
	return
}

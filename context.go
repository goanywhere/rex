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
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/securecookie"

	"github.com/goanywhere/rex/internal"
)

var cookie = internal.Cookie{
	Path:     options.String("session.cookie.path"),
	Domain:   options.String("session.cookie.domain"),
	MaxAge:   options.Int("session.cookie.maxage"),
	Secure:   options.Bool("session.cookie.secure"),
	HttpOnly: options.Bool("session.cookie.httponly"),
}

type Context struct {
	Options internal.Cookie
	Writer  http.ResponseWriter
	Request *http.Request

	size   int
	status int

	dirty  bool
	values map[string]interface{}
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := new(Context)
	ctx.Writer = w
	ctx.Request = r
	ctx.Options = cookie
	return ctx
}

// ----------------------------------------
// Session Supports
// ----------------------------------------
// Fetches the secured cookie session from http request.
func (self *Context) session() map[string]interface{} {
	if self.values == nil {
		self.values = make(map[string]interface{})
		name := options.String("session.cookie.name")
		values := self.Cookie(name)
		if values != "" {
			securecookie.DecodeMulti(name, values, &self.values, app.codecs...)
		}
	}
	return self.values
}

// Get value associated with the given key from encrypted cookie.
func (self *Context) Get(key string) interface{} {
	return self.session()[key]
}

// Set the key/value into the encrypted cookie.
func (self *Context) Set(key string, value interface{}) {
	self.session()[key] = value
	self.dirty = true
}

// Save the key/value into the encrypted cookie.
func (self *Context) Save() error {
	if !self.dirty {
		return nil
	}
	name := options.String("session.cookie.name")
	values, err := securecookie.EncodeMulti(name, self.values, app.codecs...)
	if err != nil {
		return err
	}
	self.SetCookie(name, values)
	self.dirty = false
	return nil
}

// Cookie returns the cookie value previously set.
func (self *Context) Cookie(name string) string {
	if cookie, err := self.Request.Cookie(name); err == nil {
		return cookie.Value
	}
	return ""
}

// SetCookie writes cookie to ResponseWriter.
func (self *Context) SetCookie(key, value string) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    value,
		Path:     self.Options.Path,
		Domain:   self.Options.Domain,
		MaxAge:   self.Options.MaxAge,
		Secure:   self.Options.Secure,
		HttpOnly: self.Options.HttpOnly,
	}
	if cookie.MaxAge > 0 {
		cookie.Expires = time.Now().Add(time.Duration(cookie.MaxAge) * time.Second)
	} else if cookie.MaxAge < 0 {
		cookie.Expires = time.Unix(1, 0)
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
func (self *Context) Error(status int, errors ...string) {
	if len(errors) > 0 {
		http.Error(self.Writer, errors[0], status)
	} else {
		http.Error(self.Writer, http.StatusText(status), status)
	}
}

// HTML renders cached HTML templates via `bytes.Buffer` to response.
// TODO empty loader/html
func (self *Context) HTML(filename string, data ...map[string]interface{}) {
	var buffer = new(bytes.Buffer)
	self.Writer.Header()["Content-Type"] = []string{"text/html; charset=utf-8"}

	var err error
	if len(data) == 0 {
		err = app.HTML.Get(filename).Execute(buffer, nil)
	} else {
		err = app.HTML.Get(filename).Execute(buffer, data[0])
	}

	if err != nil {
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

/* ----------------------------------------------------------------------
 * Implementations of http.Hijacker
 * ----------------------------------------------------------------------*/
func (self *Context) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := self.Writer.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}

/* ----------------------------------------------------------------------
 * Implementations of http.CloseNotifier
 * ----------------------------------------------------------------------*/
func (self *Context) CloseNotify() <-chan bool {
	return self.Writer.(http.CloseNotifier).CloseNotify()
}

/* ----------------------------------------------------------------------
 * Implementations of http.Flusher
 * ----------------------------------------------------------------------*/
func (self *Context) Flush() {
	flusher, ok := self.Writer.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}

/* ----------------------------------------------------------------------
 * Implementations of http.ResponseWriter
 * ----------------------------------------------------------------------*/
func (self *Context) WriteHeader(status int) {
	if status >= 100 && status < 512 {
		self.status = status
		self.Writer.WriteHeader(status)
	}
}

// Write: Implementation of http.ResponseWriter#Write
func (self *Context) Write(data []byte) (size int, err error) {
	size, err = self.Writer.Write(data)
	self.size += size
	return
}

/* ----------------------------------------------------------------------
 * Implementations of rex.Writer interface.
 * ----------------------------------------------------------------------*/
func (self *Context) Size() int {
	return self.size
}

// Status returns current status code of the Context.
func (self *Context) Status() int {
	return self.status
}

func (self *Context) Written() bool {
	return self.status != 0 || self.size > 0
}

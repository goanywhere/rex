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

	session *session
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
func (self *Context) Session() *session {
	if self.session == nil {
		self.session = new(session)
		self.session.values = make(map[string]interface{})
		name := options.String("session.cookie.name")
		if raw := self.Cookie(name); raw != "" {
			securecookie.DecodeMulti(name, raw, &self.session.values, app.codecs...)
		}
	}
	return self.session
}

// ----------------------------------------
// Cookie & Secure Cookie Supports
// ----------------------------------------

// Cookie fetches the value associated with the given name from request cookie.
func (self *Context) Cookie(name string) string {
	if cookie, err := self.Request.Cookie(name); err == nil {
		return cookie.Value
	}
	return ""
}

// SetCookie adds a Set-Cookie header to response with default options.
func (self *Context) SetCookie(name, value string) {
	cookie := &http.Cookie{
		Name:     name,
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

// SignedCookie decodes the signed values associated with the given name from request cookie.
func (self *Context) SignedCookie(name string, value interface{}) error {
	if raw := self.Cookie(name); raw != "" {
		return securecookie.DecodeMulti(name, raw, value, app.codecs...)
	}
	return nil
}

// SetSignedCookie encode the raw value securely and adds a Set-Cookie header to response.
func (self *Context) SetSignedCookie(name string, value interface{}) error {
	if raw, err := securecookie.EncodeMulti(name, value, app.codecs...); err == nil {
		cookie := &http.Cookie{
			Name:  name,
			Value: raw,
		}
		http.SetCookie(self.Writer, cookie)
	} else {
		return err
	}
	return nil
}

// ----------------------------------------
// HTTP Utilities
// ----------------------------------------

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
func (self *Context) Render(filename string, values ...map[string]interface{}) {
	if template, exists := app.HTML.Get(filename); exists {
		var (
			buffer = new(bytes.Buffer)
			data   map[string]interface{}
		)
		if len(data) > 0 {
			data = values[0]
		}
		if e := template.Execute(buffer, data); e != nil {
			self.Error(http.StatusInternalServerError, e.Error())
			return
		}
		self.Writer.Header()["Content-Type"] = []string{"text/html; charset=utf-8"}
		self.Writer.Write(buffer.Bytes())
	} else {
		e := fmt.Errorf("Template <%s> does not exists", filename)
		self.Error(http.StatusInternalServerError, e.Error())
		return
	}
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
	if flusher, ok := self.Writer.(http.Flusher); ok {
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

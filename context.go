/**
 *  ------------------------------------------------------------
 *  @project	web.go
 *  @file       context.go
 *  @date       2014-10-21
 *  @author     Jim Zhan <jim.zhan@me.com>
 *
 *  Copyright © 2014 Jim Zhan.
 *  ------------------------------------------------------------
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://wwself.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *  ------------------------------------------------------------
 */
package web

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/gorilla/securecookie"
)

const ContentType = "Content-Type"

var (
	prefix   string
	identity uint64

	secret string
	secure *securecookie.SecureCookie

	templates map[string]*template.Template
)

type (
	ResponseWriter interface {
		http.ResponseWriter
		http.Hijacker
		http.Flusher
		http.CloseNotifier

		Status() int
		Size() int
		Written() bool
	}

	Context struct {
		http.ResponseWriter
		status  int
		size    int
		Request *http.Request
		data    map[string]interface{}
	}
)

func init() {
	templates = make(map[string]*template.Template)

	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}

	// system pid combined with timestamp to identity current go process.
	pid := fmt.Sprintf("%d:%d", os.Getpid(), time.Now().UnixNano())
	prefix = fmt.Sprintf("%s-%s", hostname, base64.URLEncoding.EncodeToString([]byte(pid)))
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	var ctx *Context = new(Context)
	ctx.ResponseWriter = w
	ctx.Request = r
	ctx.status = http.StatusOK
	ctx.size = -1
	return ctx
}

// ---------------------------------------------------------------------------
//  Enhancements for native http.ResponseWriter
// ---------------------------------------------------------------------------
func (self *Context) Status() int {
	return self.status
}

func (self *Context) Size() int {
	return self.size
}

func (self *Context) Written() bool {
	return self.size != -1
}

// ---------------------------------------------------------------------------
//  Implementation of http.ResponseWriter#WriteHeader
// ---------------------------------------------------------------------------
func (self *Context) WriteHeader(status int) {
	if status > 0 && !self.Written() {
		self.status = status
		self.ResponseWriter.WriteHeader(status)
	}
}

// ---------------------------------------------------------------------------
//  Implementation of http.ResponseWriter#Write
// ---------------------------------------------------------------------------
func (self *Context) Write(data []byte) (n int, err error) {
	if !self.Written() {
		self.WriteHeader(http.StatusOK)
	}
	size, err := self.ResponseWriter.Write(data)
	self.size += size
	return size, err
}

// ---------------------------------------------------------------------------
//  Implementations of http.Hijackeri#Hijack
// ---------------------------------------------------------------------------
func (self *Context) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := self.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("the ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}

// ---------------------------------------------------------------------------
//  Implementations of http.CloseNotifier#CloseNotify
// ---------------------------------------------------------------------------
func (self *Context) CloseNotify() <-chan bool {
	return self.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

// ---------------------------------------------------------------------------
//  Implementations of http.Flusher#Flush
// ---------------------------------------------------------------------------
func (self *Context) Flush() {
	flusher, ok := self.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}

// ---------------------------------------------------------------------------
//  HTTP Request Context Data
// ---------------------------------------------------------------------------
func (self *Context) Id() string {
	requestId := atomic.AddUint64(&identity, 1)
	return fmt.Sprintf("%s-%06d", prefix, requestId)
}

func (self *Context) Get(key string) interface{} {
	value, exists := self.data[key]
	if !exists {
		return nil
	}
	return value
}

func (self *Context) Set(key string, value interface{}) {
	if self.data == nil {
		self.data = make(map[string]interface{})
	}
	self.data[key] = value
}

func (self *Context) Clear() {
	for key := range self.data {
		delete(self.data, key)
	}
}

func (self *Context) Delete(key string) {
	delete(self.data, key)
}

// ---------------------------------------------------------------------------
//  HTTP Cookies
// ---------------------------------------------------------------------------
func (self *Context) Cookie(name string) string {
	var value string
	if cookie, err := self.Request.Cookie(name); err == nil {
		value = cookie.Value
	}
	return value
}

func (self *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(self, cookie)
}

func (self *Context) SecureCookie(name string) string {
	var value string
	if cookie, err := self.Request.Cookie(name); err == nil {
		if err = secure.Decode(name, cookie.Value, &value); err == nil {
			return value
		}
	}
	return value
}

// SetSecureCookie signs a cookie so it cannot be forged.
func (self *Context) SetSecureCookie(cookie *http.Cookie) {
	if secret == "" {
		panic("Application secret is missing from settings file.")
	}

	// initialize SecureCookie when first set.
	if secure == nil {
		secure = securecookie.New([]byte(secret), nil)
	}

	if value, err := secure.Encode(cookie.Name, cookie.Value); err == nil {
		cookie.Value = value
	}
	self.SetCookie(cookie)
}

// ---------------------------------------------------------------------------
//  HTTP Request Helpers
// ---------------------------------------------------------------------------
func (self *Context) ClientIP() string {
	clientIP := self.Request.Header.Get("X-Real-IP")
	if clientIP == "" {
		clientIP = self.Request.Header.Get("X-Forwarded-For")
	}
	if clientIP == "" {
		clientIP, _, _ = net.SplitHostPort(self.Request.RemoteAddr)
	}
	return clientIP
}

func (self *Context) IsAjax() bool {
	return self.Request.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

func (self *Context) Query() url.Values {
	return self.Request.URL.Query()
}

// ---------------------------------------------------------------------------
//  HTTP Response Rendering
// ---------------------------------------------------------------------------
// Forcely parse the passed in template files under the pre-defined template folder,
// & panics if the error is non-nil. It also try finding the default layout page (defined
// in ctx.Options.Layout) as the render base first, the parsed template page will be
// cached in global singleton holder.
func (self *Context) parseFiles(filename string, others ...string) *template.Template {
	page, exists := templates[filename]
	if !exists {
		var files []string
		folder := Settings.GetStringMapString("folder")["templates"]

		files = append(files, filepath.Join(folder, filename))
		for _, item := range others {
			files = append(files, filepath.Join(folder, item))
		}

		page = template.Must(template.ParseFiles(files...))
		templates[filename] = page
	}
	return page
}

// Shortcut to render HTML templates, with basic layout supports.
func (self *Context) HTML(filename string, others ...string) {
	buffer := new(bytes.Buffer)
	self.Header().Set(ContentType, "text/html; charset=utf-8")
	self.WriteHeader(self.status)
	err := self.parseFiles(filename, others...).Execute(buffer, self.data)
	if err != nil {
		panic(err)
	}
	buffer.WriteTo(self)
}

func (self *Context) JSON(values map[string]interface{}) {
	var (
		data []byte
		err  error
	)

	data, err = json.Marshal(values)

	if err != nil {
		http.Error(self, err.Error(), http.StatusInternalServerError)
		return
	}

	self.Header().Set(ContentType, "application/json; charset=utf-8")
	self.WriteHeader(self.status)
	self.Write(data)
}

func (self *Context) XML(values interface{}) {
	self.Header().Set(ContentType, "application/xml; charset=utf-8")
	self.WriteHeader(self.status)
	encoder := xml.NewEncoder(self)
	encoder.Encode(values)
}

// String writes plain text back to the HTTP response.
func (self *Context) String(content string) {
	self.Header().Set(ContentType, "text/plain; charset=utf-8")
	self.WriteHeader(self.status)
	self.Write([]byte(content))
}

// Data writes binary data back into the HTTP response.
func (self *Context) Data(data []byte) {
	self.Header().Set(ContentType, "application/octet-stream")
	self.WriteHeader(self.status)
	self.Write(data)
}

// ---------------------------------------------------------------------------
//  HTTP Status Shortcuts
// ---------------------------------------------------------------------------
// 4XX Client errors
// ---------------------------------------------------------------------------
func (self *Context) Forbidden(message string) {
	http.Error(self, message, http.StatusForbidden)
}

func (self *Context) NotFound(message string) {
	http.Error(self, message, http.StatusNotFound)
}

// ---------------------------------------------------------------------------
// 5XX Server errors
// ---------------------------------------------------------------------------

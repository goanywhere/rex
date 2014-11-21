/* ----------------------------------------------------------------------
 *       ______      ___                         __
 *      / ____/___  /   |  ____  __  ___      __/ /_  ___  ________
 *     / / __/ __ \/ /| | / __ \/ / / / | /| / / __ \/ _ \/ ___/ _ \
 *    / /_/ / /_/ / ___ |/ / / / /_/ /| |/ |/ / / / /  __/ /  /  __/
 *    \____/\____/_/  |_/_/ /_/\__. / |__/|__/_/ /_/\___/_/   \___/
 *                            /____/
 *
 * (C) Copyright 2014 GoAnywhere (http://goanywhere.io).
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
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync/atomic"
	"time"

	"github.com/goanywhere/web/crypto"
	"github.com/goanywhere/web/env"
)

const ContentType = "Content-Type"

var (
	cid       uint64
	cidPrefix string
	signature *crypto.Signature
)

type Context struct {
	http.ResponseWriter
	Request *http.Request

	status int
	size   int
	data   map[string]interface{}
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := new(Context)
	ctx.ResponseWriter = w
	ctx.Request = r
	ctx.size = -1
	ctx.createSignature()
	return ctx
}

// ---------------------------------------------------------------------------
//  Internal Helpers
// ---------------------------------------------------------------------------
func (self *Context) createSignature() {
	if signature == nil {
		secret := env.Get("secret")
		if secret == "" {
			Warn("Secret key missing, using a random string now, previous cookie will be invalidate")
			pool := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*(-_+)")
			secret = crypto.RandomString(64, pool)
		}
		signature = crypto.NewSignature(secret)
	}
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
	return self.status != 0
}

// ---------------------------------------------------------------------------
//  Implementation of http.ResponseWriter#WriteHeader
// ---------------------------------------------------------------------------
func (self *Context) WriteHeader(status int) {
	if status >= 100 && status < 512 {
		self.status = status
		self.ResponseWriter.WriteHeader(status)
	}
}

// ---------------------------------------------------------------------------
//  Implementation of http.ResponseWriter#Write
// ---------------------------------------------------------------------------
func (self *Context) Write(data []byte) (size int, err error) {
	if !self.Written() {
		self.WriteHeader(http.StatusOK)
	}
	size, err = self.ResponseWriter.Write(data)
	self.size += size
	return
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
	requestId := atomic.AddUint64(&cid, 1)
	return fmt.Sprintf("%s-%06d", cidPrefix, requestId)
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
	http.SetCookie(self, cookie)
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
	http.SetCookie(self, cookie)
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
// Shortcut to render HTML templates, with basic layout supports.
func (self *Context) HTML(filename string, others ...string) {
	buffer := new(bytes.Buffer)
	self.Header().Set(ContentType, "text/html; charset=utf-8")
	if err := loadTemplates(filename, others...).Execute(buffer, self.data); err != nil {
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
	self.Write(data)
}

func (self *Context) XML(values interface{}) {
	self.Header().Set(ContentType, "application/xml; charset=utf-8")
	encoder := xml.NewEncoder(self)
	encoder.Encode(values)
}

// String writes plain text back to the HTTP response.
func (self *Context) String(format string, values ...interface{}) {
	self.Header().Set(ContentType, "text/plain; charset=utf-8")
	self.Write([]byte(fmt.Sprintf(format, values...)))
}

// Data writes binary data back into the HTTP response.
func (self *Context) Data(data []byte) {
	self.Header().Set(ContentType, "application/octet-stream")
	self.Write(data)
}

// ---------------------------------------------------------------------------
//  Context Prerequisites
// ---------------------------------------------------------------------------
func init() {
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
	// system pid combined with timestamp to identity current go process.
	pid := fmt.Sprintf("%d:%d", os.Getpid(), time.Now().UnixNano())
	cidPrefix = fmt.Sprintf("%s-%s", hostname, base64.URLEncoding.EncodeToString([]byte(pid)))
}

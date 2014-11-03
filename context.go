/**
 *  ------------------------------------------------------------
 *  @project	web
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
 *      http://www.apache.org/licenses/LICENSE-2.0
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
	"bytes"
	"encoding/json"
	"encoding/xml"
	"html/template"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"sync"

	"github.com/gorilla/securecookie"
)

const ContentType string = "Content-Type"

var (
	mutex     sync.RWMutex
	options   *opt
	contexts  map[*http.Request](map[string]interface{})
	templates map[string]*template.Template

	secret string
	secure *securecookie.SecureCookie
)

type (
	opt struct {
		Indent bool
		Layout string
	}

	Context struct {
		// I/O access
		http.ResponseWriter
		Request *http.Request

		Options *opt
		Status  int
	}
)

func init() {
	options = &opt{Indent: true, Layout: "layout.html"}
	contexts = make(map[*http.Request]map[string]interface{})
	templates = make(map[string]*template.Template)
}

func NewContext(writer http.ResponseWriter, request *http.Request) *Context {
	if ctx := contexts[request]; ctx == nil {
		contexts[request] = make(map[string]interface{})
	}
	return &Context{writer, request, options, http.StatusOK}
}

// ---------------------------------------------------------------------------
//  HTTP Request Context Data
// ---------------------------------------------------------------------------
func (self *Context) Get(key string) interface{} {
	mutex.RLock()
	defer mutex.RUnlock()

	return contexts[self.Request][key]
}

func (self *Context) Set(key string, value interface{}) {
	mutex.Lock()
	defer mutex.Unlock()

	contexts[self.Request][key] = value
}

func (self *Context) Clear() {
	mutex.Lock()
	defer mutex.Unlock()
	// TODO incorperate with map initialization
	delete(contexts, self.Request)
}

func (self *Context) Delete(key string) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(contexts[self.Request], key)
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
		folder := Settings.GetStringMapString("folder")["templates"]
		var files []string
		if self.Options.Layout != "" {
			files = append(files, filepath.Join(folder, self.Options.Layout))
		}
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
	self.WriteHeader(self.Status)
	err := self.parseFiles(filename, others...).Execute(buffer, contexts[self.Request])
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

	if self.Options.Indent {
		data, err = json.MarshalIndent(values, "", "\t")
	} else {
		data, err = json.Marshal(values)
	}

	if err != nil {
		http.Error(self, err.Error(), http.StatusInternalServerError)
		return
	}

	self.Header().Set(ContentType, "application/json; charset=utf-8")
	self.WriteHeader(self.Status)
	self.Write(data)
}

func (self *Context) XML(values interface{}) {
	self.Header().Set(ContentType, "application/xml; charset=utf-8")
	self.WriteHeader(self.Status)
	encoder := xml.NewEncoder(self)
	encoder.Encode(values)
}

// String writes plain text back to the HTTP response.
func (self *Context) String(content string) {
	self.Header().Set(ContentType, "text/plain; charset=utf-8")
	self.WriteHeader(self.Status)
	self.Write([]byte(content))
}

// Data writes binary data back into the HTTP response.
func (self *Context) Data(data []byte) {
	self.Header().Set(ContentType, "application/octet-stream")
	self.WriteHeader(self.Status)
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

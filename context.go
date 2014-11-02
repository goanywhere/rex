/**
 *  ------------------------------------------------------------
 *  @project	webapp
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
package webapp

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
)

const ContentType string = "Content-Type"

var (
	mutex     sync.RWMutex
	options   *opt
	contexts  map[*http.Request](map[string]interface{})
	templates map[string]*template.Template
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
func (ctx *Context) Get(key string) interface{} {
	mutex.RLock()
	defer mutex.RUnlock()

	return contexts[ctx.Request][key]
}

func (ctx *Context) Set(key string, value interface{}) {
	mutex.Lock()
	defer mutex.Unlock()

	contexts[ctx.Request][key] = value
}

func (ctx *Context) Clear() {
	mutex.Lock()
	defer mutex.Unlock()
	// TODO incorperate with map initialization
	delete(contexts, ctx.Request)
}

func (ctx *Context) Delete(key string) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(contexts[ctx.Request], key)
}

// ---------------------------------------------------------------------------
//  HTTP Cookies
// ---------------------------------------------------------------------------
func (ctx *Context) Cookie(name string) string {
	cookie, err := ctx.Request.Cookie(name)
	if err == nil {
		return cookie.Value
	}
	return ""
}

func (ctx *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(ctx, cookie)
}

func (ctx *Context) SecureCookie(name string) string {
	return ""
}

func (ctx *Context) SetSecureCookie(cookie *http.Cookie) {
	panic("Not Implemented Yet")
}

// ---------------------------------------------------------------------------
//  HTTP Request Helpers
// ---------------------------------------------------------------------------
func (ctx *Context) ClientIP() string {
	clientIP := ctx.Request.Header.Get("X-Real-IP")
	if clientIP == "" {
		clientIP = ctx.Request.Header.Get("X-Forwarded-For")
	}
	if clientIP == "" {
		clientIP, _, _ = net.SplitHostPort(ctx.Request.RemoteAddr)
	}
	return clientIP
}

func (ctx *Context) IsAjax() bool {
	return ctx.Request.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

func (ctx *Context) Query() url.Values {
	return ctx.Request.URL.Query()
}

// ---------------------------------------------------------------------------
//  HTTP Response Rendering
// ---------------------------------------------------------------------------

// Forcely parse the passed in template files under the pre-defined template folder,
// & panics if the error is non-nil. It also try finding the default layout page (defined
// in ctx.Options.Layout) as the render base first, the parsed template page will be
// cached in global singleton holder.
func (ctx *Context) parseFiles(filename string, others ...string) *template.Template {
	page, exists := templates[filename]
	if !exists {
		folder := Settings.GetStringMapString("folder")["templates"]
		var files []string
		if ctx.Options.Layout != "" {
			files = append(files, filepath.Join(folder, ctx.Options.Layout))
			Logger.Debug("Using default layout: " + filepath.Join(folder, ctx.Options.Layout))
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
func (ctx *Context) HTML(filename string, others ...string) {
	buffer := new(bytes.Buffer)
	ctx.Header().Set(ContentType, "text/html; charset=utf-8")
	ctx.WriteHeader(ctx.Status)
	err := ctx.parseFiles(filename, others...).Execute(buffer, contexts[ctx.Request])
	if err != nil {
		panic(err)
	}
	buffer.WriteTo(ctx)
}

func (ctx *Context) JSON(values map[string]interface{}) {
	var (
		data []byte
		err  error
	)

	if ctx.Options.Indent {
		data, err = json.MarshalIndent(values, "", "\t")
	} else {
		data, err = json.Marshal(values)
	}

	if err != nil {
		http.Error(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.Header().Set(ContentType, "application/json; charset=utf-8")
	ctx.WriteHeader(ctx.Status)
	ctx.Write(data)
}

func (ctx *Context) XML(values interface{}) {
	ctx.Header().Set(ContentType, "application/xml; charset=utf-8")
	ctx.WriteHeader(ctx.Status)
	encoder := xml.NewEncoder(ctx)
	encoder.Encode(values)
}

// String writes plain text back to the HTTP response.
func (ctx *Context) String(content string) {
	ctx.Header().Set(ContentType, "text/plain; charset=utf-8")
	ctx.WriteHeader(ctx.Status)
	ctx.Write([]byte(content))
}

// Data writes binary data back into the HTTP response.
func (ctx *Context) Data(data []byte) {
	ctx.Header().Set(ContentType, "application/octet-stream")
	ctx.WriteHeader(ctx.Status)
	ctx.Write(data)
}

// ---------------------------------------------------------------------------
//  HTTP Status Shortcuts
// ---------------------------------------------------------------------------
// 4XX Client errors
// ---------------------------------------------------------------------------
func (ctx *Context) Forbidden(message string) {
	http.Error(ctx, message, http.StatusForbidden)
}

func (ctx *Context) NotFound(message string) {
	http.Error(ctx, message, http.StatusNotFound)
}

// ---------------------------------------------------------------------------
// 5XX Server errors
// ---------------------------------------------------------------------------

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
	"net/http"
	"path/filepath"
	"sync"
)

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
	context struct {
		http.ResponseWriter
		request *http.Request

		Options *opt
	}
)

func init() {
	options = &opt{Indent: true, Layout: "layout.html"}
	contexts = make(map[*http.Request]map[string]interface{})
	templates = make(map[string]*template.Template)
}

func Context(writer http.ResponseWriter, request *http.Request) *context {
	if ctx := contexts[request]; ctx == nil {
		contexts[request] = make(map[string]interface{})
	}
	return &context{writer, request, options}
}

func (ctx *context) Get(key string) interface{} {
	mutex.RLock()
	defer mutex.RUnlock()

	return contexts[ctx.request][key]
}

func (ctx *context) Set(key string, value interface{}) {
	mutex.Lock()
	defer mutex.Unlock()

	contexts[ctx.request][key] = value
}

func (ctx *context) Delete(key string) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(contexts[ctx.request], key)
}

// ---------------------------------------------------------------------------
//  HTTP Request Helpers
// ---------------------------------------------------------------------------
func (ctx *context) ClientIP() string {
	clientIP := ctx.request.Header.Get("X-Real-IP")
	if clientIP == "" {
		clientIP = ctx.request.Header.Get("X-Forwarded-For")
	}
	if clientIP == "" {
		clientIP = ctx.request.RemoteAddr
	}
	return clientIP
}

// ---------------------------------------------------------------------------
//  HTTP Response Rendering
// ---------------------------------------------------------------------------

// Forcely parse the passed in template files under the pre-defined template folder,
// & panics if the error is non-nil. It also try finding the default layout page (defined
// in ctx.Options.Layout) as the render base first, the parsed template page will be
// cached in global singleton holder.
func (ctx *context) parseFiles(filename string, others ...string) *template.Template {
	page, exists := templates[filename]
	if !exists {
		folder := Settings.GetStringMapString("folder")["templates"]
		var files []string
		if ctx.Options.Layout != "" {
			files = append(files, filepath.Join(folder, ctx.Options.Layout))
			Debug("Using default layout: " + filepath.Join(folder, ctx.Options.Layout))
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
func (ctx *context) HTML(status int, filename string, others ...string) {
	buffer := new(bytes.Buffer)
	ctx.Header().Set(ContentType, "text/html; charset=utf-8")
	ctx.WriteHeader(status)
	err := ctx.parseFiles(filename, others...).Execute(buffer, contexts[ctx.request])
	if err != nil {
		panic(err)
	}
	buffer.WriteTo(ctx)
}

func (ctx *context) JSON(status int, values map[string]interface{}) {

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
	ctx.WriteHeader(status)
	ctx.Write(data)
}

func (ctx *context) XML(status int, values interface{}) {
	ctx.Header().Set(ContentType, "application/xml; charset=utf-8")
	ctx.WriteHeader(status)
	encoder := xml.NewEncoder(ctx)
	encoder.Encode(values)
}

// String writes plain text back to the HTTP response.
// TODO response in buffer.
func (ctx *context) String(status int, content string) {
	ctx.Header().Set(ContentType, "text/plain; charset=utf-8")
	ctx.WriteHeader(status)
	ctx.Write([]byte(content))
}

func (ctx *context) Data(status int, data []byte) {
	ctx.Header().Set(ContentType, "application/octet-stream")
	ctx.WriteHeader(status)
	ctx.Write(data)
}

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
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

var (
	templates map[string]*template.Template
)

type context struct {
	Layout   string
	Response http.ResponseWriter
	Request  *http.Request
	data     map[string]interface{}
	mutex    sync.RWMutex
}

func getTemplatePath(filename string) string {
	folder := Settings.GetStringMapString("folder")["templates"]
	return filepath.Join(folder, filename)
}

func Context(writer http.ResponseWriter, request *http.Request) *context {
	c := &context{
		Layout:   "layout.html",
		Response: writer,
		Request:  request,
		data:     make(map[string]interface{}),
	}
	return c
}

func (c *context) Get(key string) interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.data[key]
}

func (c *context) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = value
}

func (c *context) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.data[key] != nil {
		delete(c.data, key)
	}
}

// String writes plain text back to the HTTP response.
// TODO response in buffer.
func (c *context) String(status uint8, content string) {
	c.Response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Response.Write([]byte(content))
}

// Render creates combined file templates using html/template package,
// Content-Type will be determined by context.Layout's extension.
// NOTE Pre-defined base template name *MUST* be the same as the file name itself.
// TODO buffer based caching.
func (c *context) Render(status uint8, file string, others ...string) {
	// self-determined for Content-Type
	switch filepath.Ext(c.Layout) {
	case ".json":
		c.Response.Header().Set(ContentType, JSON)
	case ".xml":
		c.Response.Header().Set(ContentType, XML)
	default:
		c.Response.Header().Set(ContentType, HTML)
	}
	// Multiple file templates parsing.
	files := []string{getTemplatePath(c.Layout), getTemplatePath(file)}
	for _, item := range others {
		files = append(files, getTemplatePath(item))
	}
	layout := template.Must(template.ParseFiles(files...))
	layout.ExecuteTemplate(c.Response, c.Layout, c.data)
}

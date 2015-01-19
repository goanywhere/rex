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
/*
Package rex provides an out-of-box web server with common middleware modules.
Example:
	package main

	import (
		"net/http"
		"github.com/goanywhere/rex"
		"github.com/goanywhere/rex/web"
	)

	func index(ctx *web.Context) {
		ctx.String("Hello World")
	}

	func main() {
		rex.Get("/", index)
		rex.Run()
	}
*/
package rex

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/goanywhere/rex/internal"
	"github.com/goanywhere/rex/modules"
	"github.com/goanywhere/rex/web"
)

// default rex server with reasonable middleware modules.
var (
	port int
	mux  = web.New()

	options = internal.Options()
)

type H map[string]interface{}

// Define saves primitive values using os environment.
func Define(key string, value interface{}) error {
	return internal.Define(key, value)
}

// Option fetches the value from os's enviroment into the appointed address.
func Option(key string, ptr interface{}, fallback ...interface{}) {
	internal.Option(key, ptr, fallback...)
}

// Get adds a HTTP GET route to the default Mux.
func Get(pattern string, handler interface{}) {
	mux.Get(pattern, handler)
}

// Post adds a HTTP POST route to the default Mux.
func Post(pattern string, handler interface{}) {
	mux.Post(pattern, handler)
}

// Put adds a HTTP PUT route to the default Mux.
func Put(pattern string, handler interface{}) {
	mux.Put(pattern, handler)
}

// Delete adds a HTTP DELETE route to the default Mux.
func Delete(pattern string, handler interface{}) {
	mux.Delete(pattern, handler)
}

// Patch adds a HTTP PATCH route to the default Mux.
func Patch(pattern string, handler http.HandlerFunc) {
	mux.Patch(pattern, handler)
}

// Head adds a HTTP HEAD route to the default Mux.
func Head(pattern string, handler http.HandlerFunc) {
	mux.Head(pattern, handler)
}

// Options adds a HTTP OPTIONS route to the default Mux.
func Options(pattern string, handler http.HandlerFunc) {
	mux.Options(pattern, handler)
}

// Group creates a new application group in default Mux with the given path.
func Group(path string) *web.Mux {
	return mux.Group(path)
}

// Use appends middleware module into the default serving list.
func Use(modules ...interface{}) {
	mux.Use(modules...)
}

// Serve starts serving the requests at the pre-defined address from settings.
func Run() {
	flag.Parse()
	if port > 0 {
		options.Set("port", port)
	}
	mux.Run(fmt.Sprintf(":%d", options.Int("port")))
}

func init() {
	mux.Use(modules.Env)
	mux.Use(modules.XSRF)

	if cwd, err := os.Getwd(); err != nil {
		log.Fatalf("Failed to retrieve project root: %v", err)
	} else {
		root, _ := filepath.Abs(cwd)
		Define("root", root)
	}
	options.Load(".env")

	// cmd parameters take the priority.
	flag.IntVar(&port, "port", 0, "port to run the application server")
}

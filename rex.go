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
 *  Unless required by muxlicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 * ----------------------------------------------------------------------*/
/*
Package rex provides an out-of-box web mux with common middleware modules.
Example:
	package main

	import (
		"net/http"
		"github.com/goanywhere/rex"
	)

	func index(ctx *rex.Context) {
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
	"net/http"

	"github.com/goanywhere/rex/internal"
	"github.com/goanywhere/rex/modules"
)

// default rex mux with reasonable middleware modules.
var (
	port       int
	DefaultMux = New()
	Options    = internal.Options()
)

type H map[string]interface{}

// Define saves primitive values using os environment.
func Define(key string, value interface{}) error {
	return Options.Set(key, value)
}

// Get adds a HTTP GET route to the default DefaultMux.
func Get(pattern string, handler interface{}) {
	DefaultMux.Get(pattern, handler)
}

// Post adds a HTTP POST route to the default DefaultMux.
func Post(pattern string, handler interface{}) {
	DefaultMux.Post(pattern, handler)
}

// Put adds a HTTP PUT route to the default DefaultMux.
func Put(pattern string, handler interface{}) {
	DefaultMux.Put(pattern, handler)
}

// Delete adds a HTTP DELETE route to the default DefaultMux.
func Delete(pattern string, handler interface{}) {
	DefaultMux.Delete(pattern, handler)
}

// Head adds a HTTP HEAD route to the default DefaultMux.
func Head(pattern string, handler http.HandlerFunc) {
	DefaultMux.Head(pattern, handler)
}

// Group creates a new muxlication group in default Mux with the given path.
func Group(path string) *Mux {
	return DefaultMux.Group(path)
}

// Use muxends middleware module into the default serving list.
func Use(modules ...interface{}) {
	DefaultMux.Use(modules...)
}

// Serve starts serving the requests at the pre-defined address from settings.
func Run() {
	flag.Parse()
	if port > 0 {
		Options.Set("port", port)
	}
	DefaultMux.Run(fmt.Sprintf(":%d", Options.Int("port")))
}

func init() {
	DefaultMux.Use(modules.Env)
	DefaultMux.Use(modules.XSRF)

	// cmd parameters take the priority.
	flag.IntVar(&port, "port", 0, "port to run the muxlication mux")
}

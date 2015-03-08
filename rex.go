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
	)

	func index(ctx *rex.Context) {
		ctx.Render("Hello World")
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
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"

	"github.com/goanywhere/rex/internal"
	"github.com/goanywhere/rex/modules"
	"github.com/goanywhere/rex/web"
)

var (
	settings = internal.Settings()

	// Define saves primitive values using os environment.
	Define = settings.Set

	// Bool retrieves boolean value associated with the given key from environ.
	Bool = settings.Bool

	// Float retrieves float64 value associated with the given key from environ.
	Float = settings.Float

	// Int retrieves int value associated with the given key from environ.
	Int = settings.Int

	// Int64 retrieves int64 value associated with the given key from environ.
	Int64 = settings.Int64

	// String retrieves string value associated with the given key from environ.
	String = settings.String

	// Strings retrieves string array associated with the given key from environ.
	Strings = settings.Strings
)

// default rex mux with reasonable middleware modules.
var server = web.New()

type H map[string]interface{}

// Get adds a HTTP GET route to the default server.
func Get(pattern string, handler interface{}) {
	server.Get(pattern, handler)
}

// Post adds a HTTP POST route to the default server.
func Post(pattern string, handler interface{}) {
	server.Post(pattern, handler)
}

// Put adds a HTTP PUT route to the default server.
func Put(pattern string, handler interface{}) {
	server.Put(pattern, handler)
}

// Delete adds a HTTP DELETE route to the default server.
func Delete(pattern string, handler interface{}) {
	server.Delete(pattern, handler)
}

// Head adds a HTTP HEAD route to the default server.
func Head(pattern string, handler http.HandlerFunc) {
	server.Head(pattern, handler)
}

// Group creates a new muxlication group in default Mux with the given path.
func Group(path string) *web.Server {
	return server.Group(path)
}

func FileServer(prefix, dir string) {
	server.FileServer(prefix, dir)
}

// Use muxends middleware module into the default serving list.
func Use(modules ...web.Module) {
	server.Use(modules...)
}

// Serve starts serving the requests at the pre-defined address from settings.
var port int

func Run() {
	// common server middleware modules.
	//server.Use(modules.XSRF)
	server.Use(modules.Env)
	server.Use(modules.LiveReload)

	flag.Parse()
	server.Run(fmt.Sprintf(":%d", port))
}

func init() {

	// setup fundamental project root.
	if cwd, err := os.Getwd(); err == nil {
		root, _ := filepath.Abs(cwd)
		Define("root", root)
	} else {
		log.Fatalf("Failed to retrieve project root: %v", err)
	}

	// cmd parameters take the priority.
	flag.IntVar(&port, "port", 5000, "port to run the application server")
}

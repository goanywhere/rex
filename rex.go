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
		ctx.Send("Hello World")
	}

	func main() {
		rex.Get("/", index)
		rex.Run()
	}
*/
package rex

import (
	"flag"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/goanywhere/rex/internal"
	"github.com/goanywhere/rex/modules"
	"github.com/goanywhere/rex/web"
	"github.com/goanywhere/x/env"
)

var (
	// Serve starts serving the requests at the pre-defined address from settings.
	Port = settings.Port

	server *web.Server

	settings = internal.Settings()
)

// Shortcut to create hash map.
type M map[string]interface{}

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

func Run() {
	// common server middleware modules.
	//server.Use(modules.XSRF)
	server.Use(modules.Env)
	server.Use(modules.LiveReload)

	flag.Parse()
	server.Run(Port)
}

func init() {
	/*** project root ***/
	pc, _, _, _ := runtime.Caller(2)
	function := runtime.FuncForPC(pc)
	filename, _ := function.FileLine(0)
	root := path.Dir(filename)
	os.Setenv("rex.root", root)

	/*** custom settings ***/
	env.Load(filepath.Join(root, ".env"))
	env.Map(settings)
	env.Map(settings.Session)

	server = web.New()

	// cmd parameters take the priority.
	flag.BoolVar(&settings.Debug, "debug", settings.Debug, "flag to toggle debug mode")
	flag.IntVar(&settings.Port, "port", settings.Port, "port to run the application server")
}

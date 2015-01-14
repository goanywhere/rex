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
	"fmt"
	"log"
	"net/http"
	"time"

	. "github.com/goanywhere/rex/config"
	"github.com/goanywhere/rex/modules"
	"github.com/goanywhere/rex/web"
)

var (
	server   *web.Server
	Settings = config.Settings()
)

type H map[string]interface{}

func Get(pattern string, handler interface{}) {
	server.Get(pattern, handler)
}

func Post(pattern string, handler interface{}) {
	server.Post(pattern, handler)
}

func Put(pattern string, handler interface{}) {
	server.Put(pattern, handler)
}

func Delete(pattern string, handler interface{}) {
	server.Delete(pattern, handler)
}

func Patch(pattern string, handler http.HandlerFunc) {
	server.Patch(pattern, handler)
}

func Head(pattern string, handler http.HandlerFunc) {
	server.Head(pattern, handler)
}

func Options(pattern string, handler http.HandlerFunc) {
	server.Options(pattern, handler)
}

func Group(path string) *web.Server {
	return server.Group(path)
}

func Use(modules ...interface{}) {
	server.Use(modules...)
}

// Serve starts serving the requests at the pre-defined address from settings.
func Run() {
	go func() {
		time.Sleep(100 * time.Millisecond)
		log.Printf("Application server started [:%d]", Port)
	}()
	if err := http.ListenAndServe(fmt.Sprintf(":%d", Port), server); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}

func init() {
	server = web.New()
	server.Use(modules.Env)
	server.Use(modules.XSRF)
	server.Use(modules.Static)
	server.Use(modules.LiveReload)
	server.Use(modules.Compress)
}

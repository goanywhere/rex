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
package rex

import (
	"flag"
	"os"
	"path/filepath"
	"regexp"

	pongo "github.com/flosch/pongo2"
	cookie "github.com/gorilla/securecookie"

	"github.com/goanywhere/rex/internal"
	"github.com/goanywhere/rex/modules"
	"github.com/goanywhere/x/env"
	"github.com/goanywhere/x/fs"
)

var (
	root string
)

var flags = struct {
	debug bool
	port  int
}{
	debug: true,
	port:  5000,
}

var (
	securecookie *cookie.SecureCookie

	views map[string]*pongo.Template = make(map[string]*pongo.Template)
)

// Shortcut to create hash map.
type M map[string]interface{}

// loadViews load the html/xml documents from the pre-defined directory,
// rex will ignores directories named "layouts" & "include".
// TODO multiple paths supports.
func loadViews() {
	var (
		files   = regexp.MustCompile(`\.(html|xml)$`)
		ignores = regexp.MustCompile(`(layouts|include|\.(\w+))`)
	)
	dir := filepath.Join(root, env.String("VIEWS", "views"))

	if fs.Exists(dir) {
		filepath.Walk(dir, func(path string, info os.FileInfo, e error) error {

			if info.IsDir() {
				if ignores.MatchString(info.Name()) {
					return filepath.SkipDir
				} else {
					return nil
				}
			}

			if files.MatchString(path) {
				key, _ := filepath.Rel(dir, path)
				views[key] = pongo.Must(pongo.FromFile(path))
			}

			return e
		})
	}
}

// ---------------------------------------------------------------------------
//  Default Server Mux
// ---------------------------------------------------------------------------
var web *Server

func Get(pattern string, handler interface{}) {
	web.Get(pattern, handler)
}

func Run() {
	// common server middleware modules.
	//server.Use(modules.XSRF)
	web.Use(modules.Env)

	flag.Parse()
	web.Run()
}

func init() {
	// ----------------------------------------
	// Project Root
	// ----------------------------------------
	if root = env.String(internal.Root); root == "" {
		root = fs.Getcd(2)
		env.Set(internal.Root, root)
	}
	env.Load(filepath.Join(root, ".env"))

	loadViews()

	if securecookie == nil {
		if secrets := env.Strings("SECRET_KEYS"); len(secrets) == 2 {
			securecookie = cookie.New([]byte(secrets[0]), []byte(secrets[1]))
		}
	}

	web = New()
	// cmd arguments
	flag.BoolVar(&flags.debug, "debug", flags.debug, "flag to toggle debug mode")
	flag.IntVar(&flags.port, "port", flags.port, "port to run the application server")
}

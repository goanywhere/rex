/* ----------------------------------------------------------------------
 *       ______      ___                         __
 *      / ____/___  /   |  ____  __  ___      __/ /_  ___  ________
 *     / / __/ __ \/ /| | / __ \/ / / / | /| / / __ \/ _ \/ ___/ _ \
 *    / /_/ / /_/ / ___ |/ / / / /_/ /| |/ |/ / / / /  __/ /  /  __/
 *    \____/\____/_/  |_/_/ /_/\__. / |__/|__/_/ /_/\___/_/   \___/
 *                            /____/
 *
 * (C) Copyright 2014 GoAnywhere (http://goanywhere.io).
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

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/goanywhere/fs"
	"github.com/goanywhere/rex/crypto"
)

var (
	pattern *regexp.Regexp
	pool    = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*(-_+)")
)

// getWorkspace finds the very first workspace as project base under $GOPATH.
func getWorkspace() (workspace string, err error) {
	if gopath := os.Getenv("GOPATH"); gopath != "" {
		workspace = strings.Split(gopath, ";")[0]
	} else {
		err = os.ErrNotExist
	}
	return
}

// createProject creates the given path under $GOPATH/src along with
// a dotenv file contains the generated secret key for your web app.
func createProject(path string) (project string, err error) {
	if gopath := os.Getenv("GOPATH"); gopath != "" {
		workspace := strings.Split(gopath, ";")[0]
		project = filepath.Join(workspace, "src", path)
		if fs.Exists(project) {
			project, err = "", os.ErrExist
			return
		}
		if err = os.MkdirAll(project, os.ModePerm); err == nil {
			log.Printf("Project created: %s", project)
			filename := filepath.Join(project, ".env")
			if dotenv, err := os.Create(filename); err == nil {
				defer dotenv.Close()
				buffer := bufio.NewWriter(dotenv)
				buffer.WriteString(fmt.Sprintf("secret=\"%s\"\n", crypto.RandomString(64, pool)))
				buffer.Flush()
				log.Printf("dotenv created: %s", filename)
			}
		}
	} else {
		err = os.ErrNotExist
	}
	return
}

// setupProject copies fixes assets into newly create project,
// generated project specific values for Go files.
func setupProject(project string) (err error) {
	_, me, _, _ := runtime.Caller(1)
	scaffold := filepath.Join(filepath.Dir(me), "scaffold")
	if err = fs.Copy(filepath.Join(scaffold, "assets"), project); err == nil {
		if err = fs.Copy(filepath.Join(scaffold, "templates"), project); err == nil {
			// Project Specific Values Go Here
		}
	}
	return
}

// 1. Fetch Golang Environment
// 2. Create Workspace under Given Namespace
// 3. Generate .env under created workspace.
// 4. Copy Fixes Assets
// 5. Render Text Template Go Files.
func Create(context *cli.Context) {
	args := context.Args()
	if len(args) != 1 || !pattern.MatchString(args[0]) {
		log.Printf("Please provide a valid project name/path")
	} else {
		var err error
		if project, err := createProject(args[0]); err == nil {
			if err = setupProject(project); err == nil {
				log.Printf("Project all set: %s", project)
				os.Exit(0)
			}
		}
		log.Fatalf("Failed to create project: %s", err)
	}
}

func init() {
	if runtime.GOOS == "windows" {
		pattern = regexp.MustCompile(`\A(?:[0-9a-zA-Z\.\_\-]+\\?)+\z`)
	} else {
		pattern = regexp.MustCompile(`\A(?:[0-9a-zA-Z\.\_\-]+\/?)+\z`)
	}
}

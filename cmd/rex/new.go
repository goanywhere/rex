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
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/goanywhere/rex/internal"
	"github.com/goanywhere/x/crypto"
)

const endpoint = "https://github.com/goanywhere/rex"

var secrets = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*(-_+)")

type project struct {
	name string
	root string
}

func (self *project) create() {
	internal.Prompt("Fetching project template\n")
	var done = make(chan bool)
	internal.Loading(done)

	cmd := exec.Command("git", "clone", "-b", "scaffolds", endpoint, self.name)
	cmd.Dir = cwd
	if e := cmd.Run(); e == nil {
		self.root = filepath.Join(cwd, self.name)
		// create dotenv under project's root.
		filename := filepath.Join(self.root, ".env")
		if dotenv, err := os.Create(filename); err == nil {
			defer dotenv.Close()
			buffer := bufio.NewWriter(dotenv)
			buffer.WriteString(fmt.Sprintf("export Rex_Secret_Keys=\"%s, %s\"\n", crypto.Random(64), crypto.Random(32)))
			buffer.Flush()
			// close loading here as nodejs will take over prompt.
			done <- true
			// initialize project packages via nodejs.
			self.setup()
		}
		os.RemoveAll(filepath.Join(self.root, ".git"))
		os.Remove(filepath.Join(self.root, "README.md"))
	} else {
		// loading prompt should be closed in anyway.
		done <- true
	}
}

func (self *project) setup() {
	if e := exec.Command("npm", "-v").Run(); e == nil {
		internal.Prompt("Fetching project dependencies\n")
		cmd := exec.Command("npm", "install")
		cmd.Dir = self.root
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	} else {
		log.Fatalf("Failed to setup project dependecies: nodejs is missing.")
	}
}

func New(context *cli.Context) {
	var pattern *regexp.Regexp
	if runtime.GOOS == "windows" {
		pattern = regexp.MustCompile(`\A(?:[0-9a-zA-Z\.\_\-]+\\?)+\z`)
	} else {
		pattern = regexp.MustCompile(`\A(?:[0-9a-zA-Z\.\_\-]+\/?)+\z`)
	}

	args := context.Args()
	if len(args) != 1 || !pattern.MatchString(args[0]) {
		log.Fatal("Please provide a valid project name/path")
	} else {
		project := new(project)
		project.name = args[0]
		project.create()
	}
}

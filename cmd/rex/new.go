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
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/goanywhere/rex/crypto"
)

const endpoint = "https://github.com/goanywhere/rex-scaffolds"

var secrets = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*(-_+)")

type project struct {
	name string
	root string
}

func (self *project) create() {
	cmd := exec.Command("git", "clone", endpoint, self.name)
	cmd.Dir = cwd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if e := cmd.Run(); e == nil {
		self.root = filepath.Join(cwd, self.name)
		// create dotenv under project's root.
		filename := filepath.Join(self.root, ".env")
		if dotenv, err := os.Create(filename); err == nil {
			defer dotenv.Close()
			buffer := bufio.NewWriter(dotenv)
			buffer.WriteString(fmt.Sprintf("SecretKey=\"%s\"\n", crypto.RandomString(64, secrets)))
			buffer.Flush()
			// initialize project packages via nodejs.
			self.setup()
		}
		os.RemoveAll(filepath.Join(self.root, ".git"))
		os.RemoveAll(filepath.Join(self.root, "README.md"))
		os.RemoveAll(filepath.Join(self.root, "LICENSE"))
	}
}

func (self *project) setup() {
	if e := exec.Command("npm", "-v").Run(); e == nil {
		fmt.Println("Setting up project dependencies...")
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
		log.Printf("Please provide a valid project name/path")
	} else {
		project := new(project)
		project.name = args[0]
		project.create()
	}
}

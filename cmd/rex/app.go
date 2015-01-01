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
	"go/build"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"syscall"

	"github.com/go-fsnotify/fsnotify"
	"github.com/goanywhere/fs"
	"github.com/goanywhere/rex/http/livereload"
)

var watchList = regexp.MustCompile(`\.(go|html|css|js|jsx|less|sass|scss)$`)

type (
	app struct {
		dir    string
		binary string
		args   []string

		Script string // script for npm.
		npm    *npm
	}

	npm struct {
		exist     bool
		installed bool
		script    string
	}
)

func NewApp() *app {
	cwd, _ := os.Getwd()

	pkg, err := build.ImportDir(cwd, build.AllowBinary)
	if err != nil || pkg.Name != "main" {
		log.Fatalf("No buildable Go source files found")
	}

	app := new(app)
	app.dir = cwd
	app.binary = "rex-bin"
	if runtime.GOOS == "windows" {
		app.binary += ".exe"
	}
	return app
}

// build compiles the application into rex-bin executable
// to run & optionally compiles static assets using npm.
func (self *app) build() {
	// * try build the application into rex-bin(.exe)
	cmd := exec.Command("go", "build", "-o", self.binary)
	cmd.Dir = self.dir
	if e := cmd.Run(); e != nil {
		log.Fatalf("Failed to compile the application: %v", e)
	}

	// initialize nodejs's status.
	if self.npm == nil {
		self.npm = new(npm)
		if e := exec.Command("npm", "-v").Run(); e == nil {
			self.npm.exist = true
			self.npm.installed = false
		}
	}
	// * initialize npm packages.
	if self.npm.exist && !self.npm.installed {
		cmd := exec.Command("npm", "install")
		cmd.Dir = self.dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
	// * run specify script using npm.
	if self.npm.exist && self.npm.installed {
		cmd := exec.Command("npm", "run-script", self.Script)
		cmd.Dir = self.dir
		cmd.Dir = self.dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}

// run executes the runnerable executable under package binary root.
func (self *app) run() (gorun chan bool) {
	gorun = make(chan bool)
	go func() {
		var proc *os.Process
		for start := range gorun {
			if proc != nil {
				// try soft kill before hard one.
				if err := proc.Signal(os.Interrupt); err != nil {
					proc.Kill()
				}
				proc.Wait()
			}
			if !start {
				continue
			}
			cmd := exec.Command(filepath.Join(self.dir, self.binary))
			cmd.Dir = self.dir
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Start(); err != nil {
				log.Fatalf("Failed to start the process: %v\n", err)
			}
			proc = cmd.Process
		}
	}()
	return
}

func (self *app) rerun(gorun chan bool) {
	self.build()
	livereload.Reload()
	gorun <- true
}

// Starts activates the application server along with
// a daemon watcher for monitoring the files's changes.
func (self *app) Start() {
	// ctrl-c: listen removes binary package when application stopped.
	channel := make(chan os.Signal, 2)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-channel
		// remove the binary package on stop.
		os.Remove(filepath.Join(self.dir, self.binary))
		os.Exit(1)
	}()

	// start waiting the signal to start running.
	var gorun = self.run()
	self.build()
	gorun <- true

	wd := fs.Watchdog(self.dir)
	wd.Add(watchList, func(event *fsnotify.Event) {
		self.rerun(gorun)
	})
	wd.Start()
}

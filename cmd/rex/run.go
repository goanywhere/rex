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
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/goanywhere/rex/internal"
	"github.com/goanywhere/rex/modules/livereload"
	"github.com/goanywhere/x/fs"
)

var (
	port      int
	watchList = regexp.MustCompile(`\.(go|html)$`)
)

type app struct {
	dir    string
	binary string
	args   []string

	task string // script for npm.
}

// build compiles the application into rex-bin executable
// to run & optionally compiles static assets using npm.
func (self *app) build() {
	var done = make(chan bool)
	internal.Loading(done)

	// * try build the application into rex-bin(.exe)
	cmd := exec.Command("go", "build", "-o", self.binary)
	cmd.Dir = self.dir
	if e := cmd.Run(); e != nil {
		log.Fatalf("Failed to compile the application: %v", e)
	}

	done <- true
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
			cmd := exec.Command(self.binary, fmt.Sprintf("--port=%d", port))
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
		os.Remove(self.binary)
		os.Exit(1)
	}()

	// start waiting the signal to start running.
	var gorun = self.run()
	self.build()
	gorun <- true

	watcher := fs.NewWatcher(self.dir)
	log.Infof("Start watching: %s", self.dir)
	watcher.Add(watchList, func(filepath string) {
		log.Infof("%s updated", filepath)
		self.rerun(gorun)
	})
	watcher.Start()
}

// Run creates an executable application package with livereload supports.
func Run(ctx *cli.Context) {
	port = ctx.Int("port")

	pkg, err := build.ImportDir(cwd, build.AllowBinary)
	if err != nil || pkg.Name != "main" {
		log.Fatalf("No buildable Go source files found")
	}
	app := new(app)
	app.dir = cwd
	app.binary = filepath.Join(os.TempDir(), "rex-bin")
	if runtime.GOOS == "windows" {
		app.binary += ".exe"
	}
	app.task = ctx.String("task")
	app.Start()
}

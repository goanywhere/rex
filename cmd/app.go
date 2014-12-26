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

package cmd

import (
	"bytes"
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"syscall"

	"github.com/go-fsnotify/fsnotify"
	"github.com/goanywhere/fs"
	"github.com/goanywhere/rex/http/livereload"
)

var watchList = regexp.MustCompile(`\.(css|js|html|go)$`)

type app struct {
	pkg *build.Package
	// binary
	name string
	path string
}

func New() *app {
	cwd, _ := os.Getwd()

	pkg, err := build.ImportDir(cwd, build.AllowBinary)
	if err != nil || pkg.Name != "main" {
		log.Fatalf("No runnable Go sources found.")
	}

	app := new(app)
	app.pkg = pkg
	app.name = filepath.Base(pkg.ImportPath)

	if gobin := os.Getenv("GOBIN"); gobin != "" {
		app.path = filepath.Join(gobin, app.name)
	} else {
		app.path = filepath.Join(app.pkg.BinDir, app.name)
	}
	return app
}

// install compiled the application binary package into binary root path.
func (self *app) install() (err error) {
	cmd := exec.Command("go", "get", self.pkg.ImportPath)

	buffer := bytes.NewBuffer([]byte{})
	cmd.Stdout = buffer
	cmd.Stderr = buffer

	err = cmd.Run()

	if buffer.Len() > 0 {
		err = fmt.Errorf("Failed to compile the application: %s", buffer.String())
	}
	return
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
			cmd := exec.Command(self.path)
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
	if err := self.install(); err != nil {
		panic(fmt.Errorf("Failed to rebuild the application: %v", err))
	}
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
		os.Remove(self.path)
		os.Exit(1)
	}()

	// start waiting the signal to start running.
	var gorun = self.run()
	if err := self.install(); err == nil {
		gorun <- true
	}

	wd := fs.Watchdog(self.pkg.Dir)
	wd.Add(watchList, func(event *fsnotify.Event) {
		self.rerun(gorun)
	})
	wd.Start()
}

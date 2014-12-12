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
	"bytes"
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/go-fsnotify/fsnotify"
)

var WatchList = []string{"*.go", "*.html"}

type binary struct {
	name string
	path string
}

type App struct {
	bin *binary
	pkg *build.Package

	installed bool

	daemon *fsnotify.Watcher
}

func New(path string) *App {
	pkg, err := build.ImportDir(path, build.AllowBinary)
	if err != nil {
		panic(err)
	}
	app := new(App)
	app.pkg = pkg
	// binary file settings
	app.bin = new(binary)
	app.bin.name = "demo"
	if gobin := os.Getenv("GOBIN"); gobin != "" {
		app.bin.path = filepath.Join(gobin, app.bin.name)
	} else {
		app.bin.path = filepath.Join(app.pkg.BinDir, app.bin.name)
	}
	return app
}

// install compiled the application binary package into binary root path.
func (self *App) install() (err error) {
	log.Printf("Building binary package...")
	cmd := exec.Command("go", "get", self.pkg.ImportPath)

	buffer := bytes.NewBuffer([]byte{})
	cmd.Stdout = buffer
	cmd.Stderr = buffer

	err = cmd.Run()

	if buffer.Len() > 0 {
		err = fmt.Errorf("Failed to compiled the application: %s", buffer.String())
	} else {
		self.installed = true
	}

	return
}

// run executes the runnerable executable under package binary root.
func (self *App) run() (gorun chan bool) {
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
			cmd := exec.Command(self.bin.path)
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

// listen removes binary package when application stopped.
func (self *App) listen() {
	channel := make(chan os.Signal, 2)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-channel
		os.Remove(self.bin.path)
		os.Exit(1)
	}()
}

// watch recursively addes all folders under package directory into watch list.
func (self *App) watch() {
	if self.daemon != nil {
		self.daemon.Close()
	}
	// ensure we have a new daemon watcher eachtime we start watching.
	self.daemon, _ = fsnotify.NewWatcher()

	// watch all folders under package directory.
	filepath.Walk(self.pkg.Dir, func(path string, info os.FileInfo, e error) error {
		if info.IsDir() && !strings.HasPrefix(info.Name(), ".") {
			self.daemon.Add(path)
		}
		return e
	})
}

// inWatchList checks if the event happened on the listed watched files.
func (self *App) inWatchList(event fsnotify.Event) bool {
	var basename = filepath.Base(event.Name)
	for _, pattern := range WatchList {
		matched, _ := filepath.Match(pattern, basename)
		if matched {
			return true
		}
	}
	return false
}

// Starts activates the application server along with
// a daemon watcher for monitoring the files's changes.
// FIXME real multiple events handling instead of simple delay.
func (self *App) Start() (err error) {
	// start listening to the ctrl-c interruption.
	self.listen()

	// start waiting the signal to start running.
	var gorun = self.run()
	if err = self.install(); err == nil {
		gorun <- true
	}

	// start watching the changes.
	self.watch()
	// multiple events can be triggered on a successful write
	// (e.g. Create followed by multiple CHMOD), just wait for
	// seconds to let it calm before actual processing.
	var delay <-chan time.Time
	for {
		select {
		case event := <-self.daemon.Events:
			if self.inWatchList(event) {
				delay = time.After(3 * time.Second)
			}

		case err := <-self.daemon.Errors:
			log.Fatalf("Failed to watch the application sources: %v", err)

		case <-delay:
			// restart the daemon watcher.
			self.watch()

			err = self.install()
			if !self.installed {
				continue
			}
			gorun <- true
		}
	}
	return
}

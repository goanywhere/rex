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
	"github.com/goanywhere/rex/livereload"

	"github.com/goanywhere/cmd"
	"github.com/goanywhere/env"
	"github.com/goanywhere/fs"
)

var (
	port      int
	watchList = regexp.MustCompile(`\.(go|html|atom|rss|xml)$`)
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
	cmd.Loading(done)

	// * try build the application into rex-bin(.exe)
	command := exec.Command("go", "build", "-o", self.binary)
	command.Dir = self.dir
	if e := command.Run(); e != nil {
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
			command := exec.Command(self.binary, fmt.Sprintf("--port=%d", port))
			command.Dir = self.dir
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			if err := command.Start(); err != nil {
				log.Fatalf("Failed to start the process: %v\n", err)
			}
			proc = command.Process
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
	watcher.Add(watchList, func(filename string) {
		relpath, _ := filepath.Rel(self.dir, filename)
		log.Infof("Changes on %s detected", relpath)
		self.rerun(gorun)
	})
	watcher.Start()
}

// Run creates an executable application package with livereload supports.
func Run(ctx *cli.Context) {
	port = ctx.Int("port")

	if len(ctx.Args()) == 1 {
		cwd = ctx.Args()[0]
	}
	if abspath, err := filepath.Abs(cwd); err == nil {
		env.Set(internal.ROOT, abspath)
	} else {
		log.Fatalf("Failed to retrieve the directory: %v", err)
	}

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

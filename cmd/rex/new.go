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

	"github.com/goanywhere/cmd"
	"github.com/goanywhere/crypto"
)

const endpoint = "https://github.com/goanywhere/rex"

var secrets = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*(-_+)")

type project struct {
	name string
	root string
}

func (self *project) create() {
	cmd.Prompt("Fetching project template\n")
	var done = make(chan bool)
	cmd.Loading(done)

	command := exec.Command("git", "clone", "-b", "scaffolds", endpoint, self.name)
	command.Dir = cwd
	if e := command.Run(); e == nil {
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
		cmd.Prompt("Fetching project dependencies\n")
		command := exec.Command("npm", "install")
		command.Dir = self.root
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		command.Run()
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

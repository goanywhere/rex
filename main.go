package main

import (
	"github.com/goanywhere/rex"

	"./apps"
)

func main() {
	rex.Get("/", apps.Index)
	rex.FileServer("/static/", "build")
	rex.Run()
}

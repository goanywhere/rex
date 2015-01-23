package main

import (
	"github.com/goanywhere/rex"
	"github.com/goanywhere/rex/modules"
)

func index(ctx *rex.Context) {
	ctx.HTML("index.html")
}

func main() {
	rex.Use(modules.LiveReload)
	rex.Get("/", index)
	rex.FileServer("build")
	rex.Run()
}

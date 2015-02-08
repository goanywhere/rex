package main

import (
	"github.com/goanywhere/rex"
	"github.com/goanywhere/rex/web"
)

func index(ctx *web.Context) {
	ctx.Set("echo", "Welcome")
	ctx.Render("index.html")
}

func main() {
	rex.Get("/", index)
	rex.FileServer("/static/", "build")
	rex.Run()
}

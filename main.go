package main

import "github.com/goanywhere/rex"

func index(ctx *rex.Context) {
	ctx.HTML("index.html")
}

func main() {
	rex.Get("/", index)
	rex.FileServer("/static/", "build")
	rex.Run()
}

package main

import (
	"github.com/goanywhere/rex"
	. "github.com/goanywhere/rex/context"
)

func index(ctx *Context) {
	ctx.HTML("index.html")
}

func main() {
	server := rex.Defaults()
	server.Get("/", index)
	server.Run()
}

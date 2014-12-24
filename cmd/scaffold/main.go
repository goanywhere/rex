package main

import (
	"github.com/goanywhere/rex"
	"github.com/goanywhere/rex/web"
)

func index(ctx *web.Context) {
	ctx.HTML("index.html")
}

func main() {
	server := rex.New()
	server.Get("/", index)
	server.Run()
}

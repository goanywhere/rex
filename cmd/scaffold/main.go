package main

import (
	"github.com/goanywhere/rex"
)

func index(ctx *rex.Context) {
    ctx.HTML("index.html")
}

func main() {
	app := rex.New()
	app.Get("/", index)
	app.Serve()
}

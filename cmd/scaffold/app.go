package main

import (
	"github.com/goanywhere/rex"
	. "github.com/goanywhere/rex/context"
	"github.com/goanywhere/rex/middleware"
)

func index(ctx *Context) {
	ctx.HTML("index.html")
}

func main() {
	app := rex.Defaults()
	app.Use(middleware.Static("build"))
	app.Get("/", index)
	app.Run()
}

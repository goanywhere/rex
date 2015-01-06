package main

import (
	"github.com/goanywhere/rex"
	"github.com/goanywhere/rex/modules"
	. "github.com/goanywhere/rex/context"	
)

func index(ctx *Context) {
	ctx.HTML("index.html")
}

func main() {
	app := rex.New()
	app.Use(modules.LiveReload)
	app.Use(modules.Static(modules.Options{
		"URL": "/static",
		"Dir": "assets",
	}))
	app.Get("/", index)
	app.Run()
}

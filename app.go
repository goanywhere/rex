package main

import (
	"github.com/goanywhere/rex"
	"github.com/goanywhere/rex/modules"
	"github.com/goanywhere/rex/web"	
)

func index(ctx *web.Context) {
	ctx.HTML("index.html")
}

func main() {
	app := rex.New()
	app.Use(modules.LiveReload)
	app.Use(modules.Static(modules.Options{
		"URL": "/static",
	}))
	app.Get("/", index)
	app.Run()
}

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
	rex.Get("/", index)
	rex.Use(modules.XSRF)
	rex.Use(modules.Static)
	rex.Use(modules.LiveReload)
	rex.Run()
}

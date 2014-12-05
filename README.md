Web.*go*
======

Web.*go* is a powerful starter kit for modular web applications/services in Golang.

## Getting Started

Install the package (**go 1.3** and greater is required):

```shell
$ go get -v github.com/goanywhere/web
```


## Features
* Flexible Env-based configurations.
* Non-intrusive/Modular design, extremely easy to use.
* Awesome routing system provided by [Gorilla/Mux](http://www.gorillatoolkit.org/pkg/mux).
* Flexible middleware system based on [http.Handler](http://godoc.org/net/http#Handler) interface.
* Works nicely with other Golang packages.
* **Fully compatible with the [http.Handler](http://godoc.org/net/http#Handler)/[http.HandlerFunc](http://godoc.org/net/http#HandlerFunc) interface.**


After installing Go and setting up your [GOPATH](http://golang.org/doc/code.html#GOPATH), create your first server, we named it `server.go` here.

``` go
package main

import (
    "github.com/goanywhere/web"
)

func main() {
    app := web.New()
    app.Get("/", func(w http.ResponseWriter, r *http.Request) {
        ctx := web.NewContext(w, r)
        ctx.String("Hello World")
    })
    app.Get("/hello", func(ctx *web.Context) {
        ctx.String("Hello Again")
    })
    app.Serve()
}
```

Then start your server:
``` sh
go run server.go
```

You will now have a HTTP server running on `localhost:5000`.


## Template

The standard template (html/template) package implements data-driven templates for generating HTML output safe against code injection, sounds nice? But once you step into the real world, you will soon find your code to be spaghetti. To parse multiple files with pieces of "define", say you have a "index.html", and header source defined in "header.html", footer source in "footer.html", you will need this:

```go
template.Must(template.ParseFiles("index.html", "header.html", "footer.html"))
```

What if another page say "contact.html" will share the same header & footer? Oops & yes, you'll need to do this again,

```go
template.Must(template.ParseFiles("contact.html", "header.html", "footer.html"))
```

Inheritance? They are pretty much the same, yes, you'll have to do this over & over again like this:

```go
template.Must(template.ParseFiles("layout.html", "index.html", "header.html", "footer.html"))

template.Must(template.ParseFiles("layout.html", "contact.html", "header.html", "footer.html"))
```

Web.*go* solution? Simple, in addtition to the standard tags, we introduce two "new" (not really if you have ever used Django/Tornado/Jinja/Liquid) tags, "extends" & "include". You simply add the these two into the html pages as previous, the code will then will be like:

```go
import "github.com/goanywhere/web/template"

loader := template.NewLoader("templates")
template := loader.Parse("index.html")
```

There you Go now, simple as that.


## Context

Context is a very useful helper shipped with Web.*go*. It allows you to access incoming requests & responsed data, there are also shortcuts for rendering HTML/JSON/XML.


``` go
package main

import (
    "github.com/goanywhere/web"
)

func index (ctx *web.Context) {
    ctx.HTML("index.html")  // Context.HTML has the extends/include tag supports by default.
}

func json (ctx *web.Context) {
    ctx.JSON(web.H{"data": "Hello Web", "success": true})
}

func main() {
    app := web.New()
    app.Get("/", index)
    app.Get("/api", json)
    app.Serve()
}
```


## Settings

All settings on Web.*go* utilize system evironment via `os.Environ`. By using this approach you can compile your own settings files into the binary package for deployment without exposing the sensitive settings, it also makes configuration extremly easy & flexible via both command line & application.

``` go
package main

import (
    "github.com/goanywhere/env"
    "github.com/goanywhere/web"
)

func index (ctx *web.Context) {
    ctx.HTML("layout.html", "index.html", "header.html")
}

func main() {
    // Override default 5000 port here.
    env.Set("port", "9394")

    app := web.New()
    app.Get("/", index)
    app.Serve()
}
```

You will now have the HTTP server running on `0.0.0.0:9394`.

`web.Env` also supports custom struct for you to access the reflected values (the key is case-insensitive).

``` go
package main

import (
    "fmt"

    "github.com/goanywhere/env"
)

type Spec struct {
    App string
}

func main() {
    var spec Spec

    env.Set("app", "myapplication")
    env.SetSpec(&spec)

    fmt.Printf("App: %s", spec.App)     // output: "App: myapplication"
}
```

We also includes dotenv supports:

``` text
test1  =  value1
test2 = 'value2'
test3 = "value3"
export test4=value4
```

``` go
package main

import (
    "fmt"

    "github.com/goanywhere/env"
)

func main() {
    // This will load '.env' from current working directory (enabled by Web.go by default)
    // Use env.Set("root", "<Other Dir. You Want>") to initiate different root path for .env.
    env.Load()

    fmt.Printf("<test: %s>", env.Get("test"))     // output: "value"
    fmt.Printf("<test2: %s>", env.Get("test2"))   // output: "value2"
}
```


Hey, dude, why not just use those popular approaches, like file-based config? We know you'll be asking & we have the answer as well, [here](//12factor.net/config).


## Middleware

Middleware works between http request and the router, they are no different than the standard http.Handler. Existing middlewares from other frameworks like logging, authorization, session, gzipping are very easy to integrate into Web.*go*. As long as the middleware comply the `web.Middleware` interface (shorcut to standard `func(http.Handler) http.Handler`), you can simply add one like this:

``` go
app.Use(middleware.XSRF)
```


Since the middleware is just the standard http.Handler, writing a custom middleware is also pretty straightforward:

``` go
app.Use(func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
        web.Debug("Custom Middleware Started")
        next.ServeHTTP(writer, request)
        web.Debug("Custom Middleware Ended")
    })
})
```

## Live Reload

[Fresh](https://github.com/pilu/fresh) backs your Web.*go* application right out of the box.

To get started with Live Reload, you need to install the command line application first.
*NOTE* Please make sure you have $GOPATH & $GOBIN correctly set.

``` sh
$ cd $GOPATH/src/github.com/goanywhere/web/cmd
$ go build -o $GOPATH/bin/web
```

All set, you are good to Go now, simple as that:

``` sh
$ web serve
```



## Frameworks comes & dies, will this be supported?

Positive! Web.*go* is an internal/fundamental project at GoAnywhere. We developed it and we are going to continue using/improving it.


##Roadmap for v1.0


- [X] Sharding Supports
- [X] Env-Based Configurations
- [X] Project Home page
- [X] Test Suite
- [X] New Project Template
- [X] CLI Apps Integrations 
- [X] Improved Template Rendering
- [ ] Live Reload Integration
- [ ] Template Functions
- [ ] i18n Supports
- [ ] More Middlewares
- [ ] Form Validations
- [ ] Project Wiki
- [ ] Continuous Integration
- [ ] Performance Boost
- [ ] Stable API

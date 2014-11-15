Web.*go*
======

Web.*go* is a powerful starter kit for modular web applications/services in Golang.


## NOTE

This is a ongoing project at experiemental stage, consider it's version *ZERO* and *NOT* suitable for production usage yet.


## Features
* Flexible Env-based configurations.
* Non-intrusive/Modular design, extremely easy to use.
* Awesome routing system provided by [Gorilla/Mux](http://www.gorillatoolkit.org/pkg/mux).
* Flexible middleware system based on [http.Handler](http://godoc.org/net/http#Handler) interface.
* Works nicely with other Golang packages.
* **Fully compatible with the [http.Handler](http://godoc.org/net/http#Handler)/[http.HandlerFunc](http://godoc.org/net/http#HandlerFunc) interface.**


## Getting Started

Install the package (**go 1.3** and greater is required):

~~~
go get github.com/goanywhere/web
~~~


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


## Context

Context is a very useful helper shipped with Web.*go*. It allows you to access incoming requests & responsed data, there are also shortcuts for rendering HTML/JSON/XML.


``` go
package main

import (
    "github.com/goanywhere/web"
)

func index (ctx *web.Context) {
    ctx.HTML("layout.html", "index.html", "header.html")
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

All settings on Web.*go* utilize system evironment via `os.Environ`. By using this approach you can compile your own settings files into the binary package for deployment without exposing the sensitive settings, this also make configuration extremly easy & flexible via both command line & application.

``` go
package main

import (
    "github.com/goanywhere/web"
)

func index (ctx *web.Context) {
    ctx.HTML("layout.html", "index.html", "header.html")
}

func main() {
    // Override default 5000 port here.
    web.Env.Set("port", "9394")

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

    "github.com/goanywhere/web"
)

type Spec struct {
    App string
}

func main() {
    web.Env.Set("app", "myapplication")

    var spec Spec
    web.Env.Load(&spec)
    fmt.Printf("App: %s", spec.App)     // output: "App: myapplication"
}
```


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



## Frameworks comes & dies, will this be supported?

Positive! Web.*go* is an internal & fundamental project at GoAnywhere. We developed it and we are going to continue using/improving it.


##Roadmap for v1.0


- [X] Sharding Supports
- [X] Live Reload Integration
- [X] Env-Based Configurations
- [ ] Improved Template Rendering
- [ ] Improved Logging System
- [ ] i18n Supports
- [ ] Template Functions
- [ ] More Middlewares
- [ ] Command-Line Apps
- [ ] Validations
- [ ] Test Suite
- [ ] Documentation
- [ ] Home page
- [ ] Continuous Integration
- [ ] Performance Boost
- [ ] Stable API

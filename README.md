<a href="#"><img alt="rex" src="https://raw.githubusercontent.com/go-rex/rex/assets/images/rex.png" width="160px" height="64px"></a>
===

Rex is a powerful framework for modular web development in Golang.

## Getting Started

Install the package, along with executable binary helper (**go 1.4** and greater is required):

```shell
$ go get -v github.com/goanywhere/rex/...
```

## Features
* Flexible Env-based configurations.
* Awesome routing system provided by [Gorilla/Mux](//github.com/gorilla/mux).
* Django-syntax like templating-language backed by [flosch/pongo](//github.com/flosch/pongo2).
* Non-intrusive/Modular design, extremely easy to use.
* Standard & modular system based on [http.Handler](http://godoc.org/net/http#Handler) interface.
* Command line tools
    * Auto-compile for .go & .html
    * Browser-based Live reload supports
* **Fully compatible with the [http.Handler](http://godoc.org/net/http#Handler)/[http.HandlerFunc](http://godoc.org/net/http#HandlerFunc) interface.**


After installing Go and setting up your [GOPATH](http://golang.org/doc/code.html#GOPATH), create your first server, we named it `server.go` here.

``` go
package main

import (
    "fmt"
    "net/http"

    "github.com/goanywhere/rex"
)

func main() {
    rex.Get("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello World")
    })
    rex.Run()
}
```

Then start your server:
``` shell
rex run
```

You will now have a HTTP server running on `localhost:5000`.



## Context

Context is a very useful helper shipped with Rex. It allows you to access/send incoming requests & responsed data, there are also shortcuts for rendering HTML/JSON/XML.


``` go
package main

import (
    "github.com/goanywhere/rex"
    "github.com/goanywhere/rex/web"
)

func Index(ctx *web.Context) {
    ctx.Render("index.html")
}

func XML(ctx *web.Context) {
    ctx.Render("atom.xml")
}

func JSON(ctx *web.Context) {
    ctx.Send(rex.M{"Success": true, "Response": "This is a JSON Response"})
}

func main() {
    // rex loads the html/xml from `os.Getenv("views")` by default.
    rex.Get("/", Index)
    rex.Get("/api", JSON)
    rex.Get("/xml", XML)
    rex.Run()
}
```


## Settings

All settings on Rex can be accessed via `rex`, which essentially stored in `os.Environ`. By using this approach you can compile your own settings files into the binary package for deployment without exposing the sensitive settings, it also makes configuration extremly easy & flexible via both command line & application.

``` go
package main

import (
    "github.com/goanywhere/rex"
)

func index (ctx *rex.Context) {
    ctx.Render("index.html")
}

func main() {
    // Override default 5000 port here.
    env.Set("PORT", 9394)

    rex.Get("/", index)
    rex.Run()
}
```

You will now have the HTTP server running on `0.0.0.0:9394`.

Hey, dude, why not just use those popular approaches, like file-based config? We know you'll be asking & we have the answer as well, [here](http://12factor.net/config).


## Modules

Modules work between http requests and the router, they are no different than the standard http.Handler. Existing modules (aka. middleware) from other frameworks like logging, authorization, session, gzipping are very easy to integrate into Rex. As long as it complies the `rex.Module` interface (shorcut to standard `func(http.Handler) http.Handler`), you can simply add one like this:

``` go
app.Use(modules.XSRF)
```


Since a module is just the standard http.Handler, writing a custom module is also pretty straightforward:

``` go
app.Use(func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
        log.Printf("Custom Middleware Module Started")
        next.ServeHTTP(writer, request)
        log.Printf("Custom Middleware Module Ended")
    })
})
```

## Frameworks comes & dies, will this be supported?

Positive! Rex is an internal/fundamental project at GoAnywhere. We developed it and we are going to continue using/improving it.


##Roadmap for v1.0


- [X] Env-Based Configurations
- [X] Project Home page
- [X] Test Suite
- [X] New Project Template
- [X] CLI Apps Integrations
- [X] Improved Template Rendering
- [X] Performance Boost
- [X] Hot-Compile Runner
- [X] Live Reload Integration
- [X] Template Functions
- [X] Common Modules
- [X] Cache Framework
- [ ] Project Wiki
- [ ] Continuous Integration
- [ ] Stable API

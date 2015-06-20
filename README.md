<a href="#"><img alt="rex" src="https://raw.githubusercontent.com/go-rex/rex/assets/images/rex.png" width="160px" height="64px"></a>
===

Rex is a powerful toolkit for modular web development in Golang, designed to work directly with net/http.

<img alt="wrk" src="https://raw.githubusercontent.com/goanywhere/rex/assets/images/wrk.png">

## Getting Started

Install the package, along with executable binary helper (**go 1.4** and greater is required):

```shell
$ go get -v github.com/goanywhere/rex/...
```

## Features
* Flexible Env-based configurations.
* Awesome routing system provided by [Gorilla/Mux](//github.com/gorilla/mux).
* Group routing system with middleware modules supports
* Non-intrusive/Modular design, extremely easy to use.
* Standard & modular system based on [http.Handler](http://godoc.org/net/http#Handler) interface.
* Command line tools
    * Auto-compile/reload for .go & .html sources
    * Browser-based Live reload supports for HTML templates
* **Fully compatible with the [http.Handler](http://godoc.org/net/http#Handler)/[http.HandlerFunc](http://godoc.org/net/http#HandlerFunc) interface.**


After installing Go and setting up your [GOPATH](http://golang.org/doc/code.html#GOPATH), create your first server.

``` go
package main

import (
    "io"
    "net/http"

    "github.com/goanywhere/rex"
)

func main() {
    rex.Get("/", func(w http.ResponseWriter, r *http.Request) {
        io.WriteString(w, "Hello World")
    })
    rex.Run()
}
```

Then start your server:
``` shell
rex run
```

You will now have a HTTP server running on `localhost:5000`.




## Settings

All settings on Rex can be accessed via `env`, which essentially stored in `os.Environ`. By using this approach you can compile your own settings files into the binary package for deployment without exposing the sensitive settings, it also makes configuration extremly easy & flexible via both command line & application.

``` go
package main

import (
    "github.com/goanywhere/rex"
    "github.com/goanywhere/x/env"
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

Modules (aka. middleware) work between http requests and the router, they are no different than the standard http.Handler. Existing modules from other frameworks like logging, authorization, session, gzipping are very easy to integrate into Rex. As long as it complies the standard `func(http.Handler) http.Handler` signature, you can simply add one like this:

``` go
app.Use(modules.XSRF)
```


Since a module is just the standard http.Handler, writing a custom module is also pretty straightforward:

``` go
app.Use(func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Custom Middleware Module Started")
        next.ServeHTTP(writer, request)
        log.Printf("Custom Middleware Module Ended")
    })
})
```

Using prefixed (aka. subrouter) router is exactly same as the main one:

```go
app := rex.new()
app.Get("/", func(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "index page")
})

user := app.Group("/users")
user.Use(func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("this is a protected page")
        next.ServeHTTP(writer, request)
    })
})
```


## Frameworks comes & dies, will this be supported?

Positive! Rex is an internal/fundamental project at GoAnywhere. We developed it and we are going to continue using/improving it.


##Roadmap for v1.0


- [X] Env-Based Configurations
- [X] Test Suite
- [X] New Project Template
- [X] CLI Apps Integrations
- [X] Performance Boost
- [X] Hot-Compile Runner
- [X] Live Reload Integration
- [X] Common Modules
- [ ] Full Test Converage
- [ ] Improved Template Rendering
- [ ] Project Wiki
- [ ] Continuous Integration
- [ ] Stable API

<a href="#"><img alt="rex" src="https://raw.githubusercontent.com/goanywhere/rex/assets/images/rex.png" width="160px" height="64px"></a>
===
[![Build Status](https://travis-ci.org/goanywhere/rex.svg?branch=master)](https://travis-ci.org/goanywhere/rex) [![GoDoc](https://godoc.org/github.com/goanywhere/rex?status.svg)](http://godoc.org/github.com/goanywhere/rex)

Rex is a library for modular web development in [Go](http://golang.org/), designed to work directly with net/http.

## Intro

Nah, not another **Web Framework**, we have that enough.The more we spend on [Go](http://golang.org/), the more clearly we realize that most lightweight, pure-stdlib conventions really do scale to large groups of developers and diverse project ecosystems. You absolutely don’t need a *Web Framework* like you normally do in other languages, simply because your code base has grown beyond a certain size. Or you believe it might grow beyond a certain size! You truly ain’t gonna need it. What we really need is just a suitable routing system, along with some common toolkits for web development, the standard idioms and practices will continue to function beautifully at scale.

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
    app := rex.New()
    app.Get("/", func(w http.ResponseWriter, r *http.Request) {
        io.WriteString(w, "Hello World")
    })
    app.Run()
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
    "io"

    "github.com/goanywhere/env"
    "github.com/goanywhere/rex"
)

func index(w http.ResponseWriter, r *http.Request) {
    io.WriteString("Hey you")
}

func main() {
    // Override default 5000 port here.
    env.Set("PORT", 9394)

    app := rex.New()
    app.Get("/", index)
    app.Run()
}
```

You will now have the HTTP server running on `0.0.0.0:9394`.

Hey, dude, why not just use those popular approaches, like file-based config? We know you'll be asking & we have the answer as well, [here](http://12factor.net/config).


## Middleware

Middlware modules work between http requests and the router, they are no different than the standard http.Handler. Existing middleware modules from other frameworks like logging, authorization, session, gzipping are very easy to integrate into Rex. As long as it complies the standard `func(http.Handler) http.Handler` signature, you can simply add one like this:

``` go
app.Use(middleware.XSRF)
```


Since a middleware module is just the standard http.Handler, writing custom middleware is also pretty straightforward:

``` go
app.Use(func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Custom Middleware Module Started")
        next.ServeHTTP(w, r)
        log.Printf("Custom Middleware Module Ended")
    })
})
```

Using prefixed (aka. subrouter) router is exactly same as the main one:

```go
app := rex.New()
app.Get("/", func(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "index page")
})

user := app.Group("/users")
user.Use(func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("this is a protected page")
        next.ServeHTTP(w, r)
    })
})
```

## Benchmark?

Rex is built upon [Gorilla/Mux](//github.com/gorilla/mux), designed to work with starndard `net/http` directly, which means it can run as fast as stdlib can without compromise. Here is a simple [wrk](https://github.com/wg/wrk) HTTP benchmark on a RMBP (2.8 GHz Intel Core i5 with 16GB memory) machine.

<img alt="wrk" src="https://raw.githubusercontent.com/goanywhere/rex/assets/images/wrk.png">


## Frameworks come & die, will this be supported?

Positive! Rex is an internal/fundamental project at GoAnywhere. We developed it and we are going to continue using/improving it.


##Roadmap for v1.0

- [X] Env-Based Configurations
- [X] CLI Apps Integrations
- [X] Performance Boost
- [X] Hot-Compile Runner
- [X] Live Reload Integration
- [X] Common Middleware Modules
- [X] Continuous Integration
- [X] Full Test Converage
- [ ] Unified Rendering
- [ ] Project Wiki
- [ ] Stable API

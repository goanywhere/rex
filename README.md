Webapp
======

Webapp is a powerful starter kit for modular web applications/services in Golang.


## Features
* Modular design, extremely easy to use.
* File-based configurations with YAML, TOML or JSON supports.
* Non-intrusive design, yet a mature scaffold is still optional.
* Awesome routing system provided by [Gorilla/Mux](http://www.gorillatoolkit.org/pkg/mux).
* Flexible middleware system based on [http.Handler](http://godoc.org/net/http#Handler) interface.
* Works nicely with other Golang packages.
* **Fully compatible with the [http.Handler](http://godoc.org/net/http#Handler)/[http.HandlerFunc](http://godoc.org/net/http#HandlerFunc) interface.**


## Getting Started

Install the webapp package (**go 1.3** and greater is required):

~~~
go get github.com/goanywhere/webapp
~~~


After installing Go and setting up your [GOPATH](http://golang.org/doc/code.html#GOPATH), create your first server, we named it `server.go` here.

``` go
package main

import (
    "fmt"
    "net/http"

    "github.com/goanywhere/webapp"
)

func main() {
    app := webapp.New()
    app.GET("/", func(writer http.ResponseWriter, request *http.Request) {
        fmt.Fprint(writer, "Hello World! ")
    })
    app.Serve()
}
```

Then start your server:
``` sh
go run server.go
```

You will now have a HTTP server running on `localhost:9394`.


## Context

Context is a very useful helper shipped with Webapp. It allows you to access incoming requests & responsed data, there are also shortcuts for rendering HTML/JSON/XML.


``` go
package main

import (
    "net/http"

    "github.com/goanywhere/webapp"
)

func index (writer http.ResponseWriter, request *http.Request) {
    context := webapp.Context(writer, request)
    context.Options.Layout = "layout.html"
    context.HTML(http.StatusOK, "index.html", "header.html")
}

func json (writer http.ResponseWriter, request *http.Request) {
    context := webapp.Context(writer, request)
    context.JSON(http.StatusOK, webapp.H{"data": "Hello Webapp", "success": true})
}

func main() {
    app := webapp.New()
    app.GET("/", index)
    app.GET("/api", json)
    app.Serve()
}
```



## Middleware

Middleware works between http request and the router, they are no different than the standard http.Handler. Existing middlewares from other frameworks like logging, authorization, session, gzipping are very easy to integrate into webapp. As long as the middleware comply the `webapp.Middleware` interface (which is pretty standard), you can simply add one like this:

``` go
app.Use(SampleMiddleware())
```


Since the middleware is just the standard http.Handler, writing a custom middleware is also pretty straightforward:

``` go
app.Use(func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
        fmt.Fprint(writer, "Custom Middleware Started")
        next.ServeHTTP(writer, request)
        fmt.Fprint(writer, "Custom Middleware Ended")
    })
})
```

``` go
func CustomMiddleware() webapp.Middleware {
    return func (next http.Handler) http.Handler {
        return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
            fmt.Fprint(writer, "Custom Middleware Started")
            next.ServeHTTP(writer, request)
            fmt.Fprint(writer, "Custom Middleware Ended")
        }
    }
}

app.Use(CustomMiddleware())
```

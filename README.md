Webapp
======

Webapp is a powerful starter kit for modular web applications/services in Golang.


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

# Goshi

Goshi is an extension of http.ResponseWriter written in Go.

[![license](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://github.com/nanoninja/goshi/blob/master/LICENSE)   [![godoc](https://godoc.org/github.com/nanoninja/goshi?status.svg)](https://godoc.org/github.com/nanoninja/goshi)
[![build Status](https://travis-ci.org/nanoninja/goshi.svg)](https://travis-ci.org/nanoninja/goshi)
[![Coverage Status](https://coveralls.io/repos/github/nanoninja/goshi/badge.svg?branch=master)](https://coveralls.io/github/nanoninja/goshi?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/nanoninja/goshi)](https://goreportcard.com/report/github.com/nanoninja/goshi)  [![codebeat](https://codebeat.co/badges/58e89ce4-2fd8-4a93-b624-afdbbb44a6e3)](https://codebeat.co/projects/github-com-nanoninja-goshi)

## Installation

    go get github.com/nanoninja/goshi

## Getting Started

After installing Go and setting up your
[GOPATH](http://golang.org/doc/code.html#GOPATH), create your first `.go` file.

``` go
package main

import (
	"fmt"
	"net/http"

	"github.com/nanoninja/goshi"
)

func HomeHandler(rw http.ResponseWriter, r *http.Request) {
    w := rw.(goshi.ResponseWriter)

    if !w.Written() {
        fmt.Fprintf(rw, "Welcome to the home page!\n")
    }

    fmt.Fprintf(rw, "Status : %d\n", w.Status())
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", HomeHandler)
    http.ListenAndServe(":3000", goshi.Middleware(mux))
}
```

## License

Goshi is licensed under the Creative Commons Attribution 3.0 License, and code is licensed under a [BSD license](https://github.com/nanoninja/goshi/blob/master/LICENSE).

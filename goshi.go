// Copyright 2016 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package goshi is an extension of http.ResponseWriter written in Go.
package goshi

import (
	"bufio"
	"io"
	"net"
	"net/http"
)

var _ http.ResponseWriter = (*responseWriter)(nil)

// A ResponseWriter interface is used by an HTTP handler to
// construct an HTTP response. It wraps http.ResponseWriter.
//    func HomeHandler(rw http.ResponseWriter, r *http.Request) {
//        w := rw.(goshi.ResponseWriter)
//
//        if !w.Written() {
//            fmt.Fprintf(rw, "Welcome to the home page!\n")
//        }
//
//        fmt.Fprintf(rw, "Status : %d\n", w.Status())
//    }
//
//    func main() {
//        mux := http.NewServeMux()
//        mux.HandleFunc("/", HomeHandler)
//        http.ListenAndServe(":3000", goshi.Middleware(mux))
//    }
type ResponseWriter interface {
	http.ResponseWriter
	http.CloseNotifier
	http.Flusher
	http.Hijacker

	// Returns the number of bytes already written into the response http body.
	Len() int

	// Returns the HTTP response status code of the current request.
	Status() int

	// Returns true if the response body was already written.
	Written() bool

	// WriteString writes the string into the response body.
	WriteString(string) (int, error)
}

type responseWriter struct {
	http.ResponseWriter
	status int
	length int
}

// New returns an instance of ResponseWriter.
func New(rw http.ResponseWriter) ResponseWriter {
	return &responseWriter{
		ResponseWriter: rw,
	}
}

func (rw *responseWriter) Write(b []byte) (n int, err error) {
	rw.writeHeader()
	n, err = rw.ResponseWriter.Write(b)
	rw.length += n
	return
}

func (rw *responseWriter) WriteString(s string) (n int, err error) {
	rw.writeHeader()
	n, err = io.WriteString(rw.ResponseWriter, s)
	rw.length += n
	return
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(rw.status)
}

func (rw *responseWriter) writeHeader() {
	if !rw.Written() {
		rw.WriteHeader(http.StatusOK)
	}
}

func (rw *responseWriter) Len() int {
	return rw.length
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) Written() bool {
	return rw.Status() != 0
}

func (rw *responseWriter) CloseNotify() <-chan bool {
	return rw.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (rw *responseWriter) Flush() {
	rw.ResponseWriter.(http.Flusher).Flush()
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return rw.ResponseWriter.(http.Hijacker).Hijack()
}

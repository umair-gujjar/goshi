// Copyright 2016 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goshi

import (
	"bufio"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	_ ResponseWriter      = &responseWriter{}
	_ http.ResponseWriter = ResponseWriter(&responseWriter{})
	_ http.CloseNotifier  = &responseWriter{}
	_ http.Flusher        = &responseWriter{}
	_ http.Hijacker       = &responseWriter{}
)

func TestResponseWriterWriteString(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := New(rec)
	want := "Lorem ipsum dolor sit amet."

	rw.WriteString(want)

	if got := rec.Body.String(); got != want {
		t.Errorf("WriteString() got %s; want", got)
	}
}

func TestResponseWriterLen(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := New(rec)

	rw.WriteString("Ut enim ad minim veniam.")

	if got, want := rw.Len(), rec.Body.Len(); got != want {
		t.Errorf("Len() got %d; want %d", got, want)
	}
}

func TestResponseWriterWritingHeader(t *testing.T) {
	handler := func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
	}
	rec := httptest.NewRecorder()
	rw := New(rec)

	handler(rw, (*http.Request)(nil))
	if got, want := rw.Status(), 500; got != want && rec.Code != want {
		t.Errorf("Status() = %d; want %d", got, want)
	}
}

func TestResponseWriterDefaultStatus(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := New(rec)

	if status, code := rw.Status(), 0; status != code && rec.Code != code {
		t.Errorf("Status() = %d; want %d", status, code)
	}

	rw.Write([]byte("data"))
	if got, want := rw.Status(), http.StatusOK; got != want && rec.Code != want {
		t.Errorf("Status() = %d; want %d", got, want)
	}
}

type closeNotifyRecorder struct {
	*httptest.ResponseRecorder
	clientClosed chan bool
}

func newCloseNotifyRecorder() *closeNotifyRecorder {
	return &closeNotifyRecorder{
		httptest.NewRecorder(),
		make(chan bool, 1),
	}
}

func (cn *closeNotifyRecorder) closeClient() {
	cn.clientClosed <- true
}

func (cn *closeNotifyRecorder) CloseNotify() <-chan bool {
	return cn.clientClosed
}

func TestResponseWriterCloseNotify(t *testing.T) {
	rec := newCloseNotifyRecorder()
	rw := New(rec)
	cn := rw.(http.CloseNotifier).CloseNotify()
	closed := false

	rec.closeClient()

	select {
	case <-cn:
		closed = true
	case <-time.After(time.Second):
	}
	if closed == false {
		t.Errorf("CloseNotify() got = %v, want %v", closed, true)
	}
}

func TestResponseWriterFlusher(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := New(rec)
	f, ok := rw.(http.Flusher)
	if !ok {
		t.Errorf("rw.(http.Flusher) got %T; want http.Flusher", f)
	}

	rec.Write([]byte("data"))
	f.Flush()

	if flushed := rec.Flushed; !flushed {
		t.Errorf("Flush() got %v; want %v", rec.Flushed, true)
	}
}

type hijackResponseRecorder struct {
	ResponseWriter
}

func newHijackResponseRecorder() *hijackResponseRecorder {
	return &hijackResponseRecorder{
		New(httptest.NewRecorder()),
	}
}

func (h hijackResponseRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, nil
}

func TestResponseWriterHijacker(t *testing.T) {
	rec := newHijackResponseRecorder()
	rw := New(rec)
	conn, buf, err := rw.(http.Hijacker).Hijack()
	if conn != nil && buf != nil && err != nil {
		t.Errorf("Hijack() got %v, %v, %v; want nil, nil, nil", conn, buf, err)
	}
}

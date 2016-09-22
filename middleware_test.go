// Copyright 2016 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goshi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware(t *testing.T) {
	handlerFunc := func(rw http.ResponseWriter, r *http.Request) {
		http.Error(rw, "Something failed", http.StatusInternalServerError)
	}
	handler := Middleware(http.HandlerFunc(handlerFunc))
	rw := httptest.NewRecorder()
	handler.ServeHTTP(rw, (*http.Request)(nil))
	if got, want := rw.Code, http.StatusInternalServerError; got != want {
		t.Errorf("Unexpected Code got = %d; want %d", got, want)
	}
}

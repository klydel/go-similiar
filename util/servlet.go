// Copyright 2013 Adam Peck

package util

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
)

// Replies to http.Request r without a message-body with the response of http.Handler h.
func Head(w http.ResponseWriter, r *http.Request, h http.Handler) {
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, r)
	rr.Flush()

	for k, v := range rr.HeaderMap {
		for _, kv := range v {
			w.Header().Set(k, kv)
		}
	}
	w.WriteHeader(rr.Code)
}

// Returns 405 Method Not Allowed.
func MethodNotAllowed(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

// Return 501 Not Implemented.
func NotImplemented(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

// Reflects http.Request r in the response entity-body.
func Trace(w http.ResponseWriter, r *http.Request) {
	b := bytes.NewBufferString(r.Method)
	b.WriteString(" ")
	b.WriteString(r.URL.String())
	b.WriteString(" ")
	b.WriteString(r.URL.Scheme)
	for k, v := range r.Header {
		b.WriteString("\r\n")
		b.WriteString(k)
		b.WriteString(": ")
		b.WriteString(strings.Join(v, ", "))
	}
	b.WriteString("\r\n")

	w.Header().Set("Content-Type", "message/http")
	w.Write(b.Bytes())
}

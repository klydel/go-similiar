// Copyright 2013 Adam Peck

package util

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func checkResponseRecorder(t *testing.T, rr *httptest.ResponseRecorder, code int, header map[string][]string, body string) {
	if !rr.Flushed {
		t.Fatal("ResponseRecorder not flushed")
	}

	if a, e := rr.Code, code; a != e {
		t.Error(a, "!=", e)
	}
	if a, e := map[string][]string(rr.Header()), header; !reflect.DeepEqual(a, e) {
		t.Error(a, "!=", e)
	}
	if a, e := rr.Body.Bytes(), []byte(body); !bytes.Equal(a, e) {
		t.Error(a, "!=", e)
	}
}

func TestHead(t *testing.T) {
	r, err := http.NewRequest("FOOBAR", "http://foobar/", nil)
	r.Header.Add("Foo", "Bar")
	r.Header.Add("Foo", "Bar")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	Head(rr, r, http.NotFoundHandler())
	rr.Flush()

	checkResponseRecorder(t, rr, http.StatusNotFound, map[string][]string{"Content-Type": []string{"text/plain; charset=utf-8"}}, "")
}

func TestMethodNotAllowed(t *testing.T) {
	rr := httptest.NewRecorder()
	MethodNotAllowed(rr)
	rr.Flush()

	checkResponseRecorder(t, rr, http.StatusMethodNotAllowed, map[string][]string{"Content-Type": []string{"text/plain; charset=utf-8"}}, "Method Not Allowed\n")
}

func TestNotImplemented(t *testing.T) {
	rr := httptest.NewRecorder()
	NotImplemented(rr)
	rr.Flush()

	checkResponseRecorder(t, rr, http.StatusNotImplemented, map[string][]string{"Content-Type": []string{"text/plain; charset=utf-8"}}, "Not Implemented\n")
}

func TestTrace(t *testing.T) {
	r, err := http.NewRequest("FOOBAR", "http://foobar/", nil)
	r.Header.Add("Foo", "Bar")
	r.Header.Add("Foo", "Bar")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	Trace(rr, r)
	rr.Flush()

	checkResponseRecorder(t, rr, http.StatusOK, map[string][]string{"Content-Type": []string{"message/http"}}, "FOOBAR http://foobar/ http\r\nFoo: Bar, Bar\r\n")
}

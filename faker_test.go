package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHTTPCode(t *testing.T) {
	codeInt, codeStr, err := getHTTPCode("404")
	if codeInt != 404 || codeStr != "Not Found" || err != nil {
		t.Errorf("Incorrect! Got %d, %s, %d", codeInt, codeStr, err)
	}
	codeInt, codeStr, err = getHTTPCode("123456")
	if err.Error() != "There's no related http code" {
		t.Errorf("Incorrect! Got %d, %s, %d", codeInt, codeStr, err)
	}
}

func TestHandler(t *testing.T) {
	steps := []struct {
		method     string
		path       string
		data       []byte
		expCode    int
		expData    []byte
		remoteAddr string
	}{
		{
			method:  "GET",
			path:    "/",
			data:    []byte(""),
			expCode: 200,
			expData: []byte("Hello there !"),
		},
		{
			method:  "GET",
			path:    "/echo/blablabla",
			data:    []byte(""),
			expCode: 200,
			expData: []byte("blablabla"),
		},
		{
			method:  "GET",
			path:    "/ip",
			data:    []byte(""),
			expCode: 200,
			expData: []byte("1.2.3.4"),
		},
		{
			method:  "GET",
			path:    "/sleep/1",
			data:    []byte(""),
			expCode: 200,
			expData: []byte("Slept for 1 second(s)"),
		},
		{
			method:  "GET",
			path:    "/sleep/asdf",
			data:    []byte(""),
			expCode: 500,
			expData: []byte("Cannot sleep specified time\n"),
		},
		{
			method:  "GET",
			path:    "/418",
			data:    []byte(""),
			expCode: 418,
			expData: []byte("418 I'm a teapot\n"),
		},
		{
			method:  "POST",
			path:    "/",
			data:    []byte("blablabla"),
			expCode: 200,
			expData: []byte("blablabla"),
		},
		{
			method:  "BLABLABLA",
			path:    "/",
			data:    []byte(""),
			expCode: 501,
			expData: []byte("Not implemented\n"),
		},
	}
	for _, step := range steps {
		req, _ := http.NewRequest(step.method, step.path, bytes.NewReader(step.data))
		req.RemoteAddr = "1.2.3.4:12345"
		rec := httptest.NewRecorder()
		a := &app{code: 200}
		http.HandlerFunc(a.handler).ServeHTTP(rec, req)

		if step.expCode == rec.Code && string(step.expData) == rec.Body.String() {
			t.Logf(`[TEST PASSED] Data: method %v, path %v, data %v
			Expected: code %v, data %v
			Actual: code %v, data %v`, step.method, step.path, string(step.data), step.expCode, string(step.expData), rec.Code, rec.Body.String())
		} else {
			t.Errorf(`[TEST FAILED] Data: method %v, path %v, data %v
Expected: code %v, data "%v"
Actual: code %v, data "%v"`, step.method, step.path, string(step.data), step.expCode, string(step.expData), rec.Code, rec.Body.String())
		}
	}
}

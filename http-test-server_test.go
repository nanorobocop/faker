package main

import "testing"

func TestGetHTTPCode(t *testing.T) {
	codeInt, codeStr, err := getHTTPCode("404")
	if codeInt != 404 || codeStr != "Not Found" || err != nil {
		t.Errorf("Incorrect! Got %d, %s, %d", codeInt, codeStr, err)
	}
}

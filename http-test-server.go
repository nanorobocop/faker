package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	codeStr := r.URL.Path[1:]
	codeInt, err := strconv.Atoi(codeStr)
	if err == nil {
		message := http.StatusText(codeInt)
		if message != "" {
			codeMessage := codeStr + " " + message
			http.Error(w, codeMessage, codeInt)
			return
		}
	}
	fmt.Fprintf(w, "Hello there %s!", r.URL.Path[1:])
}

func handlerEcho(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", strings.TrimLeft(r.URL.Path, "/echo/"))
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/echo/", handlerEcho)
	log.Fatal(http.ListenAndServe(":8090", nil))
}

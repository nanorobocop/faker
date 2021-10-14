package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type app struct {
	code     *int
	resp     *string
	respType *string
}

var (
	port = flag.Int("port", 8080, "port number")
)

func getHTTPCode(s string) (codeInt int, codeStr string, err error) {
	codeInt, err = strconv.Atoi(s)
	if err != nil {
		return 0, "", errors.New("There's no related http code")
	}
	codeStr = http.StatusText(codeInt)
	if codeStr == "" {
		return 0, "", errors.New("There's no related http code")
	}
	return codeInt, codeStr, nil
}

func (a *app) handlerEcho(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.Path, "/")

	if a.code != nil {
		w.WriteHeader(*a.code)
	}
	fmt.Fprintf(w, "%s", strings.Join(urlParts[2:], "/"))
}

func (a *app) handlerIP(w http.ResponseWriter, r *http.Request) {
	remoteAddr, _, _ := net.SplitHostPort(r.RemoteAddr)

	if a.code != nil {
		w.WriteHeader(*a.code)
	}
	fmt.Fprintf(w, "%s", remoteAddr)
}

func (a *app) handlerSleep(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.Path, "/")
	seconds, err := strconv.Atoi(urlParts[2])
	if err != nil {
		http.Error(w, "Cannot sleep specified time", 500)
		return
	}

	time.Sleep(time.Duration(seconds) * time.Second)

	if a.code != nil {
		w.WriteHeader(*a.code)
	}
	fmt.Fprintf(w, "Slept for %v second(s)", seconds)
}

func (a *app) handlerCode(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.Path, "/")
	codeInt, codeStr, err := getHTTPCode(urlParts[1])
	if err != nil {
		http.Error(w, "Failed to parse code, code should be integer", 400)
		return
	}

	codeMessage := strconv.Itoa(codeInt) + " " + codeStr
	http.Error(w, codeMessage, codeInt)
}

func (a *app) handlerHeaders(w http.ResponseWriter, r *http.Request) {
	if a.code != nil {
		w.WriteHeader(*a.code)
	}

	if r.Host != "" {
		fmt.Fprintf(w, "Host: %s\n", r.Host)
	}
	for k, v := range r.Header {
		for _, vv := range v {
			fmt.Fprintf(w, "%s: %s\n", k, vv)
		}
	}
}

func (a *app) handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		urlParts := strings.Split(r.URL.Path, "/")
		if urlParts[1] == "echo" {
			a.handlerEcho(w, r)
			return
		} else if r.URL.Path == "/ip" {
			a.handlerIP(w, r)
			return
		} else if r.URL.Path == "/headers" {
			a.handlerHeaders(w, r)
			return
		} else if len(urlParts) == 3 && urlParts[1] == "sleep" {
			a.handlerSleep(w, r)
			return
		} else if _, err := strconv.Atoi(urlParts[1]); err == nil {
			a.handlerCode(w, r)
			return
		} else {
			if a.respType != nil {
				w.Header().Set("Content-Type", *a.respType)
			}

			if a.code != nil {
				w.WriteHeader(*a.code)
			}

			if a.resp != nil {
				w.Write([]byte(*a.resp))
				return
			}

			fmt.Fprintf(w, "Hello there %s!", r.URL.Path[1:])
		}
	case "POST":
		if a.code != nil {
			w.WriteHeader(*a.code)
		}
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		s := buf.String()
		fmt.Fprintf(w, "%v", s)
	default:
		http.Error(w, "Not implemented", 501)
	}
}

func main() {
	a := &app{
		code:     flag.Int("code", 200, "response code"),
		resp:     flag.String("resp", "", "response content"),
		respType: flag.String("resp-type", "", "response content"),
	}

	flag.Parse()

	http.HandleFunc("/", a.handler)
	fmt.Printf("Running on port :%d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

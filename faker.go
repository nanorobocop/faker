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

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type app struct {
	code     *int
	resp     *string
	respType *string
}

var (
	port     = flag.Int("port", 8080, "port number")
	code     = flag.Int("code", 200, "response code")
	resp     = flag.String("resp", "", "response content")
	respType = flag.String("resp-type", "", "response content")
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
func (a *app) handlerPost(w http.ResponseWriter, r *http.Request) {
	if a.code != nil {
		w.WriteHeader(*a.code)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	s := buf.String()
	fmt.Fprintf(w, "%v", s)
}

func (a *app) handlerDefault(w http.ResponseWriter, r *http.Request) {
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

func (a *app) handlerNotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", 501)
}

type handler struct {
	name    string
	handler http.HandlerFunc
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	LoggerHandler(LoggerSettings{},
		MetricsHandler(h.name, h.handler)).ServeHTTP(w, r)
}

func (a *app) handler(w http.ResponseWriter, r *http.Request) {
	var h handler

	switch r.Method {
	case "GET", "HEAD":
		urlParts := strings.Split(r.URL.Path, "/")
		if urlParts[1] == "echo" {
			h = handler{"echo", a.handlerEcho}
		} else if r.URL.Path == "/ip" {
			h = handler{"ip", a.handlerIP}
		} else if r.URL.Path == "/headers" {
			h = handler{"headers", a.handlerHeaders}
		} else if len(urlParts) == 3 && urlParts[1] == "sleep" {
			h = handler{"sleep", a.handlerSleep}
		} else if _, err := strconv.Atoi(urlParts[1]); err == nil {
			h = handler{"code", a.handlerCode}
		} else {
			h = handler{"default", a.handlerDefault}
		}
	case "POST":
		h = handler{"post", a.handlerPost}
	default:
		h = handler{"not implemented", a.handlerNotImplemented}
	}
	h.ServeHTTP(w, r)
}

func main() {
	a := &app{
		code:     code,
		resp:     resp,
		respType: respType,
	}

	flag.Parse()

	prometheus.MustRegister(inFlightGauge, counter, duration, responseSize)

	http.HandleFunc("/", a.handler)
	http.Handle("/metrics", promhttp.Handler())

	fmt.Printf("Running on port :%d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type LogLevel int

const (
	Metadata LogLevel = iota
	Headers
	Body
)

var (
	maxBodyBytesDefault = 10 * 1024 * 1024
)

type LoggerSettings struct {
	Handler      string
	Level        LogLevel
	MaxBodyBytes *int
}

type LogMessage struct {
	Date            string `json:"date"`
	DurationNs      int    `json:"durationNs"`
	Method          string `json:"method"`
	URL             string `json:"url"`
	Handler         string `json:"handler"`
	ResponseCode    int    `json:"responseCode"`
	RequestHeaders  string `json:"requestHeaders,omitempty"`
	RequestBody     string `json:"requestBody,omitempty"`
	ResponseHeaders string `json:"responseHeaders,omitempty"`
	ResponseBody    string `json:"responseBody,omitempty"`
}

func LoggerHandler(settings LoggerSettings, next http.Handler) http.Handler {
	if settings.MaxBodyBytes == nil {
		settings.MaxBodyBytes = &maxBodyBytesDefault
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := newResponseWriter(w, *settings.MaxBodyBytes)

		start := time.Now()
		next.ServeHTTP(rw, r)
		duration := time.Now().Sub(start)

		log := LogMessage{
			Date:         start.Format(time.RFC3339Nano),
			DurationNs:   int(duration.Nanoseconds()),
			Method:       r.Method,
			URL:          r.URL.String(),
			Handler:      settings.Handler,
			ResponseCode: rw.code,
		}
		bytes, err := json.Marshal(log)
		if err != nil {
			fmt.Printf("Failed to marshal log message: %+v\n", err)
		}

		fmt.Println(string(bytes))
	})
	return handler
}

type responseWriter struct {
	w            http.ResponseWriter
	maxBodyBytes int
	body         *bytes.Buffer
	code         int
	headers      http.Header
}

func newResponseWriter(w http.ResponseWriter, maxBodyBytes int) *responseWriter {
	return &responseWriter{
		w:            w,
		maxBodyBytes: maxBodyBytes,
		body:         &bytes.Buffer{},
		headers:      map[string][]string{},
	}
}

func (rw *responseWriter) Header() http.Header {
	return rw.headers
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	maxWrite := len(b)
	if rw.body.Len()+len(b) > rw.maxBodyBytes {
		maxWrite = rw.maxBodyBytes - rw.body.Len()
	}
	_, err := rw.body.Write(b[:maxWrite])
	if err != nil {
		return 0, err
	}

	return rw.w.Write(b)
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.code = code
	rw.w.WriteHeader(code)
}

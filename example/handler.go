package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ahmdrz/rate-limit"
)

func main() {
	mux := http.DefaultServeMux
	mux.HandleFunc("/", mainHandler)
	limiter := ratelimit.NewHandler(mux, 10, 1*time.Minute, ratelimit.DefaultHandler)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", "8080"),
		Handler: limiter,
	}
	server.ListenAndServe()
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

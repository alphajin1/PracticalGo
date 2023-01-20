package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func apiHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func healthCheckerHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "ok")
}

func setupHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", healthCheckerHandler) // signature : handler func(ResponseWriter, *Request)
	mux.HandleFunc("/api", apiHandler)
}

func main() {
	listenAdder := os.Getenv("LISTEN_ADDR")
	if len(listenAdder) == 0 {
		listenAdder = ":8080" // 127.0.0.1:8080 사용을 권장
	}

	mux := http.NewServeMux() // return *ServeMux
	setupHandlers(mux)
	log.Fatal(http.ListenAndServe(listenAdder, mux))
}

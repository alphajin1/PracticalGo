package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
	"time"
)

func handlePing(w http.ResponseWriter, r *http.Request) {
	log.Println("ping: Got a request")
	fmt.Fprintf(w, "pong")
}

func doSomeWork() {
	time.Sleep(2 * time.Second)
}

func handleUserAPI(w http.ResponseWriter, r *http.Request) {
	log.Println("I started processing the request")

	doSomeWork()

	// 새로운 Request with Context 생성
	req, err := http.NewRequestWithContext(
		r.Context(),
		"GET",
		"http://localhost:8080/ping", nil,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 연습 문제 7.1 어떤 시점이 중단되는지?
	trace := &httptrace.ClientTrace{
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Info: %+v\n", dnsInfo) // (2)
		},
		ConnectStart: func(network, addr string) {
			// TCP 연결 수립이 시작된 경우에 값이 설정됨
			fmt.Printf("Connect Start: %s - %s\n", network, addr) // (3)
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Printf("Got Conn: %+v\n", connInfo) // (4)
		},
		WroteRequest: func(rInfo httptrace.WroteRequestInfo) {
			// HTTP 요청이 완료된 이후에 값이 설정됨
			fmt.Printf("Wrote Request: %+v", rInfo) // (5)
		},
	}
	ctxTrace := httptrace.WithClientTrace(req.Context(), trace)
	req = req.WithContext(ctxTrace)
	// 연습문제 부분 종료

	client := &http.Client{}
	log.Println("Outgoing HTTP request") // (1)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)

	fmt.Fprint(w, string(data))
	log.Println("I finished processing the request")
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	// 3초로 하면 순서를 ClientTrace 순서를 알 수 있음
	// 1초로 하면은 DNSDone 하기도 전에 Request 를 종료함을 알 수 있음
	//timeoutDuration := 1 * time.Second
	timeoutDuration := 3 * time.Second

	userHandler := http.HandlerFunc(handleUserAPI)
	hTimeout := http.TimeoutHandler(
		userHandler,
		timeoutDuration,
		"I ran out of time",
	)

	mux := http.NewServeMux()
	mux.Handle("/api/users/", hTimeout)
	mux.HandleFunc("/ping", handlePing)

	log.Fatal(http.ListenAndServe(listenAddr, mux))
}

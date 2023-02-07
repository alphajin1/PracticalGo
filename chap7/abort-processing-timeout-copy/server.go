package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func handleUserAPI(w http.ResponseWriter, r *http.Request) {
	log.Println("I started processing the request")
	time.Sleep(5 * time.Second)

	log.Println("Before continuing, i will check if the timeout has already expired")
	// 이 부분은 Timeout 으로 인한 HTTP 503 응답 후에 handler function 의 실행을 중단시키기 위한 부분이다.
	// 여기가 없으면, handler function 이 중단되지 않고 끝까지 실행됨
	// nil 이외에 값이 반환되면, client 연결이 끊어졌다는 의미이다.
	// ctx 를 사용하는 것은 클라이언트의 연결 상태를 유지한 채로 요청 처리를 진행하는 것 -> 요청을 처리할지, 말지 결정이 가능하다.
	if r.Context().Err() != nil {
		log.Printf("Aborting further processing: %v\n", r.Context().Err())
		return
	}

	fmt.Fprintf(w, "Hello world!")
	log.Println("I finished processing the request")
}

func main() {
	timeoutDuration := 4 * time.Second

	userHandler := http.HandlerFunc(handleUserAPI)
	// http.Handler 객체를 wrapping 하여 새로운 http.Handler 를 반환하는 미들웨어
	// msg: HTTP 5034 응답과 함께 클라이언트에게 전송될 메세지
	hTimeout := http.TimeoutHandler(userHandler, timeoutDuration, "I ran out of time")

	mux := http.NewServeMux()
	mux.Handle("/api/users/", hTimeout)

	log.Fatal(http.ListenAndServe(":8080", mux))
}

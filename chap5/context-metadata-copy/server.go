package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
)

type requestContextKey struct{} // 외부에서 참고 불가능 (내보내지지 않음)
type requestContextValue struct {
	requestID string
}

func addRequestID(r *http.Request, requestID string) *http.Request {
	c := requestContextValue{
		requestID: requestID,
	}

	/*
		1. 현재의 콘텍스트를 가져와서
		2. 새로운 콘텍스트 부착, Arguments : context, 식별 interface, 데이터 자체의 interface 객체
			WithValue() 규칙
				1) 키로 사용되는 데이터 타입은 반드시 string과 같은 기본 데이터 타입이면 안 됨
				2) 어느 한 패키지에는 키로 사용할 내보내지지 않은(unexported) 커스텀 struct 타입을 정의해야 함
					내보내지지 않은 데이터 타입을 사용함으로써 키가 외부 패키지에서 우연히 중복되지 않음을 확신할 수 있음
				3) 콘텍스트 내에는 요청 스코프 데이터만 저장되어야 함.
	*/
	currentCtx := r.Context()
	newCtx := context.WithValue(currentCtx, requestContextKey{}, c)
	return r.WithContext(newCtx)
}

func logRequest(r *http.Request) {
	ctx := r.Context()
	v := ctx.Value(requestContextKey{})
	if m, ok := v.(requestContextValue); ok {
		log.Printf("Processing request: %s", m.requestID)
	}
}

func processRequest(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	fmt.Fprintf(w, "Request processed")
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	requestID := "request-456-def"
	r = addRequestID(r, requestID)
	processRequest(w, r)
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api", apiHandler)

	log.Fatal(http.ListenAndServe(listenAddr, mux))
}

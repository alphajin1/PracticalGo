package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func longRunningProcess(logWriter *io.PipeWriter) {
	for i := 0; i <= 20; i++ {
		// io.PipeWriter 는 io.Writer 의 구현체
		fmt.Fprintf(logWriter, `{"id": %d, "user_ip", "172.121.19.21", "event": "click_on_add_cart"}`, i)
		fmt.Fprintln(logWriter)
		time.Sleep(1 * time.Second)
	}

	logWriter.Close()
}

func progressStreamer(logReader *io.PipeReader, w http.ResponseWriter, done chan struct{}) {

	buf := make([]byte, 500)
	// interface 구현 여부를 판별
	f, flushSupported := w.(http.Flusher)
	defer logReader.Close()

	w.Header().Set("Content-Type", "text/plain")
	// 사용자에게 데이터를 보여주기 전에 먼저 클라이언트 측에서 데이터를 버퍼링하지 않도록 브라우저에게 알림
	w.Header().Set("X-Content-Type-Options", "nosniff")

	for {
		n, err := logReader.Read(buf)
		if err == io.EOF {
			break
		}

		w.Write(buf[:n])
		if flushSupported {
			// 응답 데이터가 클라이언트에게 사용 가능하게 하기 위해 호출한다.
			f.Flush()
		}
	}

	done <- struct{}{}
}

func longRunningProcessHandler(w http.ResponseWriter, r *http.Request) {
	done := make(chan struct{})
	logReader, logWriter := io.Pipe()
	go longRunningProcess(logWriter)
	go progressStreamer(logReader, w, done)

	<-done
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/job", longRunningProcessHandler)
	err := http.ListenAndServe(listenAddr, mux)
	if err != nil {
		log.Fatalf("Server could not start listening on %s. Error: %v", listenAddr, err)
	}

}

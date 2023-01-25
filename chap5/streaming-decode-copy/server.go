package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type logLine struct {
	UserIP string `json:"user_ip"`
	Event  string `json:"event"`
}

func decodeHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)

	// 연습문제 5.2 엄격한 JSON 디코딩
	dec.DisallowUnknownFields()
	for {

		var l logLine
		err := dec.Decode(&l)
		if err == io.EOF {
			break
		}

		var e *json.UnmarshalTypeError
		if errors.As(err, &e) {
			log.Println(err)
			continue
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(l.UserIP, l.Event)
	}

	fmt.Fprintf(w, "OK")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/decode", decodeHandler)
	http.ListenAndServe(":8080", mux)

	//TEST
	//curl -X POST http://localhost:8080/decode \
	//-d '
	//{"user_ip": "172.121.19.21", "event": "click_on_add_cart"}
	//{"user_ip": "172.121.19.21", "event": "click_on_checkout"}
	//'

	//curl -X POST http://localhost:8080/decode \
	//-d '
	//{"user_ip": "172.121.19.21", "event": "click_on_add_cart", "user_data" : "something"}
	//'
}

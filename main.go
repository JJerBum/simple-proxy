package main

import (
	"fmt"
	"io"
	"net/http"
)

var customTransport = http.DefaultTransport

// 대리 요청
func RequestOnBehalf(w http.ResponseWriter, r *http.Request, targetURL string) {
	// proxy서버에 요청된 Method, URL, Body를 이용해 proxy 요청과 같은 proxy 서버에서 타켓 서버로 요청할  새로운 HTTP 요청을 생성합니다.
	proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Error creating proxy reqeust", http.StatusInternalServerError)
		return
	}

	// 원본 요청의 헤더를 proxyReq으로 복사합니다.
	for name, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}

	fmt.Println(proxyReq)
	// custom transport를 사용하여 proxy reqeust를 요청 보낸다.
	resp, err := customTransport.RoundTrip(proxyReq)
	if err != nil {
		http.Error(w, "Error seding proxy request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 프록시 응답의 헤더를 원본 응답으로 복사합니다.
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// 원본 응답의 상태 코드를 프록시 응답의 상태 코드로 설정합니다.
	w.WriteHeader(resp.StatusCode)

	// 프록시 응답의 본문을 원본 응답에 복사합니다.
	io.Copy(w, resp.Body)
}

func main() {

	http.HandleFunc("/asdf", func(w http.ResponseWriter, r *http.Request) {
		RequestOnBehalf(w, r, "http://localhost:9190/")
	})

	http.ListenAndServe(":10000", nil)
}

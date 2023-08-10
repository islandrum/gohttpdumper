package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type HelloWorldHandler struct{}

func (h HelloWorldHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}

type LoggerMiddleware struct {
	Handler http.Handler
}

func (l LoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestUuid, err := uuid.NewRandom()
	if err != nil {
		log.Fatal("Could not generate uuid for request")
	}
	log.Printf("%s: %s %s %s", requestUuid, r.Method, r.URL.Path, r.Proto)
	for name, values := range r.Header {
		for _, value := range values {
			log.Printf("%s: Header: %s, Value: %s", requestUuid, name, value)
		}
	}
	for name, values := range r.URL.Query() {
		for _, value := range values {
			log.Printf("%s: Param: %s, Value: %s", requestUuid, name, value)
		}
	}
	body, _ := io.ReadAll(r.Body)
	log.Printf("%s: Body: %s", requestUuid, string(body))
	l.Handler.ServeHTTP(w, r)
}

func main() {
	handler := HelloWorldHandler{}
	middleware := LoggerMiddleware{Handler: handler}

	http.Handle("/", middleware)
	http.ListenAndServe(":5555", nil)
}

package main

import (
	"encoding/json"
	"net/http"
)

type ServiceHandler struct {
	helloHandler
	goodbyeHandler
}

type helloHandler interface {
	Hello(w http.ResponseWriter, r *http.Request)
}

type goodbyeHandler interface {
	Goodbye(w http.ResponseWriter, r *http.Request)
}

type httpHelloHandler struct{}
type httpGoodbyeHandler struct{}

func (h *httpHelloHandler) Hello(w http.ResponseWriter, r *http.Request) {
	res := "Hello World"

	jData, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func (h *httpGoodbyeHandler) process(r *http.Request) ([]byte, error) {
	return []byte("\"Peace Out\""), nil
}

func (h *httpGoodbyeHandler) Goodbye(w http.ResponseWriter, r *http.Request) {
	res, err := h.process(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (s *ServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch path {
	case "/hello":
		s.Hello(w, r)
	case "/goodbye":
		s.Goodbye(w, r)
	}
}

func getServiceHandler(hello helloHandler, goodbye goodbyeHandler) *ServiceHandler {
	result := &ServiceHandler{helloHandler: hello, goodbyeHandler: goodbye}
	return result
}

func main() {

	goodbyeHandler := &httpGoodbyeHandler{}
	helloHandler := &httpHelloHandler{}
	service := getServiceHandler(helloHandler, goodbyeHandler)

	http.ListenAndServe(":8080", service)
}

package main

import (
	"encoding/json"
	"net/http"
)

type handler interface {
	validate(params interface{}) (err error)
	prepare(params interface{}) (err error)
	process(params interface{}) (res interface{}, err error)
}

type baseHandler struct {
	handler
}

type helloHandler struct {
	*baseHandler
}

type goodbyeHandler struct {
	*baseHandler
}

type ServiceHandler struct{}

func (h *baseHandler) validate(params interface{}) (err error) {
	return nil
}

func (h *baseHandler) prepare(params interface{}) (err error) {
	return nil
}

func (h *helloHandler) process(params interface{}) (interface{}, error) {
	return "Hello World", nil
}

func (h *goodbyeHandler) process(params interface{}) (interface{}, error) {
	return "Peace Out", nil
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

func (s *ServiceHandler) Hello(w http.ResponseWriter, r *http.Request) {
	handler := newHelloHandler()
	s.handle(w, r, handler)
}

func (s *ServiceHandler) Goodbye(w http.ResponseWriter, r *http.Request) {
	handler := newGoodbyeHandler()
	s.handle(w, r, handler)
}

func newHelloHandler() *helloHandler {
	result := &helloHandler{baseHandler: new(baseHandler)}
	return result
}

func newGoodbyeHandler() *goodbyeHandler {
	result := &goodbyeHandler{baseHandler: new(baseHandler)}
	return result
}

func (s *ServiceHandler) handle(w http.ResponseWriter, r *http.Request, handler handler) {
	if err := handler.validate(r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := handler.prepare(r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := handler.process(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jData, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func main() {
	service := new(ServiceHandler)
	http.ListenAndServe(":8080", service)
}

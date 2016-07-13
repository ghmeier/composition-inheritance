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

type helloHandler struct {
	*handler
}

type goodbyeHandler struct {
	*handler
}

type ServiceHandler struct{}

func (h *handler) validate(params interface{}) (err error) {
	return nil
}

func (h *handler) prepare(params interface{}) (err error) {
	return nil
}

func (h *helloHandler) process(params interface{}) (interface{}, error) {
	return "Hello World", nil
}

func (h *goodbyeHandler) process(params interface{}) (interface{}, error) {
	return "Peace Out", nil
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
	result := &helloHandler{handler: new(handler)}
	return result
}

func newGoodbyeHandler() *goodbyeHandler {
	result := &goodbyeHandler{handler: new(handler)}
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
	http.HandleFunc("/hello", service.Hello)
	http.HandleFunc("/goodbye", service.Goodbye)

	http.ListenAndServe(":8080", nil)
}

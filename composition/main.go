package main

import (
	"encoding/json"
	"net/http"

	"github.com/facebookgo/inject"
)

type ServiceHandler struct {
	httpHelloHandler   `inject:"inline"`
	httpGoodbyeHandler `inject:"inline"`
}

type helloHandler interface {
	Hello(w http.ResponseWriter, r *http.Request)
}

type goodbyeHandler interface {
	GoodBye(w http.ResponseWriter, r *http.Request)
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

func main() {

	serviceHandler := ServiceHandler{}
	helloHandler := httpHelloHandler{}
	goodbyeHandler := httpGoodbyeHandler{}

	injectionGraph := inject.Graph{}
	if err := injectionGraph.Provide(
		&inject.Object{Value: &serviceHandler},
		&inject.Object{Value: &helloHandler},
		&inject.Object{Value: &goodbyeHandler},
	); err != nil {
		//handle error
	}

	if err := injectionGraph.Populate(); err != nil {
		//handle err
	}

	http.HandleFunc("/hello", serviceHandler.Hello)
	http.HandleFunc("/goodbye", serviceHandler.Goodbye)

	http.ListenAndServe(":8080", nil)
}

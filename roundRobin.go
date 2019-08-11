package main

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type UpstreamNumber struct {
	Upstream
	Count int
}

type RoundRobin struct{}

func (rr *RoundRobin) RoundRobinRunner(interf, path, method string, upstream Upstream) {
	r := mux.NewRouter()
	upNum := UpstreamNumber{upstream, 0}
	r.HandleFunc(path, upNum.roundRobinHandle).Methods(method)
	http.ListenAndServe(interf, r)
}

func (upstream *UpstreamNumber) roundRobinHandle(w http.ResponseWriter, r *http.Request) {

	w.Write(sendRequest(upstream.Backends[upstream.Count]))

	upstream.Count = (upstream.Count + 1) % len(upstream.Backends)

}
func sendRequest(url string) []byte {

	response, err := http.Get(url)

	// Check(err)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	// Check(err)
	if err != nil {
		panic(err)
	}
	return body
}

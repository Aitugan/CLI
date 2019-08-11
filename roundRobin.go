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

func (rr *RoundRobin) RoundRobinRunner(muxer *mux.Router,interf, path, method string, upstream Upstream) {
	upNum := UpstreamNumber{upstream, 0}
	muxer.HandleFunc(path, upNum.roundRobinHandle).Methods(method)
}

func (upstream *UpstreamNumber) roundRobinHandle(w http.ResponseWriter, r *http.Request) {

	w.Write(sendRequest(upstream.Backends[upstream.Count]))

	upstream.Count = (upstream.Count + 1) % len(upstream.Backends)

}
func sendRequest(url string) []byte {

	response, err := http.Get(url)

	Check(err)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	Check(err)
	
	return body
}

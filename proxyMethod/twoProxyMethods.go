package proxyMethod

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Config struct {
	Interface string     `json:"interface"`
	Upstreams []Upstream `json:"upstreams"`
}

type Upstream struct {
	Path        string   `json:"path"`
	Method      string   `json:"method"`
	Backends    []string `json:"backends"`
	ProxyMethod string   `json:"proxyMethod"`
}

type UpstreamNumber struct {
	Upstream
	Count int
}

//////////

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func LoadConfiguration(file string) (Config, error) {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	Check(err)
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config, nil
}

func SendRequest(url string) []byte {

	response, err := http.Get(url)

	Check(err)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	Check(err)
	return body
}

//////

func RunServer(filename string) {

	config, err := LoadConfiguration(filename)
	Check(err)
	muxer := mux.NewRouter()

	for i := 0; i < len(config.Upstreams); i++ {
		if config.Upstreams[i].ProxyMethod == "round-robin" {
			RoundRobinRunner(muxer, config.Interface, config.Upstreams[i].Path, config.Upstreams[i].Method, config.Upstreams[i])
		} else {
			AnycastRunner(muxer, config.Interface, config.Upstreams[i].Path, config.Upstreams[i].Method, config.Upstreams[i])
		}
	}
	http.ListenAndServe(config.Interface, muxer)
}

//////

func (upstream Upstream) anycastHandler(w http.ResponseWriter, r *http.Request) {

	mainCH := make(chan []byte, 1)

	for _, backend := range upstream.Backends {
		go func(url string, ch chan<- []byte) {
			ch <- SendRequest(url)
		}(backend, mainCH)
	}
	w.Write(<-mainCH)
}

func AnycastRunner(muxer *mux.Router, interf, path, method string, upstream Upstream) {
	muxer.HandleFunc(path, upstream.anycastHandler).Methods(method)
}

////////

func (upstream *UpstreamNumber) roundRobinHandle(w http.ResponseWriter, r *http.Request) {

	w.Write(SendRequest(upstream.Backends[upstream.Count]))

	upstream.Count = (upstream.Count + 1) % len(upstream.Backends)

}

func RoundRobinRunner(muxer *mux.Router, interf, path, method string, upstream Upstream) {
	upNum := UpstreamNumber{upstream, 0}
	muxer.HandleFunc(path, upNum.roundRobinHandle).Methods(method)
}

// Package proxymethod Runs server, reads JSON configuration file, gets information, checks which proxy method to run on server
package proxymethod

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type Top struct {
	Configure []Config `json:"config"`
}

// Config structure with information from JSON configuration file
type Config struct {
	Interface string     `json:"interface"`
	Upstreams []Upstream `json:"upstreams"`
}

// Upstream structure with information from JSON configuration file
type Upstream struct {
	Path        string   `json:"path"`
	Method      string   `json:"method"`
	Backends    []string `json:"backends"`
	ProxyMethod string   `json:"proxyMethod"`
}

// UpstreamNumber structure with Upstream structure as the first element
type UpstreamNumber struct {
	Upstream
	Count int
}

type ServerNumber struct {
	Top
	CountServ int
}

type Server struct {
	srv             *http.Server
	stopped         bool
	router          *mux.Router
	gracefulTimeout time.Duration
}

var (
	Stopped = false
)

func New(srv *http.Server) *Server {
	router := mux.NewRouter()
	srv.Handler = router
	graceTimeout := 5 * time.Second

	return &Server{
		srv,
		false,
		router,
		graceTimeout,
	}
}

// Check if error exists or not
func Check(err error) {
	if err != nil {
		panic(err)
	}
}

// LoadConfiguration reads JSON configuration file returns configuration struct
func LoadConfiguration(file string) (Top, error) {
	var topConfig Top
	configFile, err := os.Open(file)
	defer configFile.Close()
	Check(err)
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&topConfig)
	return topConfig, nil
}

// SendRequest takes one of provided backends as an argument to http.Get(), takes response and returns Body of response in []byte format
func SendRequest(url string) []byte {
	response, err := http.Get(url)
	Check(err)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	Check(err)
	return body
}

// RunServer gets information from JSON configuration file, calls mux.NewRouter() and checks each element in upstreams array. Then it checks it proxy method and runs server
func RunServer(filename string) {
	topConfig, err := LoadConfiguration(filename)
	Check(err)
	var wg sync.WaitGroup
	for j := 0; j < len(topConfig.Configure); j++ {
		srv := New(&http.Server{
			Addr:         topConfig.Configure[j].Interface,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  10 * time.Second,
		})
		muxer := srv.router
		for i := 0; i < len(topConfig.Configure[j].Upstreams); i++ {
			switch topConfig.Configure[j].Upstreams[i].ProxyMethod {
			case "round-robin":
				RoundRobinRunner(muxer, topConfig.Configure[j].Upstreams[i].Path, topConfig.Configure[j].Upstreams[i].Method, topConfig.Configure[j].Upstreams[i])
			case "anycast":
				AnycastRunner(muxer, topConfig.Configure[j].Upstreams[i].Path, topConfig.Configure[j].Upstreams[i].Method, topConfig.Configure[j].Upstreams[i])
			}
		}
		wg.Add(1)
		go func() {
			srv.ListenAndServe()
			wg.Done()
		}()
	}
	wg.Wait()

}

func (srv *Server) ListenAndServe() error {
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		if err := srv.Shutdown(); err != nil {
			log.Printf("Error:%v\n", err)
		} else {
			log.Println("Server stopped")
		}
	}()
	return srv.srv.ListenAndServe()
}

func (srv *Server) Shutdown() error {
	Stopped = true
	ctx, cancel := context.WithTimeout(context.Background(), srv.gracefulTimeout)
	defer cancel()
	time.Sleep(srv.gracefulTimeout)
	return srv.srv.Shutdown(ctx)
}

// AnycastHandler sends request to provided backends, gets their HTML source code and writes it to webserver, counts till the server with the next index.
// Each time server reloads takes and writes the first response it gets
func (upstream *Upstream) AnycastHandler(w http.ResponseWriter, r *http.Request) {
	if Stopped {
		w.WriteHeader(503)
		return
	}
	select {
	case <-r.Context().Done():
		w.WriteHeader(503)
	default:

		mainCH := make(chan []byte, 1)
		for _, backend := range upstream.Backends {
			go func(url string, ch chan<- []byte) {
				ch <- SendRequest(url)
			}(backend, mainCH)
		}
		w.Write(<-mainCH)
	}
}

// AnycastRunner runs server with anycast proxy method
func AnycastRunner(muxer *mux.Router, path, method string, upstream Upstream) {
	muxer.HandleFunc(path, upstream.AnycastHandler).Methods(method)
}

// RoundRobinHandle sends request to provided backends, gets their HTML source code and writes it to webserver, counts till the server with the next index.
// Each time server reloads takes and writes next response by queue
func (upstream *UpstreamNumber) RoundRobinHandle(w http.ResponseWriter, r *http.Request) {
	if Stopped {
		w.WriteHeader(503)
		return
	}
	select {
	case <-r.Context().Done():
		w.WriteHeader(503)
	default:
		w.Write(SendRequest(upstream.Backends[upstream.Count]))
		upstream.Count = (upstream.Count + 1) % len(upstream.Backends)

	}
}

// RoundRobinRunner runs server with round-robin proxy method
func RoundRobinRunner(muxer *mux.Router, path, method string, upstream Upstream) {

	upNum := UpstreamNumber{upstream, 0}
	muxer.HandleFunc(path, upNum.RoundRobinHandle).Methods(method)
}

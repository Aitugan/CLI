package proxymethod

// package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	// Server
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
	// isRabbitMQ      bool
}

type RequestSender struct {
	clnt *http.Client
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

func (sr *RequestSender) RequestStatus(method, url, server string) *http.Response {
	request, err := http.NewRequest(method, server+url, nil)
	Check(err)
	response, err := sr.clnt.Do(request)
	Check(err)
	return response
}

// SendRequest takes one of provided backends as an argument to http.Get(), takes response and returns Body of response in []byte format
func SendRequest(url string, method string, ch chan *http.Response) error { //) []byte {
	defer func() {
		if toRecover := recover(); toRecover != nil {
			fmt.Println("Recovered in serve", toRecover)
		}
	}()

	req, err := http.NewRequest(method, url, nil)
	Check(err)
	client := &http.Client{}
	resp, err := client.Do(req)
	Check(err)

	headers, err := json.Marshal(resp.Header)
	Check(err)
	log.Printf("%s %s %s", method, url, headers)

	ch <- resp
	return nil

}

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
	log.Println("shutting down")
	return srv.srv.Shutdown(ctx)
}


func (upstream *Upstream) AnycastHandler(w http.ResponseWriter, r *http.Request) {
	// srv := Server{}
	// if srv.stopped {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in reliable request", r)
		}
	}()
	response := make(chan *http.Response)
	for i := 0; i < 2; i++ {
		go func(upstream *Upstream, ch chan *http.Response, w http.ResponseWriter) {
			if Stopped {
				w.WriteHeader(503)
				return
			}
			select {
			case <-r.Context().Done():
				w.WriteHeader(503)
			default:

				defer func() {
					if toRecover := recover(); toRecover != nil {
						log.Println("Reliably requested", r)
					}
					select {

					case <-time.After(time.Second * 10):

						log.Println("Time out. Nothing in 10 seconds. Anycast")

					}

				}()
				mainCH := make(chan *http.Response)
				for _, backend := range upstream.Backends {
					go func(url string, ch chan *http.Response) {
						SendRequest(url, upstream.Method, ch)
					}(backend, mainCH)
				}
				select {
				case bd := <-mainCH:
					// w.Write([]byte(bd.Body))
					io.Copy(w, bd.Body)
				case <-time.After(time.Second * 10):
					log.Println("Time out error")

				}
			}
		}(upstream, response, w)

	}
}

// AnycastRunner runs server with anycast proxy method
func AnycastRunner(muxer *mux.Router, path, method string, upstream Upstream) {
	muxer.HandleFunc(path, upstream.AnycastHandler).Methods(method)
}


// RoundRobinHandle sends request to provided backends, gets their HTML source code and writes it to webserver, counts till the server with the next index.
// Each time server reloads takes and writes next response by queue
func (upstream *UpstreamNumber) RoundRobinHandle(w http.ResponseWriter, r *http.Request) { //(*http.Response) {
	// if upstream.Upstream.Server.stopped {
	response := make(chan *http.Response)
	for range upstream.Backends {
		go func(upstream *UpstreamNumber, ch chan *http.Response) {
			if Stopped {
				w.WriteHeader(503)
				return
			}
			select {
			case <-r.Context().Done():
				w.WriteHeader(503)
			default:
				// ch := make(chan *http.Response)
				// defer close(ch)
				var wg sync.WaitGroup
				// wg.Add(1)
				SendRequest(upstream.Backends[upstream.Count], upstream.Upstream.Method, ch)
				resp := make(chan *http.Response)
				select {
				case d := <-resp:
					ch <- d
				case <-time.After(time.Second * 10):
					log.Println("Time out. Nothing in 10 seconds. Round-robin")
					wg.Wait()
				}
				upstream.Count = (upstream.Count + 1) % len(upstream.Backends)

			}
		}(upstream, response)

	}
}

// RoundRobinRunner runs server with round-robin proxy method
func RoundRobinRunner(muxer *mux.Router, path, method string, upstream Upstream) {
	upNum := UpstreamNumber{upstream, 0}
	muxer.HandleFunc(path, upNum.RoundRobinHandle).Methods(method)
}

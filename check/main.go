package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

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

type HTMLT struct {
	HTML string
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

func (h *HTMLT) bHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, h.HTML) //"Some data Some data Some data Some data Some data")
}

////////////////////////////////////////

////////////////////////////////////////

func main() {
	filename := "config.json"
	config, err := LoadConfiguration(filename)
	Check(err)
	fmt.Println(config.Upstreams[0].Path)

	// for _, upstream := range config.Upstreams {

	// 	if upstream.ProxyMethod == "round-robin" {
	// 		rr := NewRoundRobin{}
	// 	}

	// }
	one := NewHandler(config.Upstreams[0].Backends[0])
	two := NewHandler(config.Upstreams[0].Backends[1])
	for i := 0; i < len(config.Upstreams[0].Backends); i++ { //_, i := range config.Upstreams[0].Backends {
		handler := NewRoundRobin(one, two) //config.Upstreams[0].Backends[0], config.Upstreams[0].Backends[1])
		server := &http.Server{
			Handler: handler,
			Addr:    config.Upstreams[0].Path,
		}
		server.ListenAndServe()
		req, err := http.NewRequest(config.Upstreams[0].Method, config.Upstreams[0].Backends[i], nil)
		Check(err)
		req.Close = true
		// for backend, _ := range config.Backends {

		// func indexHandler(w http.ResponseWriter, r *http.Request) {
		// 	fmt.Fprintf(w, "<h1>GOLANG</h1>")
		// }

		// func main() {
		// 	http.HandleFunc("/", indexHandler)
		// 	http.ListenAndServe(":8080", nil)
		// }
		htmlHand := HTMLT{}
		htmlHand.HTML = Resp(req)
		// http.Handle("/",Resp(req)
		http.HandleFunc("/", htmlHand.bHandler)
		http.ListenAndServe(":8080", nil)
	}
	// handler := NewRoundRobin(one, two) //config.Upstreams[0].Backends[0], config.Upstreams[0].Backends[1])

	// server := &http.Server{
	// 	Handler: handler,
	// 	Addr:    config.Upstreams[0].Path,
	// }

	// fmt.Println("Server is running\n")
	// fmt.Printf("%v will be used as configuration file. If you changed your opinion, use reload command\n\n", filename)
	// server.ListenAndServe(":8080", nil)
	// for _, backend := range config.Upstreams[0].Backends {

	// 	req, err := http.NewRequest(config.Upstreams[0].Method, backend, nil)
	// 	Check(err)
	// 	req.Close = true
	// 	// for backend, _ := range config.Backends {
	// 	Resp(req)
	// }
	// return nil

}

func Resp(request *http.Request) string {
	resp, err := http.DefaultClient.Do(request)
	// Check(err)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	// // Check(err)
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(b))
	return string(b)
	//# fmt.Printf("Response received %v\n\n", b)
}

func NewRoundRobin(handlers ...http.Handler) *RoundRobin {
	return &RoundRobin{
		Counter:  0,
		Mut:      &sync.Mutex{},
		Handlers: handlers,
	}
}

func NewHandler(msg string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", msg)
	})
	return mux
}

type RoundRobin struct {
	Counter  int
	Mut      *sync.Mutex
	Handlers []http.Handler
}

func (rr *RoundRobin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rr.Mut.Lock()
	defer rr.Mut.Unlock()
	rr.Handlers[rr.Counter].ServeHTTP(w, r)
	rr.Counter = (rr.Counter + 1) % len(rr.Handlers)
}

func (rr *RoundRobin) AddHandler(handler http.Handler) {
	rr.Mut.Lock()
	defer rr.Mut.Unlock()

	rr.Handlers = append(rr.Handlers, handler)
}

/////////////////

/////////////////

/////////////////

/////////////////

/////////////////

/////////////////

/////////////////

/////////////////
/////////////////
/////////////////

/////////////////

/////////////////

/////////////////

/////////////////

/////////////////

/////////////////

/////////////////

/////////////////

/////////////////

// package main

// import (
// 	"fmt"
// 	"net/http"
// 	"html/template"
// )

// type NewsAggPage struct {
// 	Title string
// 	News string
// }

// func indexHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "<h1>GOLANG</h1> %s", r)
// }

// func newsAggHandler(w http.ResponseWriter, r *http.Request) {
// 	p := NewsAggPage(Title:"Amazing News Aggregator", News:"Some news")
// 	t,_ := template.ParseFiles('basictemplating.html')
// 	t.execute(w,p)

// }

// func main() {
// 	http.HandleFunc("/afgg", indexHandler)
// 	http.HandleFunc("/agg/", newsAggHandler)
// 	http.ListenAndServe(":8080", nil)
// }

// package main

// type RoundRobin struct {
// 	Counter  int
// 	Mut      *sync.Mutex
// 	Handlers []http.Handler
// }

// func (rr *RoundRobin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	rr.Mut.Lock()
// 	defer rr.Mut.Unlock()
// 	rr.Handlers[rr.Counter].ServeHTTP(w, r)
// 	rr.Counter = (rr.Counter + 1) % len(rr.Handlers)
// }

// func (rr *RoundRobin) AddHandler(handler http.Handler) {
// 	rr.Mut.Lock()
// 	defer rr.Mut.Unlock()

// 	rr.Handlers = append(rr.Handlers, handler)
// }

// func NewRoundRobin(handlers ...http.Handler) *RoundRobin {
// 	return &RoundRobin{
// 		Counter:  0,
// 		Mut:      &sync.Mutex{},
// 		Handlers: handlers,
// 	}
// }

// func NewHandler(msg string) http.Handler {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintf(w, "Hello, %s!", msg)
// 	})
// 	return mux
// }

// func Resp(request *http.Request) {
// 	resp, err := http.DefaultClient.Do(request)
// 	// Check(err)
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer resp.Body.Close()

// 	b, err := ioutil.ReadAll(resp.Body)
// 	// Check(err)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("Response received %v\n\n", b)
// }

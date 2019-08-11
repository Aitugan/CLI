package main

func AnycastRunner(muxer *mux.Router, interf, path, method string, upstream Upstream) {
	muxer.HandleFunc(path, upstream.anycastHandler).Methods(method)

}

func (upstream Upstream) anycastHandler(w http.ResponseWriter, r *http.Request) {

	mainCH := make(chan []byte, 1)

	for _, backend := range upstream.Backends {
		go func(url string, ch chan){
			ch <- sendRequest(url)
		}(backend,mainCH)//sendRequestToChannel(backend, fibonacciNumbersResultChannel)
	}

	w.Write(<-fibonacciNumbersResultChannel)

}

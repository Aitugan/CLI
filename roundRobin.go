package main

import (
	"net/url"
	"sync"
)

type RoundRobin interface {
	Next() *url.URL
}

type Roundrobin struct {
	urls []*url.URL
	mu   *sync.Mutex
	next int
}

func New(urls []*url.URL) (RoundRobin) {
	return &Roundrobin{
		urls: urls,
		mu:   new(sync.Mutex),
	}
}

func (r *Roundrobin) Next() *url.URL {
	r.mu.Lock()
	sc := r.urls[r.next]
	r.next = (r.next + 1) % len(r.urls)
	r.mu.Unlock()
	return sc
}

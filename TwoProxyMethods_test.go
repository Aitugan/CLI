package proxymethod

import (
	"net/http"
	"testing"
)

type TestCase struct {
	method      string
	url         string
	server      string
	expectedURL []string
}

var reqClient = &RequestSender{
	&http.Client{},
}
var tests = []TestCase{
	{
		"GET",
		"/oses/anycast",
		"http://localhost:8081",
		[]string{"https://www.apple.com/",
			"https://www.microsoft.com/",
			"https://www.ubuntu.com/",
			"https://www.archlinux.org/"},
	},
	{
		"GET",
		"/oses/rr",
		"http://localhost:8081",
		[]string{"https://www.apple.com/",
			"https://www.microsoft.com/",
			"https://www.ubuntu.com/",
			"https://www.archlinux.org/"},
	},
}

// Test Status Code
func TestResponseCodeRequester(t *testing.T) {
	for _, test := range tests {
		response := reqClient.RequestStatus(test.method, test.url, test.server)

		if response.StatusCode != 200 {
			t.Fail()
		}
	}
}

// Test upstream
func TestRequestUrlRequester(t *testing.T) {
	var correct = false
	for _, test := range tests {
		response := reqClient.RequestStatus(test.method, test.url, test.server)
		for _, expected := range test.expectedURL {
			if response.Header.Get("URL") == expected {
				correct = true
				break
			}
		}
		if !correct {
			t.Fail()
		}
	}
}

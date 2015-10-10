package main

import (
	"bytes"
	"net/http"
	"sync"
)

var hitUrl = func(method string, url string, body string, headers http.Header) (resp *http.Response, err error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}

	req.Header = headers
	return client.Do(req)
}

type urlTarget struct {
	method  string
	url     string
	body    string
	headers http.Header
}

func (tg *urlTarget) Hit(tracker *sync.WaitGroup) {
	defer tracker.Done()
	hitUrl(tg.method, tg.url, tg.body, tg.headers)
	//TODO Need to return results from the call but can't just return, need to use a channel
	return
}

type urlIterator struct {
	urls     []string
	position int
}

func (it *urlIterator) Next(continuous bool) bool {
	if it.position >= len(it.urls) {
		return false
	}
	it.position++
	return true
}
func (it *urlIterator) Value() Target {
	// Why position-1?  Because the initial value will be 0, the first call to
	// next will increment to 1 but we wont have got the value from position 0
	// at that point.
	retVal := urlTarget{method: "GET", url: it.urls[it.position-1]}
	return &retVal
}
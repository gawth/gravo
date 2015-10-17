package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
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

func (tg *urlTarget) Hit(tracker *sync.WaitGroup, t Timer, h OutputHandler) {
	defer tracker.Done()
	t.Start()

	res, err := hitUrl(tg.method, tg.url, tg.body, tg.headers)
	t.End()

	//TODO Need to return results from the call but can't just return, need to use a channel
	if err != nil {
		log.Println(err)
		return
	}

	h.DealWithIt(*res)
	log.Println(fmt.Sprintf("Call took %v\n", t.GetTime()))
	return

	//TODO Should return this lot and do it externally
	image, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Println(err)
		return
	}
	//TODO should be timing this as well...again, do it outside of the hit call...
	log.Println(fmt.Sprintf("Got %d bytes, %d meg in %v\n", len(image), len(image)/1024/1024, t.GetTime()))
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

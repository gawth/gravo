package main

import (
	"bytes"
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"
)

type stubOutput struct {
}

func (so *stubOutput) DealWithIt(r http.Response) {
	return
}

type stubTimer struct {
}

func (t *stubTimer) Start() {
	return
}
func (t *stubTimer) End() {
	return
}
func (t *stubTimer) GetTime() time.Duration {
	return time.Second
}

func TestUrlHit(t *testing.T) {
	var called = false
	var expectedUrl = "a url"
	var expectedMethod = "a method"
	var expectedBody = "a body"
	var expectedHeaders = map[string][]string{}
	tracker := &sync.WaitGroup{}

	hitUrl = func(method string, url string, body string, headers http.Header) (resp *http.Response, err error) {
		var response http.Response
		response.Body = nopCloser{bytes.NewBufferString("test data")}

		called = true
		return &response, nil
	}

	target := urlTarget{method: expectedMethod, url: expectedUrl, body: expectedBody, headers: expectedHeaders}

	tracker.Add(1)
	target.Hit(tracker, &stubTimer{}, &stubOutput{})
	tracker.Wait()

}

func TestUrlIterator(t *testing.T) {
	testurls := []string{"one", "two"}
	it := urlIterator{urls: testurls}
	called := 0

	expected := reflect.TypeOf(&urlTarget{})

	for it.Next(false) {
		typ := reflect.TypeOf(it.Value())
		if typ != expected {
			t.Errorf("Expected %v but got %v", expected, typ)
		}
		called++
	}
	if called != len(testurls) {
		t.Errorf("TestUrlIterator: Expected next %v times but got %v", len(testurls), called)
	}

}

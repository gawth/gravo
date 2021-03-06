package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"
)

type stubOutput struct {
	expectedBody []byte
	actualBody   []byte
}

func (so *stubOutput) DealWithIt(r http.Response, t Timer) {
	so.actualBody, _ = ioutil.ReadAll(r.Body)
	r.Body.Close()

	return
}

func (so *stubOutput) LogInfo(s string) {

}
func (so *stubOutput) Start() {

}

type stubTimer struct {
	start     time.Time
	callStart int
	end       time.Time
	callEnd   int
	duration  time.Duration
}

func (t *stubTimer) Start() {
	t.callStart++
	return
}
func (t *stubTimer) End() {
	t.callEnd++
	return
}
func (t *stubTimer) GetDuration() time.Duration {
	return t.duration
}

func (t *stubTimer) GetStart() time.Time {
	return t.start
}

func (t *stubTimer) GetEnd() time.Time {
	return t.end
}

func StubTimer(start time.Time, end time.Time) stubTimer {
	return stubTimer{start: start, end: end, duration: end.Sub(start)}
}

func TestUrlHit(t *testing.T) {
	var called = false
	var expectedURL = "a url"
	var expectedMethod = "a method"
	var expectedBody = "a body"
	var expectedHeaders = map[string][]string{}

	testdata := []byte("Test Data")
	expectedRes := bytes.NewBuffer(testdata)

	tracker := &sync.WaitGroup{}
	timer := stubTimer{}
	outer := stubOutput{testdata, nil}

	hitURL = func(method string, url string, body string, headers http.Header) (resp *http.Response, err error) {
		var response http.Response
		response.Body = nopCloser{expectedRes}

		called = true
		return &response, nil
	}

	target := urlTarget{method: expectedMethod, url: expectedURL, body: expectedBody, headers: expectedHeaders}

	tracker.Add(1)
	target.Hit(tracker, &timer, &outer)
	tracker.Wait()

	if !called {
		t.Errorf("TestUrlHit: Hit not called")

	}
	if timer.callStart != 1 {
		t.Errorf("TestUrlHit: Expected start time to be called once, was called %v", timer.start)
	}
	if timer.callEnd != 1 {
		t.Errorf("TestUrlHit: Expected end time to be called once, was called %v", timer.end)
	}

	if !bytes.Equal(outer.expectedBody, outer.actualBody) {
		t.Errorf("TestUrlHit: Expected output of '%v' but got '%v'", string(outer.expectedBody), string(outer.actualBody))
	}

}

func TestUrlIterator(t *testing.T) {
	testurls := []string{"one", "two"}
	it := urlIterator{urls: testurls}
	called := 0

	expected := reflect.TypeOf(&urlTarget{})

	for it.Next(false) {
		v, _ := it.Value()
		typ := reflect.TypeOf(v)
		if typ != expected {
			t.Errorf("Expected %v but got %v", expected, typ)
		}
		called++
	}
	if called != len(testurls) {
		t.Errorf("TestUrlIterator: Expected next %v times but got %v", len(testurls), called)
	}

}

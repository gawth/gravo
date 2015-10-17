package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func TestSoap(t *testing.T) {
	var expected = "http://testhost:1234/path"
	var actual = ""
	var getSOAPCalled = 0

	callTarget = func(target string, method string, header http.Header, body string) (resp *http.Response, err error) {
		actual = target

		var response http.Response
		response.Body = nopCloser{bytes.NewBufferString("test data")}
		return &response, nil
	}

	getSOAPBody = func(filename string) (string, error) {
		getSOAPCalled++
		return "this is the body", nil
	}
	c := config{Soap: true, SoapFile: "soap.txt", Target: target{Host: "testhost", Port: "1234", Path: "path"}, Rate: runrate{Rrate: 5, Rtype: "S"}}
	doSoap(c)

	if expected != actual {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
	if getSOAPCalled != 1 {
		t.Errorf("getSOAP was called should have been called once, not %d", getSOAPCalled)
	}
}

func TestTemplate(t *testing.T) {
	var expected = "The ip is 9999"
	var actual = ""
	var getSOAPCalled = 0

	callTarget = func(target string, method string, header http.Header, body string) (resp *http.Response, err error) {
		actual = body

		var response http.Response
		response.Body = nopCloser{bytes.NewBufferString("test data")}
		return &response, nil
	}

	getSOAPBody = func(filename string) (string, error) {
		getSOAPCalled++
		return "The ip is {{.ip}}", nil
	}
	// 1.Needs to read in the message body
	// 2.Read in the vars file using the first line as the key and then
	// subsequent lines as values
	// 3.For each value row call the service having combined the body
	// template with the row of data
	//
	// First iteration, hard code a map for use with body

	c := config{Soap: true, SoapFile: "soap.txt", Target: target{Host: "testhost", Port: "1234", Path: "path"}, Rate: runrate{Rrate: 5, Rtype: "S"}}
	doSoap(c)

	if expected != actual {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
	if getSOAPCalled != 1 {
		t.Errorf("getSOAP was called should have been called once, not %d", getSOAPCalled)
	}
}

func TestGetTimeUnit(t *testing.T) {
	cases := []struct {
		flag   string
		result time.Duration
	}{
		{"m", time.Millisecond},
		{"s", time.Second},
		{"M", time.Minute},
		{"H", time.Hour},
		{"x", time.Second},
	}

	for _, c := range cases {
		res := getTimeUnit(c.flag)
		if res != c.result {
			t.Errorf("getTimeUnit failed.  Expecting %v but got %v", c.result, res)
		}
	}
}

type stubTarget struct {
	hits int
}

func (tg *stubTarget) Hit(tracker *sync.WaitGroup, t Timer, h OutputHandler) {
	defer tracker.Done()
	tg.hits++
	fmt.Println("Hit %d", tg.hits)
	return
}

type stubIterator struct {
	current int
	finish  int
	target  Target
}

func (s stubIterator) Value() Target {
	return s.target
}
func (s *stubIterator) Next(forever bool) bool {
	if s.current < s.finish {
		s.current++
		return true
	}
	return false
}

func TestRunLoad(t *testing.T) {

	tmp := stubTarget{hits: 0}
	timer := stubTimer{}
	outer := stubOutput{}

	it := stubIterator{current: 0, finish: 5, target: &tmp}

	c := config{Verbose: true, Target: target{Host: "testhost", Port: "1234", Path: "path"}, Requests: 5, Rate: runrate{Rrate: 1, Rtype: "m"}}

	runLoad(c, &it, &timer, &outer)
	if tmp.hits != it.finish {
		t.Errorf("Expected %d hits but got %d", it.finish, tmp.hits)
	}

}

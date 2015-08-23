package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func TestDoStuffHappyPath(t *testing.T) {
	var called = 0
	var expected = 5

	// Override http interaction
	callTarget = func(target string) (resp *http.Response, err error) {
		var response http.Response
		response.Body = nopCloser{bytes.NewBufferString("test data")}

		called++
		return &response, nil
	}

	c := config{Target: target{Host: "testhost", Port: "1234", Path: "path"}, Requests: 5, Rate: runrate{Rrate: 5, Rtype: "S"}}
	doStuff(c)

	if called != expected {
		t.Errorf("Got %d calls to callTarget, expected %d", called, expected)
	}
}

func TestDoStuffHttpError(t *testing.T) {
	// Override http interaction
	callTarget = func(target string) (resp *http.Response, err error) {
		var response http.Response
		response.Body = nopCloser{bytes.NewBufferString("test data")}

		return nil, errors.New("Expect an error")
	}

	c := config{Target: target{Host: "testhost", Port: "1234", Path: "path"}, Requests: 5, Rate: runrate{Rrate: 5, Rtype: "S"}}
	doStuff(c)

}

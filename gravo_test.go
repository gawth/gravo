package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"reflect"
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

func TestUsingaURLFile(t *testing.T) {
	var expected = []string{"a url", "another URL", "and another"}
	var actual = []string{}

	// Override http interaction
	callTarget = func(target string) (resp *http.Response, err error) {
		// Remember what we've been called for
		actual = append(actual, target)

		var response http.Response
		response.Body = nopCloser{bytes.NewBufferString("test data")}
		return &response, nil
	}

	// Override getUrls to return expected URL
	getUrls = func(filename string) (urls []string, err error) {
		return expected, nil
	}

	// Specif the URL file in the config
	c := config{Target: target{Host: "testhost", Port: "1234", File: "urls.txt"}, Rate: runrate{Rrate: 5, Rtype: "S"}}

	doStuff(c)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %s but got %s", expected, actual)
	}

}

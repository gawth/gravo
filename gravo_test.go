package main

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"testing"
	"time"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

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
	fmt.Printf("Hit %d\n", tg.hits)
	return
}

type stubIterator struct {
	current  int
	finish   int
	target   Target
	valError bool
}

func (s stubIterator) Value() (Target, error) {
	if s.valError {
		return nil, errors.New("This is an error")
	}
	return s.target, nil
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
func TestRunLoadWithErrorFromNext(t *testing.T) {

	tmp := stubTarget{hits: 0}
	timer := stubTimer{}
	outer := stubOutput{}

	it := stubIterator{current: 0, finish: 2, target: &tmp, valError: true}

	c := config{Verbose: true, Target: target{Host: "testhost", Port: "1234", Path: "path"}, Requests: 2, Rate: runrate{Rrate: 1, Rtype: "m"}}

	runLoad(c, &it, &timer, &outer)
	if tmp.hits != 0 {
		t.Errorf("Expected %d hits but got %d", 0, tmp.hits)
	}

}

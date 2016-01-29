package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestStatsHandlerHappyPath(t *testing.T) {
	target := chartHandler{}
	client := &http.Client{}

	data := http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString("Some data")),
	}
	tm := stubTimer{}

	target.Start()
	target.DealWithIt(data, &tm)
	target.DealWithIt(data, &tm)

	req, err := http.NewRequest("GET", "http://localhost:8080/stats", nil)
	if err != nil {
		t.Errorf("TestStatsHandlerHappyPath: unable to create request")
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("TestStatsHandlerHappyPath: Got http error")
	}
	results, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Errorf("TestStatsHandlerHappyPath: Unable to process body")
	}

	var out []string
	err = json.Unmarshal(results, &out)
	if err != nil {
		t.Errorf("TestStatsHandlerHappyPath: Unable to parse json response")
	}

	if len(out) != 2 {
		t.Errorf("TestStatsHandlerHappyPath: Expected '%v' to be length 2 but was length %v", out, len(out))
	}
}

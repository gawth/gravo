package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestStatsHandlerHappyPath(t *testing.T) {
	testfile := "testfile"
	target := ChartHandler(testfile, make(chan bool), NullHandler())

	data := http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString("Some data")),
	}
	tm := stubTimer{}

	target.Start()
	target.DealWithIt(data, &tm)
	target.DealWithIt(data, &tm)

	req, err := http.NewRequest("GET", "http://localhost:8910/stats", nil)
	if err != nil {
		t.Errorf("TestStatsHandlerHappyPath: unable to create request")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("TestStatsHandlerHappyPath: Got http error")
	}
	results, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Errorf("TestStatsHandlerHappyPath: Unable to process body")
	}

	var out []metric
	err = json.Unmarshal(results, &out)
	if err != nil {
		t.Errorf("TestStatsHandlerHappyPath: Unable to parse json response:%v", results)
	}

	if len(out) != 2 {
		t.Errorf("TestStatsHandlerHappyPath: Expected '%v' to be length 2 but was length %v", out, len(out))
	}
	req, err = http.NewRequest("GET", "http://localhost:8910/results/"+testfile, nil)
	if err != nil {
		t.Errorf("TestStatsHandlerHappyPath: unable to create request")
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("TestStatsHandlerHappyPath: Got http error from results")
	}
	if resp.StatusCode != 200 {
		t.Errorf("TestStatsHandlerHappyPath: Got a HTTP %v rather than a 200 from results", resp.StatusCode)
	}
}

func TestStatsHandlerMissingFilename(t *testing.T) {
	if os.Getenv("CALL") == "1" {
		target := ChartHandler("", make(chan bool), NullHandler())
		target.Start()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestStatsHandlerMissingFilename")
	cmd.Env = append(os.Environ(), "CALL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Errorf("Call should have aborted with an error but returned %v", err)
}

func TestStatsHandlerPreLoadedData(t *testing.T) {
	tm := time.Now()
	target := ChartHandler("fred", make(chan bool), NullHandler()).(*chartHandler)

	data := []metric{{tm, 123}, {tm, 345}}
	target.data = data
	testServer := httptest.NewServer(http.HandlerFunc(target.statsHandler))
	defer testServer.Close()

	resp, err := http.Get(testServer.URL)
	if err != nil {
		t.Errorf("TestStatsHandlerPreLoadedData: Test server returned an error")
	}

	results, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Errorf("TestStatsHandlerPreLoadedData: Unable to process body")
	}

	var out []metric
	err = json.Unmarshal(results, &out)
	if err != nil {
		t.Errorf("TestStatsHandlerPreLoadedData: Unable to parse json response")
	}

	if len(out) != 2 {
		t.Errorf("TestStatsHandlerPreLoadedData: Expected '%v' to be length 2 but was length %v", out, len(out))
	}
}

func TestGenerateFilename(t *testing.T) {
	tm, _ := time.Parse("2006-01-02 03:04", "2016-04-04 12:30")
	res := generateFilename(tm)
	tar := "20160404_123000"
	if res != tar {
		t.Errorf("TestGenerateFilename: %v should have been %v", res, tar)
	}

}

func TestDealWithIt(t *testing.T) {
	target := chartHandler{logger: make(chan metric)}
	var response http.Response
	metricReceived := false

	response.Body = nopCloser{bytes.NewBuffer([]byte("Expected Data"))}

	testTime := time.Now()
	timer := StubTimer(testTime, testTime.Add(time.Hour))

	go func() {
		met := <-target.logger
		if met.Datetime != testTime {
			t.Errorf("TestDealWithIt: Expected %v but got %v", testTime, met.Datetime)
		}
		metricReceived = true
	}()
	target.DealWithIt(response, &timer)

	if !metricReceived {
		t.Errorf("TestDealWithIt: Channel didn't fire - no metric data came through")
	}

}

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
	stub := StubHandler().(*stubHandler)
	target := ChartHandler(testfile, make(chan bool), stub)

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
	if stub.startCalled == 0 {
		t.Errorf("TestChartHandlerStartParentCalled: Failed to call Start on the parent")
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
	target := ChartHandler("", make(chan bool), nil).(*chartHandler)
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
func TestDealWithItParentTest(t *testing.T) {
	stub := StubHandler().(*stubHandler)
	target := ChartHandler("", make(chan bool), stub).(*chartHandler)

	expectedData := []byte("Expected Data")
	var response http.Response
	response.Body = nopCloser{bytes.NewBuffer(expectedData)}

	testTime := time.Now()
	timer := StubTimer(testTime, testTime.Add(time.Hour))

	go func() {
		<-target.logger
	}()
	target.DealWithIt(response, &timer)

	if stub.dealCalled == 0 {
		t.Errorf("TestDealWithItParent: Failed to call DealWithIt on the parent")
	}

	result, err := ioutil.ReadAll(stub.savedBody)
	if err != nil {
		t.Errorf("TestDealWithItParent: Failed to read expected data from parent call, err %v", err)
	}
	if !bytes.Equal(result, expectedData) {
		t.Errorf("TestDealWithItParent: Expected '%v' but got '%v'", expectedData, result)
	}

}
func TestChartHandlerLogInfoParent(t *testing.T) {
	stub := StubHandler().(*stubHandler)
	target := ChartHandler("", make(chan bool), stub).(*chartHandler)

	target.LogInfo("Blahh")

	if stub.logCalled == 0 {
		t.Errorf("TestChartHandlerLogInfoParent: Failed to call LogInfo on the parent")
	}
}

func TestChartHandlerParseData(t *testing.T) {
	jsondata := metric{Datetime: time.Now(), Val: 1234}
	rawBytes, _ := json.Marshal(jsondata)
	buf := bytes.NewBuffer(rawBytes)
	buf.Write([]byte("\nand another line plus blanks\n\n\n"))
	buf.Write(rawBytes)
	data := ioutil.NopCloser(buf)

	stub := StubHandler().(*stubHandler)
	target := ChartHandler("", make(chan bool), stub).(*chartHandler)

	target.parseData(data)
	if len(target.data) != 2 {
		t.Errorf("TestChartHandlerParseData: Expected two metrics to be stored, got %v", len(target.data))
	}
}

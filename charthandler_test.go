package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestStatsHandlerMissingFilename(t *testing.T) {
	if os.Getenv("CALL") == "1" {
		target := ChartHandler(make(chan bool), &Results{name: ""}, NullHandler())
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
	target := ResultsServer("fred")

	data := []metric{{tm, 123}, {tm, 345}}
	target.data = Results{name: "Res", metrics: data}
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
	target := ChartHandler(make(chan bool), &Results{}, nil).(*chartHandler)
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
	target := ChartHandler(make(chan bool), &Results{}, stub).(*chartHandler)

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
	target := ChartHandler(make(chan bool), &Results{}, stub).(*chartHandler)

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

	target := ResultsServer("fred")

	target.parseData(data)
	if len(target.data.metrics) != 2 {
		t.Errorf("TestChartHandlerParseData: Expected two metrics to be stored, got %v", len(target.data.metrics))
	}
}

func TestResultsList(t *testing.T) {
	target := ResultsServer("ferd")
	files := ResultsList{}
	files.filelist = []string{"file1", "file2"}
	target.savedResults = files

	testServer := httptest.NewServer(http.HandlerFunc(target.resultsList))
	defer testServer.Close()

	resp, err := http.Get(testServer.URL)
	if err != nil {
		t.Errorf("TestResultsList: Test server returned an error")
	}

	res, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Errorf("TestResultsList: Unable to process body")
	}

	if !strings.Contains(string(res), target.data.Name()) {
		t.Errorf("TestResultsList: Expected to find '%v' in the results:%v", target.data.Name(), string(res))
	}
	for _, str := range files.filelist {
		if !strings.Contains(string(res), str) {
			t.Errorf("TestResultsList: Expected to find '%v' in the results", str)
		}
	}
}

package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestStandardOutputLogInfo(t *testing.T) {
	stubP := StubHandler().(*stubHandler)
	target := StandardOutput(false, nil, stubP)

	target.LogInfo("Blahh")

	if stubP.logCalled == 0 {
		t.Errorf("TestStandardOutputLogInfo: Failed to call LogInfo on the parent")
	}

}
func TestStandardOutputStart(t *testing.T) {
	stubP := StubHandler().(*stubHandler)
	target := StandardOutput(false, nil, stubP)

	target.Start()

	if stubP.startCalled == 0 {
		t.Errorf("TestStandardOutputStart: Failed to call Start on the parent")
	}

}
func TestStandardOutputDealWithIt(t *testing.T) {
	stubP := StubHandler().(*stubHandler)
	target := StandardOutput(false, nil, stubP)

	expectedData := []byte("Expected Data")
	var response http.Response
	response.Body = nopCloser{bytes.NewBuffer(expectedData)}

	testTime := time.Now()
	timer := StubTimer(testTime, testTime.Add(time.Hour))

	target.DealWithIt(response, &timer)

	if stubP.dealCalled == 0 {
		t.Errorf("TestStandardOutputDealWithIt: Failed to call DealWithIt on the parent")
	}
	result, err := ioutil.ReadAll(stubP.savedBody)
	if err != nil {
		t.Errorf("TestStandardOutputDealWithIt: Failed to read expected data from parent call, err %v", err)
	}
	if !bytes.Equal(result, expectedData) {
		t.Errorf("TestStandardOutputDealWithIt: Expected '%v' but got '%v'", expectedData, result)
	}

}

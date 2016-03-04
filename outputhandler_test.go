package main

import (
	"bytes"
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

	var response http.Response
	response.Body = nopCloser{bytes.NewBuffer([]byte("Expected Data"))}

	testTime := time.Now()
	timer := StubTimer(testTime, testTime.Add(time.Hour))

	target.DealWithIt(response, &timer)

	if stubP.dealCalled == 0 {
		t.Errorf("TestStandardOutputDealWithIt: Failed to call DealWithIt on the parent")
	}

}

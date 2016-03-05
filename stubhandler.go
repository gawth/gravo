package main

import (
	"io"
	"net/http"
)

type stubHandler struct {
	dealCalled  int
	savedBody   io.ReadCloser
	logCalled   int
	startCalled int
}

func (this *stubHandler) DealWithIt(r http.Response, t Timer) {
	this.savedBody = r.Body
	this.dealCalled++
}

func (this *stubHandler) LogInfo(s string) {
	this.logCalled++
}

func (this *stubHandler) Start() {
	this.startCalled++
}

func StubHandler() OutputHandler {
	return &stubHandler{}
}

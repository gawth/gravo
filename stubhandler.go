package main

import "net/http"

type stubHandler struct {
	dealCalled int
}

func (this *stubHandler) DealWithIt(r http.Response, t Timer) {
	this.dealCalled++
}

func (this *stubHandler) LogInfo(s string) {

}

func (this *stubHandler) Start() {

}

func StubHandler() OutputHandler {
	return &stubHandler{dealCalled: 0}
}

package main

import "net/http"

type nullHandler struct {
}

func (this *nullHandler) DealWithIt(r http.Response, t Timer) {

}

func (this *nullHandler) LogInfo(s string) {

}

func (this *nullHandler) Start() {

}

func NullHandler() OutputHandler {
	return &nullHandler{}
}

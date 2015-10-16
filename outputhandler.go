package main

import (
	"net/http"
)

type standardOutput struct {
}

func (so *standardOutput) DealWithIt(r http.Response) {
	return
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type standardOutput struct {
	Verbose bool
}

func (so *standardOutput) DealWithIt(r http.Response, t Timer) {
	payload, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(fmt.Sprintf(",%d,%d,%v\n", len(payload), len(payload)/1024/1024, t.GetTime()))
	if so.Verbose {
		fmt.Fprintln(os.Stderr, string(payload))
	}
	return
}

func (so *standardOutput) LogInfo(s string) {
	if so.Verbose {
		fmt.Fprintln(os.Stderr, s)
	}
}

func (so *standardOutput) Start() {
	fmt.Println("timestamp,bytes, meg, duration")
}

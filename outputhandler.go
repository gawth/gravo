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
	V       Validator
}

func (so *standardOutput) DealWithIt(r http.Response, t Timer) {
	payload, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		log.Println(err)
		return
	}

	var isValid string = "unknown"
	if so.V != nil {
		isValid = fmt.Sprintf("%t", so.V.IsValid(payload))
	}
	log.Println(fmt.Sprintf(", %db, %.5fmb, %v, %v", len(payload), float64(len(payload))/1024/1024, t.GetDuration(), isValid))
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
	fmt.Println("timestamp, bytes, meg, duration, valid")
}

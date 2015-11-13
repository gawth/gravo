package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type standardOutput struct {
}

func (so *standardOutput) DealWithIt(r http.Response) {
	//TODO Should return this lot and do it externally
	image, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		log.Println(err)
		return
	}
	//TODO should be timing this as well...again, do it outside of the hit call...
	log.Println(fmt.Sprintf("Got %d bytes, %d meg\n", len(image), len(image)/1024/1024))
	return
}

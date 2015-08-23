package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var callTarget = func(target string) (resp *http.Response, err error) {
	return http.Get(target)
}

func individualCall(c config) {
	t0 := time.Now()
	res, err := callTarget("http://" + c.Target.Host + ":" + c.Target.Port + "/" + c.Target.Path)
	if err != nil {
		log.Println(err)
		return
	}
	t1 := time.Now()

	image, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("Got %d bytes, %d meg in %v\n", len(image), len(image)/1024/1024, t1.Sub(t0))

}

func doStuff(c config) {

	var waitfor = (time.Second / (time.Duration(c.Rate.Rrate) * time.Second)) * time.Second

	fmt.Printf("Attacking for %d requests at a rate of %v\n", c.Requests, waitfor)
	for i := 0; i < c.Requests; i++ {

		fmt.Printf("Loop:%d", i)
		go individualCall(c)

		fmt.Println("Sleeping...")
		time.Sleep(waitfor)

	}
}

func main() {

	c := LoadConfig("gravo.yml")
	fmt.Printf("Config: %v\n", c)

	doStuff(c)

}

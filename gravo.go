package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var callTarget = func(target string) (resp *http.Response, err error) {
	return http.Get(target)
}

var getUrls = func(filename string) (urls []string, err error) {
	return []string{"url stub"}, nil
}

func individualCall(u string, tracker *sync.WaitGroup) {
	defer tracker.Done()

	t0 := time.Now()
	res, err := callTarget(u)
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
	tracker := &sync.WaitGroup{}

	r, err := c.RequestCount()
	if err != nil {
		log.Fatal("error: %v", err)
	}

	if err != nil {
		log.Fatal("error: %v", err)
	}

	fmt.Printf("Attacking for %d requests at a rate of %v\n", r, waitfor)
	for i := 0; i < r; i++ {

		fmt.Printf("Loop:%d", i)

		tracker.Add(1)
		u, err := c.Target.Url(i)
		if err != nil {
			log.Fatal("error: %v", err)
		}
		go individualCall(u, tracker)

		fmt.Println("Sleeping...")
		time.Sleep(waitfor)

	}

	tracker.Wait()
	fmt.Println("Done!!")
}

func main() {

	c := LoadConfig("gravo.yml")
	fmt.Printf("Config: %v\n", c)

	doStuff(c)

}

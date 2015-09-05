package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

func logInfo(c config, s string) {
	if c.Verbose {
		fmt.Printf(s)
	}
}

var callTarget = func(target string) (resp *http.Response, err error) {
	return http.Get(target)
}

var getUrls = func(filename string) ([]string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(content), "\n"), nil
}

func individualCall(u string, c config, tracker *sync.WaitGroup) {
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
	logInfo(c, fmt.Sprintf("Got %d bytes, %d meg in %v\n", len(image), len(image)/1024/1024, t1.Sub(t0)))

}

func validUrl(url string) bool {
	if len(url) == 0 {
		return false
	}

	return true
}

func doStuff(c config) {

	var waitfor = (time.Second / (time.Duration(c.Rate.Rrate) * time.Second)) * time.Second
	tracker := &sync.WaitGroup{}

	r := c.RequestCount()

	logInfo(c, fmt.Sprintf("Attacking for %d requests at a rate of %v\n", r, waitfor))
	for i := 0; i < r; i++ {

		logInfo(c, fmt.Sprintf("Loop:%d\n", i))

		u, err := c.Target.Url(i)
		if err != nil {
			log.Fatal("error: %v", err)
		}

		if validUrl(u) {
			tracker.Add(1)
			go individualCall(u, c, tracker)

			logInfo(c, fmt.Sprintf("Sleeping...\n"))
			time.Sleep(waitfor)
		}
	}

	tracker.Wait()
	logInfo(c, fmt.Sprintf("Done!!"))
}

func main() {

	c := InitialiseConfig("gravo.yml")
	fmt.Printf("Config: %v\n", c)

	doStuff(c)

}

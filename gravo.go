package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"text/template"
	"time"
)

// Target provides an interface for which the hit method will be called
// as part of the load test
type Target interface {
	Hit(tracker *sync.WaitGroup)
}

// Iterator is an interface that is used to iterate over a series of targets
type Iterator interface {
	Next(bool) bool
	Value() Target
}

func logInfo(c config, s string) {
	if c.Verbose {
		fmt.Printf(s)
	}
}
func getTimeUnit(unit string) time.Duration {
	switch unit {
	case "s":
		return time.Second
	case "m":
		return time.Millisecond
	case "M":
		return time.Minute
	case "H":
		return time.Hour
	default:
		return time.Second
	}
}

var callTarget = func(target string, method string, headers http.Header, body string) (resp *http.Response, err error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, target, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header = headers
	return client.Do(req)
}

var getSOAPBody = func(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
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
	res, err := callTarget(u, "GET", map[string][]string{}, "")
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

// runLoad takes config and an iterator.  It uses the iterator to repeatedly
// call the hit method on the value returned by the iterator.  The frequency of
// the calls is based on the Rate defined in the config.
// Go channels and routines are used to ensure the calls to the hit method
// are carried out independently.  The Hit method takes a sync function that
// needs to be called when it has completed
func runLoad(c config, i Iterator) {
	var waitfor = time.Duration(c.Rate.Rrate) * getTimeUnit(c.Rate.Rtype)

	logInfo(c, fmt.Sprintf("Interval %d", waitfor))
	ticker := time.NewTicker(waitfor)
	done := make(chan bool)

	tracker := &sync.WaitGroup{}

	go func(done chan bool) {
		for t := range ticker.C {
			logInfo(c, fmt.Sprintf("Tickt at %s", t))
			if i.Next(false) {
				tracker.Add(1)
				go i.Value().Hit(tracker)
			} else {
				ticker.Stop()
				done <- true
			}

		}
	}(done)

	<-done
	tracker.Wait()
	logInfo(c, fmt.Sprintf("Done!!"))

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

func doSoap(c config) {
	h := http.Header{}
	h.Add("Host", c.Target.Host)
	h.Add("Content-Type", "text/xml; charset=utf-8")

	// Can replace this by a template load from file call
	t, err := getSOAPBody(c.SoapFile)

	tmpl, err := template.New("soap").Parse(t)
	if err != nil {
		log.Fatal("error: %v", err)
	}

	m := make(map[string]string)
	m["ip"] = "9999"

	var body bytes.Buffer
	err = tmpl.Execute(&body, m)
	if err != nil {
		log.Fatal("error: %v", err)
	}

	h.Add("Content-Length", string(len(body.String())))

	url, err := c.Target.Url(0)
	if err != nil {
		log.Fatal("error: %v", err)
	}

	resp, err := callTarget(url, "POST", h, body.String())

	if err != nil {
		log.Println(err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("error: %v", err)
	}

	in := string(b)
	log.Println(in)

}

func main() {

	c := InitialiseConfig("gravo.yml")
	fmt.Printf("Config: %v\n", c)

	if c.Soap {
		doSoap(c)
	} else {
		doStuff(c)
	}

}

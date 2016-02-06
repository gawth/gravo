package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Timer is an interface used to log timing of calls
type Timer interface {
	Start()
	End()
	GetDuration() time.Duration
	GetStart() time.Time
	GetEnd() time.Time
}

// Validator is used to check the output of the call to the target is valid
type Validator interface {
	IsValid([]byte) bool
}

//OutputHandler takes the output from the service request and deals with it
type OutputHandler interface {
	DealWithIt(http.Response, Timer)
	LogInfo(string)
	Start()
}

// Target provides an interface for which the hit method will be called
// as part of the load test
type Target interface {
	Hit(*sync.WaitGroup, Timer, OutputHandler)
}

// Iterator is an interface that is used to iterate over a series of targets
type Iterator interface {
	Next(bool) bool
	Value() (Target, error)
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

// runLoad takes config and an iterator.  It uses the iterator to repeatedly
// call the hit method on the value returned by the iterator.  The frequency of
// the calls is based on the Rate defined in the config.
// Go channels and routines are used to ensure the calls to the hit method
// are carried out independently.  The Hit method takes a sync function that
// needs to be called when it has completed
func runLoad(c config, i Iterator, ti Timer, o OutputHandler) {
	var waitfor = time.Duration(c.Rate.Rrate) * getTimeUnit(c.Rate.Rtype)

	o.LogInfo(fmt.Sprintf("Interval %d\n", waitfor))
	ticker := time.NewTicker(waitfor)
	done := make(chan bool)

	tracker := &sync.WaitGroup{}

	o.Start()
	go func(done chan bool) {
		for t := range ticker.C {
			o.LogInfo(fmt.Sprintf("Tickt at %s\n", t))
			if i.Next(false) {
				v, err := i.Value()
				if err != nil {
					o.LogInfo(fmt.Sprintf("Error from iterator value so skipping this value :%s\n", err))
					continue
				}
				tracker.Add(1)
				go v.Hit(tracker, ti, o)
			} else {
				ticker.Stop()
				done <- true
			}

		}
	}(done)

	<-done
	tracker.Wait()
	o.LogInfo(fmt.Sprintf("Done!!\n"))

}

func main() {
	c := initialiseConfig("config.yml")
	var resultsFile string
	flag.StringVar(&resultsFile, "file", "", "results file from previous run")
	flag.Parse()

	var validator Validator
	var output standardOutput
	if len(c.Regex) > 0 {
		validator = &regexValidator{c.validator}
		output = standardOutput{c.Verbose, validator}
	} else {
		output = standardOutput{Verbose: c.Verbose}
	}

	if len(resultsFile) > 0 {
		results := chartHandler{filename: resultsFile, completed: make(chan bool)}
		results.Start()
		<-results.completed
	} else if len(c.DataFile) > 0 {
		iterator := dataIterator{url: c.Target.urls[0], columns: c.columns, data: c.data, template: c.template, verb: c.Verb, headers: c.Headers}
		runLoad(c, &iterator, &timer{}, &output)
	} else {
		iterator := urlIterator{urls: c.Target.urls, verb: c.Verb}
		runLoad(c, &iterator, &timer{}, &output)
	}

}

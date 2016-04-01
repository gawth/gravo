package main

import (
	"log"
	"net/http"
	"time"
)

type metric struct {
	Datetime time.Time
	Val      int64
}

type chartHandler struct {
	completed chan bool
	data      *Results
	logger    chan metric
	parent    OutputHandler
}

func (ch *chartHandler) DealWithIt(r http.Response, t Timer) {
	ch.logger <- metric{t.GetStart(), t.GetDuration().Nanoseconds()}

	ch.parent.DealWithIt(r, t)

	return
}
func (ch *chartHandler) updateData() {
	for {
		d := <-ch.logger
		ch.data.Save(d)
	}
}

func (ch *chartHandler) LogInfo(s string) {
	ch.parent.LogInfo(s)
}

func (ch *chartHandler) Start() {
	if len(ch.data.Name()) == 0 {
		log.Fatal("Must specify a filename for results")
	}
	go ch.updateData()
	ch.parent.Start()
}

func ChartHandler(channel chan bool, results *Results, parent OutputHandler) OutputHandler {
	var val *chartHandler
	if parent == nil {
		parent = NullHandler()
	}
	val = &chartHandler{data: results, completed: channel, parent: parent}
	val.logger = make(chan metric)
	return val
}

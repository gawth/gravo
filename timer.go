package main

import (
	"time"
)

type timer struct {
	started  time.Time
	finished time.Time
}

func (t *timer) Start() {
	t.started = time.Now()
	return
}

func (t *timer) End() {
	t.finished = time.Now()
	return
}

func (t *timer) GetTime() time.Duration {
	return t.finished.Sub(t.started)
}

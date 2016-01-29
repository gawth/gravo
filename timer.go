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

func (t *timer) GetDuration() time.Duration {
	return t.finished.Sub(t.started)
}

func (t *timer) GetStart() time.Time {
	return t.started
}

func (t *timer) GetEnd() time.Time {
	return t.finished
}

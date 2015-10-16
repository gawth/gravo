package main

import (
	"time"
)

type timer struct {
}

func (t *timer) Start() {
	return
}

func (t *timer) End() {
	return
}

func (t *timer) GetTime() time.Duration {
	return time.Second
}

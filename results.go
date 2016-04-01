package main

import ()

type Results struct {
	name    string
	metrics []metric
}

func (res *Results) Name() string {
	return res.name
}
func (res *Results) Save(data metric) {
	res.metrics = append(res.metrics, data)
}

func (res *Results) GetMetrics() []metric {
	return res.metrics
}

func NewResults(name string) Results {
	return Results{name: name}
}

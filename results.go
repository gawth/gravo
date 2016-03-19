package main

import ()

type Results struct {
	metrics []metric
}

func (res *Results) Save(data metric) {
	res.metrics = append(res.metrics, data)
}

func (res *Results) GetMetrics() []metric {
	return res.metrics
}

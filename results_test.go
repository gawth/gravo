package main

import (
	"testing"
	"time"
)

func TestResultsSave(t *testing.T) {
	target := Results{}
	var data metric

	for i := int64(0); i < 5; i++ {
		data = metric{time.Now(), i}
		target.Save(data)

		if int64(len(target.GetMetrics())) != i+1 {
			t.Errorf("TestResultsSaveAndGet: Expected GetMetrics to return %v not %v", i+1, len(target.GetMetrics()))
		}

		if target.GetMetrics()[i].Val != i {
			t.Errorf("TestResultsSaveAndGet: Expected GetMetrics last value to have a val of %v not %v", i, target.GetMetrics()[i].Val)
		}
	}
}

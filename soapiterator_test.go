package main

import (
	"reflect"
	"testing"
)

func TestSoapIterator(t *testing.T) {
	var col1 = "c1"
	var col2 = "c2"
	var col3 = "c3"
	var testcols = []string{col1, col2, col3}
	var testurl = "test url"

	var td = []map[string]string{{
		col1: "data1",
		col2: "data2",
		col3: "data3",
	}}

	var it = soapIterator{url: testurl, columns: testcols, data: td}
	var called = 0

	expected := reflect.TypeOf(&urlTarget{})

	for it.Next(false) {
		typ := reflect.TypeOf(it.Value())
		if typ != expected {
			t.Errorf("TestSoapIterator: Expected %v but got %v", expected, typ)
		}
		called++

		// the vall to reflect.ValueOf returns a Value type from the reflect package
		// In this case that value type contains an interface so we get at the Interface using the
		// Interface function.  We can then cast that to urlTarget
		//
		var concreteTarget = reflect.ValueOf(it.Value()).Interface().(*urlTarget)
		if concreteTarget.url != testurl {
			t.Errorf("TestSoapterator: Expected url %v but got %v", testurl, concreteTarget.url)
		}
	}
	if called != len(td) {
		t.Errorf("TestSoapterator: Expected next %v times but got %v", len(td), called)
	}
}

package main

import (
	"reflect"
	"testing"
)

func TestSoapIterator(t *testing.T) {
	var col1 = "c1"
	var col2 = "c2"
	var col3 = "c3"
	var testdata = make(map[string]string)
	testdata[col1] = "data1"
	testdata[col2] = "data2"
	testdata[col3] = "data3"
	var testcols = []string{col1, col2, col3}

	var it = soapIterator{url: "test url", columns: testcols, data: testdata}
	var called = 0

	expected := reflect.TypeOf(&urlTarget{})

	for it.Next(false) {
		typ := reflect.TypeOf(it.Value())
		if typ != expected {
			t.Errorf("TestSoapIterator: Expected %v but got %v", expected, typ)
		}
		called++
	}
	if called != 1 {
		t.Errorf("TestSoapterator: Expected next %v times but got %v", 1, called)
	}
}

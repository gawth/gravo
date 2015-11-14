package main

import (
	"reflect"
	"testing"
	"text/template"
)

func TestSoapIterator(t *testing.T) {
	var col1 = "c1"
	var col2 = "c2"
	var col3 = "c3"
	var testcols = []string{col1, col2, col3}
	var testurl = "test url"

	var td = [][]string{{
		"data1",
		"data2",
		"data3",
	}}

	var it = soapIterator{url: testurl, columns: testcols, data: td}
	var called = 0

	expected := reflect.TypeOf(&urlTarget{})

	for it.Next(false) {
		val, _ := it.Value()
		typ := reflect.TypeOf(val)
		if typ != expected {
			t.Errorf("TestSoapIterator: Expected %v but got %v", expected, typ)
		}
		called++

		// the vall to reflect.ValueOf returns a Value type from the reflect package
		// In this case that value type contains an interface so we get at the Interface using the
		// Interface function.  We can then cast that to urlTarget
		//
		var concreteTarget = reflect.ValueOf(val).Interface().(*urlTarget)
		if concreteTarget.url != testurl {
			t.Errorf("TestSoapterator: Expected url %v but got %v", testurl, concreteTarget.url)
		}
	}
	if called != len(td) {
		t.Errorf("TestSoapterator: Expected next %v times but got %v", len(td), called)
	}
}

func TestSoapTemplate(t *testing.T) {
	var col1 = "c1"
	var col2 = "c2"
	var col3 = "c3"
	var testcols = []string{col1, col2, col3}
	var testurl = "test url"

	var td = [][]string{
		{
			"data1",
			"data2",
			"data3",
		},
		{
			"data1",
			"data2",
			"data3",
		},
	}

	var tmpText = `Hello! {{.c1}} and {{.c2}} and finally {{.c3}}`
	var expectedBody = "Hello! data1 and data2 and finally data3"

	tmpl, err := template.New("test").Parse(tmpText)
	if err != nil {
		t.Errorf("TestSoapTemplate: Failed to parse the test template:%v", err)
	}

	var it = soapIterator{url: testurl, columns: testcols, data: td, template: tmpl}
	var called = 0

	for it.Next(false) {
		called++

		val, _ := it.Value()

		var concreteTarget = reflect.ValueOf(val).Interface().(*urlTarget)
		if concreteTarget.body != expectedBody {
			t.Errorf("TestSoapTemplate: Incorrect body '%v'", concreteTarget.body)
		}
	}
	if called != len(td) {
		t.Errorf("TestSoapTemplate: Expected next %v times but got %v", len(td), called)
	}
}

func TestErrorCreatingSoapTemplate(t *testing.T) {
	var col1 = "c1"
	var col2 = "c2"
	var col3 = "c3"
	var testcols = []string{col1, col2, col3}
	var testurl = "test url"

	var dataError = [][]string{
		{
			"data1",
			"data2",
		},
	}

	tmpl, err := template.New("test").Parse("")
	if err != nil {
		t.Errorf("TestErrorCreatingSoapTemplate: Failed to parse the test template:%v", err)
	}

	var it = soapIterator{url: testurl, columns: testcols, data: dataError, template: tmpl}
	for it.Next(false) {
		_, err := it.Value()
		if err == nil {
			t.Errorf("TestErrorCreatingSoapTemplate: Expected an error but didn't get one")
		}

	}
}

package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type dataIterator struct {
	url      string
	columns  []string
	data     [][]string
	template *template.Template
	position int
	verb     string
	headers  http.Header
}

func (it *dataIterator) Next(continuous bool) bool {
	if it.position >= len(it.data) {
		return false
	}
	it.position++
	return true
}
func (it *dataIterator) Value() (Target, error) {
	var body bytes.Buffer

	if len(it.data[it.position-1]) != len(it.columns) {
		return &urlTarget{}, fmt.Errorf("dataIterator: Incorrect number of data items line %v.  Expected %v but got %v", it.position-1, len(it.columns), len(it.data[it.position-1]))
	}
	if it.template != nil {

		var tmpMap = make(map[string]string)

		// Use the template plus data and columns to get the body data
		for i := 0; i < len(it.columns); i++ {
			tmpMap[it.columns[i]] = it.data[it.position-1][i]
		}

		err := it.template.Execute(&body, tmpMap)
		if err != nil {
			log.Fatal(err)
		}
	}

	//h := http.Header{}
	//h.Add("Content-Type", "text/xml; charset=utf-8")
	//h.Add("Content-Type", "application/x-www-form-urlencoded")
	retVal := urlTarget{method: it.verb, url: it.url, headers: it.headers, body: body.String()}
	return &retVal, nil
}

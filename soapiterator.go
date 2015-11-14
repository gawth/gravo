package main

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
)

type soapIterator struct {
	url      string
	columns  []string
	data     [][]string
	template *template.Template
	position int
}

func (it *soapIterator) Next(continuous bool) bool {
	if it.position >= len(it.data) {
		return false
	}
	it.position++
	return true
}
func (it *soapIterator) Value() (Target, error) {
	var body bytes.Buffer

	if len(it.data[it.position-1]) != len(it.columns) {
		return &urlTarget{}, fmt.Errorf("soapIterator: Incorrect number of data items line %v.  Expected %v but got %v", it.position-1, len(it.columns), len(it.data[it.position-1]))
	}
	if it.template != nil {

		var tmpMap = make(map[string]string)

		// Use the template plus data and columns to get the body data
		for i := 0; i < len(it.columns); i++ {
			tmpMap[it.columns[i]] = it.data[it.position-1][i]
		}

		err := it.template.Execute(&body, tmpMap)
		if err != nil {
			log.Fatal("error: %v", err)
		}
	}

	retVal := urlTarget{method: "POST", url: it.url, body: body.String()}
	return &retVal, nil
}

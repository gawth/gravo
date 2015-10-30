package main

import ()

type soapIterator struct {
	url      string
	columns  []string
	data     map[string]string
	position int
}

func (it *soapIterator) Next(continuous bool) bool {
	if it.position >= len(it.data) {
		return false
	}
	it.position++
	return true
}
func (it *soapIterator) Value() Target {
	// Why position-1?  Because the initial value will be 0, the first call to
	// next will increment to 1 but we wont have got the value from position 0
	// at that point.
	retVal := urlTarget{method: "POST", url: it.url, body: ""}
	return &retVal
}

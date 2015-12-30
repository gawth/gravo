package main

import (
	"regexp"
	"testing"
)

func TestRegexValiator(t *testing.T) {
	reggie := regexp.MustCompile(`.*`)
	target := regexValidator{reg: reggie}

	if !target.IsValid([]byte("abc")) {
		t.Errorf("TestRegexValidator: Failed to validate wildcard regexp")
	}

}
func TestRegexInvalid(t *testing.T) {
	reggie := regexp.MustCompile(`aaa`)
	target := regexValidator{reg: reggie}

	if target.IsValid([]byte("abc")) {
		t.Errorf("TestRegexValidator: Incorrectly validated 'aaa'")
	}

}

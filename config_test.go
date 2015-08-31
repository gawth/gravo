package main

import (
	"testing"
)

func TestConfigTargetUrlSimpleUrl(t *testing.T) {
	targ := target{Host: "testhost", Port: "1234", Path: "apath"}
	var expected = "http://testhost:1234/apath"

	actual, _ := targ.Url(123)

	if actual != expected {
		t.Errorf("Got '%v' instead of '%v'", actual, expected)
	}

}

func TestConfigTargetUrlFromUrls(t *testing.T) {
	targ := target{File: "afile", Host: "testhost", Port: "1234", Path: "apath"}
	var expected = "http://testhost:1234/apath"

	targ.urls = []string{expected}

	actual, _ := targ.Url(0)

	if actual != expected {
		t.Errorf("Got '%v' instead of '%v'", actual, expected)
	}
}
func TestConfigTargetUrlFromUrlsOutOfBoundsCheck(t *testing.T) {
	targ := target{File: "afile", Host: "testhost", Port: "1234", Path: "apath"}

	targ.urls = []string{"irrelevant value"}

	_, err := targ.Url(1)

	if err == nil {
		t.Errorf("Url call should have generated an error")
	}
}

func TestConstructUrl(t *testing.T) {
	targ := target{Host: "testhost", Port: "1234", Path: "apath"}
	var expected = "http://testhost:1234/apath"

	actual := targ.ConstructUrl()
	if actual != expected {
		t.Errorf("Got '%v' instead of '%v'", actual, expected)
	}
}

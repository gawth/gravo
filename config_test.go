package main

import (
	"reflect"
	"testing"
)

func TestConstructUrl(t *testing.T) {
	targ := target{Host: "testhost", Port: "1234", Path: "apath"}
	var expected = "http://testhost:1234/apath"

	actual := targ.ConstructURL()
	if actual != expected {
		t.Errorf("Got '%v' instead of '%v'", actual, expected)
	}
}

func TestDeleteBlanksFunction(t *testing.T) {
	targetStrings := []string{"line one", "", "", "line two", "line three", ""}

	actualStrings := deleteBlanks(targetStrings)

	if len(actualStrings) != 3 {
		t.Errorf("deleteBlanks: Expected '%v' lines, got '%v'", 3, actualStrings)
	}
}

func TestLoadUrls(t *testing.T) {
	targ := target{Host: "testhost", Port: "1234", Path: "apath"}
	expectedUrls := []string{"http://testhost:1234/apath"}

	targ.loadUrls()

	if !reflect.DeepEqual(targ.urls, expectedUrls) {
		t.Errorf("loadUrls: %v does not match expected %v", targ.urls, expectedUrls)
	}

}

func TestLoadUrlsFromFileError(t *testing.T) {
	targ := target{Host: "testhost", Port: "1234", Path: "apath", File: "no file"}

	targ.loadUrls()

	if len(targ.urls) != 0 {
		t.Errorf("loadUrls from file: Somehow managed to load some urls %v", targ.urls)
	}

}

func TestCovertYaml(t *testing.T) {
	data := []byte("target:\n    host: test")

	c := convertYaml(data)

	if c.Target.Host != "test" {
		t.Errorf("convertYaml: Expected config to have host %v but was %v", "test", c)
	}
}

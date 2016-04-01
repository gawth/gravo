package main

import (
	"os"
	"testing"
)

func TestResultsListGetList(t *testing.T) {
	target := NewResultsList(NewResults("fred"), "")

	list := target.GetList()

	if len(list) != 1 {
		t.Errorf("TestResultsList: Expected GetList to return %v not %v", 1, len(list))
	}
}

var folder string = "testfiles"
var files []string = []string{"file1", "file2", "file3"}

func setUpTestFiles() {
	os.Mkdir(folder, 0777)
	for _, name := range files {
		os.Create(folder + "/" + name)
	}
}
func addFolder() {
	os.Mkdir(folder+"/another", 0777)
}
func destroyTestFiles() {
	os.RemoveAll(folder)
}

func TestResultsListReadFile(t *testing.T) {
	setUpTestFiles()
	defer destroyTestFiles()

	target := NewResultsList(Results{}, folder)

	list := target.GetList()

	if len(list) != len(files)+1 {
		t.Errorf("TestResultsListReadFile: Expected GetList to return %v not %v",
			len(files)+1, len(list))
	}
}

func TestResultsListReadFileWithFolder(t *testing.T) {
	setUpTestFiles()
	addFolder()
	defer destroyTestFiles()

	target := NewResultsList(Results{}, folder)

	list := target.GetList()

	if len(list) != len(files)+1 {
		t.Errorf("TestResultsListReadFileWithFolder: Expected GetList to return %v not %v",
			len(files)+1, len(list))
	}
}

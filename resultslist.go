package main

import (
	"io/ioutil"
	"log"
)

type ResultsList struct {
	liveResults Results
	folder      string
	filelist    []string
}

func (rl *ResultsList) GetList() []string {
	tmp := []string{rl.liveResults.Name()}
	tmp = append(tmp, rl.filelist...)
	return tmp
}

func (rl *ResultsList) GetResults(name string) string {
	return ""
}
func (rl *ResultsList) InitResultsList() {
	if len(rl.folder) == 0 {
		return
	}
	files, err := ioutil.ReadDir(rl.folder)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if !file.IsDir() {
			rl.filelist = append(rl.filelist, file.Name())
		}
	}
}

func NewResultsList(liveResults Results, folder string) ResultsList {
	res := ResultsList{liveResults: liveResults, folder: folder}
	res.InitResultsList()
	return res
}

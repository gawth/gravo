package main

import ()

type ResultsList struct {
	liveResults Results
}

func (rl *ResultsList) GetList() []string {
	return []string{rl.liveResults.Name()}
}

func NewResultsList(liveResults Results) ResultsList {
	return ResultsList{liveResults: liveResults}
}

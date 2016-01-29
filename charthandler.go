package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type chartHandler struct {
	data   []string
	logger chan string
}

func (ch *chartHandler) DealWithIt(r http.Response, t Timer) {
	payload, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		log.Println(err)
		return
	}

	ch.logger <- fmt.Sprintf("%v, %v, %v", t.GetStart(), t.GetDuration(), len(payload))

	return
}

func (ch *chartHandler) LogInfo(s string) {
}

func (ch *chartHandler) statsHandler(w http.ResponseWriter, r *http.Request) {
	j, err := json.Marshal(ch.data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func (ch *chartHandler) updateData() {
	for {
		d := <-ch.logger
		ch.data = append(ch.data, d)
	}
}

func (ch *chartHandler) Start() {
	ch.logger = make(chan string)

	go ch.updateData()

	r := mux.NewRouter()
	r.HandleFunc("/stats", ch.statsHandler).Methods("GET")
	http.Handle("/", r)
	fmt.Println("Listening on port 8080")
	go http.ListenAndServe(":8080", nil)
}

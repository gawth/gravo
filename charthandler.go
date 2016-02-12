package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type metric struct {
	Datetime time.Time
	Val      int64
}

type chartHandler struct {
	filename  string
	completed chan bool
	data      []metric
	logger    chan metric
}

func (ch *chartHandler) DealWithIt(r http.Response, t Timer) {
	_, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		log.Println(err)
		return
	}

	ch.logger <- metric{t.GetStart(), t.GetDuration().Nanoseconds()}

	return
}

func (ch *chartHandler) LogInfo(s string) {
}

func (ch *chartHandler) statsHandler(w http.ResponseWriter, r *http.Request) {
	j, err := json.MarshalIndent(ch.data, "", "   ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("results.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, nil)
}

func (ch *chartHandler) updateData() {
	for {
		d := <-ch.logger
		ch.data = append(ch.data, d)
	}
}

func (ch *chartHandler) loadData() {
	f, err := os.Open(ch.filename)
	defer f.Close()
	dec := json.NewDecoder(f)
	// read open bracket
	t, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}

	// while the array contains values
	for dec.More() {

		var m metric
		// decode an array value (Message)
		err := dec.Decode(&m)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%v: %v\n", m.Datetime, m.Val)
		ch.data = append(ch.data, m)
	}

	// read closing bracket
	t, err = dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%T: %v\n", t, t)
}

func (ch *chartHandler) Start() {
	ch.logger = make(chan metric)

	if len(ch.filename) > 0 {
		ch.loadData()
	}

	go ch.updateData()

	r := mux.NewRouter()
	r.HandleFunc("/stats", ch.statsHandler).Methods("GET")
	r.HandleFunc("/results", resultsHandler).Methods("GET")
	http.Handle("/", r)
	fmt.Println("Listening on port 8080")
	go http.ListenAndServe(":8080", nil)
}

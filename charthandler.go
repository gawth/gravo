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

	"github.com/gorilla/handlers"
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
	parent    OutputHandler
}

func (ch *chartHandler) DealWithIt(r http.Response, t Timer) {
	_, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		log.Println(err)
		return
	}

	ch.logger <- metric{t.GetStart(), t.GetDuration().Nanoseconds()}

	ch.parent.DealWithIt(r, t)

	return
}

func (ch *chartHandler) LogInfo(s string) {
	ch.parent.LogInfo(s)
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
		fmt.Printf("%v\n", d)
	}
}

func (ch *chartHandler) loadData() {
	f, err := os.Open(ch.filename)
	defer f.Close()

	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Fatal(err)
	}
	dec := json.NewDecoder(f)
	// read open bracket
	_, err = dec.Token()
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

		ch.data = append(ch.data, m)
	}

	// read closing bracket
	_, err = dec.Token()
	if err != nil {
		log.Fatal(err)
	}
}

func (ch *chartHandler) Start() {

	if len(ch.filename) == 0 {
		log.Fatal("Must specify a filename for results")
	}
	ch.loadData()

	go ch.updateData()

	r := mux.NewRouter()
	r.HandleFunc("/stats", ch.statsHandler).Methods("GET")
	r.HandleFunc("/results", resultsHandler).Methods("GET")
	r.HandleFunc("/results/"+ch.filename, resultsHandler).Methods("GET")

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	fmt.Println("Listening on port http://localhost:8910/results/" + ch.filename)
	go http.ListenAndServe(":8910", loggedRouter)

	ch.parent.Start()
}

func ChartHandler(resultsFile string, channel chan bool, parent OutputHandler) OutputHandler {
	var val *chartHandler
	if parent == nil {
		parent = NullHandler()
	}
	val = &chartHandler{filename: resultsFile, completed: channel, parent: parent}
	val.logger = make(chan metric)
	return val
}

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
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
	ch.parseData(bufio.NewReader(f))
}
func (ch *chartHandler) parseData(stream io.Reader) {

	scanner := bufio.NewScanner(stream)

	for scanner.Scan() {

		var m metric

		// Ignore lines that don't convert to metric
		if err := json.Unmarshal(scanner.Bytes(), &m); err == nil {
			ch.data = append(ch.data, m)
		}
		if len(ch.data) == 0 {
			log.Fatalln("Attempted to process file but contains no valid metrics")
		}
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

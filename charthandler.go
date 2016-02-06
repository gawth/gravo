package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type chartHandler struct {
	filename  string
	completed chan bool
	data      []string
	logger    chan string
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
	d, err := ioutil.ReadFile(ch.filename)
	if err != nil {
		log.Panic(err)
	}
	ch.data = deleteBlanks(strings.Split(string(d), "\n"))
}

func (ch *chartHandler) Start() {
	ch.logger = make(chan string)

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

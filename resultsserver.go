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

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type FileView struct {
	Filename string
	FileList []string
}

type resultsServer struct {
	data         Results
	savedResults ResultsList
}

func (ch *resultsServer) statsHandler(w http.ResponseWriter, r *http.Request) {
	j, err := json.MarshalIndent(ch.data.metrics, "", "   ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
func (ch *resultsServer) resultsList(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatal(err)
	}
	view := FileView{Filename: ch.data.Name(), FileList: ch.savedResults.GetList()}
	t.Execute(w, view)
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("results.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, nil)
}

func (ch *resultsServer) loadData() {
	f, err := os.Open(ch.data.Name())
	defer f.Close()

	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Fatal(err)
	}
	ch.parseData(bufio.NewReader(f))
}
func (ch *resultsServer) parseData(stream io.Reader) {

	scanner := bufio.NewScanner(stream)

	for scanner.Scan() {

		var m metric

		// Ignore lines that don't convert to metric
		if err := json.Unmarshal(scanner.Bytes(), &m); err == nil {
			ch.data.Save(m)
		}
		if len(ch.data.metrics) == 0 {
			log.Fatalln("Attempted to process file but contains no valid metrics")
		}
	}
}

func (ch *resultsServer) Start() {

	if len(ch.data.Name()) == 0 {
		log.Fatal("Must specify a filename for results")
	}
	ch.loadData()

	r := mux.NewRouter()
	r.HandleFunc("/stats", ch.statsHandler).Methods("GET")
	r.HandleFunc("/results", resultsHandler).Methods("GET")
	r.HandleFunc("/results/"+ch.data.Name(), resultsHandler).Methods("GET")
	r.HandleFunc("/", ch.resultsList).Methods("GET")

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	fmt.Println("Listening on port http://localhost:8910/results/" + ch.data.Name())
	go http.ListenAndServe(":8910", loggedRouter)

}

func ResultsServer(resultsFile string) resultsServer {
	val := resultsServer{data: Results{name: resultsFile, metrics: []metric{}}}
	return val
}

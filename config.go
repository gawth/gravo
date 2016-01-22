package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"
)

var configFile = "gravo.yml"

// deleteBlanks is a simple function to remove blank lines from an array of strings
// Used to tidy up URL file input
func deleteBlanks(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// getUrls reads URLS from a file, removes blank lines and then returns.
var getUrls = func(filename string) ([]string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return deleteBlanks(strings.Split(string(content), "\n")), nil
}

type target struct {
	Host string
	Port string
	Path string
	File string
	urls []string
}

func (t *target) ConstructURL() string {
	return "http://" + t.Host + ":" + t.Port + "/" + t.Path
}

// loadUrls initialises the URLs on the config either with URLs from a file
// if available OR by using the config data provided.
func (t *target) loadUrls() {

	// If we're not using a file then just construct the URL
	if t.File == "" {
		t.urls = []string{t.ConstructURL()}
		return

	}
	// If we've not get any URLs try and get from fie
	if t.urls == nil {
		t.urls, _ = getUrls(t.File)
	}
}

type runrate struct {
	Rrate int
	Rtype string
}
type config struct {
	Target       target
	Requests     int
	Rate         runrate
	Verbose      bool
	TemplateFile string
	DataFile     string
	Verb         string
	Regex        string
	columns      []string
	data         [][]string
	template     *template.Template
	validator    *regexp.Regexp
	Headers      http.Header
}

func readConfigFile(file string) []byte {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func convertYaml(raw []byte) config {
	c := config{}

	err := yaml.Unmarshal(raw, &c)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func loadTemplate(filename string) (*template.Template, error) {
	t, err := template.ParseFiles(filename)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func initialiseConfig(file string) config {
	c := convertYaml(readConfigFile(file))
	c.Target.loadUrls()

	if len(c.Regex) > 0 {
		// Compile and save the regex from the config
		c.validator = regexp.MustCompile(c.Regex)
	}

	if len(c.DataFile) > 0 {
		tmp, err := loadTemplate(c.TemplateFile)
		if err != nil {
			log.Fatal(fmt.Sprintf("Config: Unable to load template %v err: %v", file, err))
		}
		c.template = tmp

		//TODO Need to clean this up
		f, err := os.Open(c.DataFile)
		if err != nil {
			log.Fatal(fmt.Sprintf("Config: Unable to open data file (%v): %v", c.DataFile, err))
		}
		reader := csv.NewReader(bufio.NewReader(f))

		raw, err := reader.ReadAll()
		if err != nil {
			log.Fatal(fmt.Sprintf("Config: Unable to read data file (%v): %v", c.DataFile, err))
		}
		if len(raw) < 2 {
			log.Fatal(fmt.Sprintf("Config: Too few lines (%v) in data file", len(raw)))
		}
		// First row "should" contain headers (how to check?)
		c.columns = raw[0]
		// The rest of the file is the data
		c.data = raw[1:len(raw)]
	}
	return c
}

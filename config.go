package gravo

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
	"text/template"
)

var configFile = "gravo.yml"

func deleteBlanks(s []string) []string {
	//TODO Need a test for this!
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

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

func (t *target) LoadUrls() {

	// If we're not using a file then just construct the URL
	if t.File == "" {
		t.urls = []string{t.ConstructURL()}

	}
	// If we've not get any URLs try and get from fie
	if t.urls == nil {
		t.urls, _ = getUrls(t.File)
	}
}
func (t *target) URL(index int) (string, error) {
	if len(t.urls) == 0 {
		return t.ConstructURL(), nil
	}
	if index >= len(t.urls) {
		return "", errors.New(fmt.Sprintf("Attempted to get URL at %d from URLs length %d", index, len(t.urls)))
	}
	return t.urls[index], nil
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
	Soap         bool
	SoapFile     string
	soapTemplate *template.Template
}

func (c *config) RequestCount() int {
	// If requests has been set to > 0 then assume we're dealing with a single repeated
	// request to a url for now.  Later on we might get a bit more complicated and use
	// combo of requests and file URLs
	//
	if c.Requests != 0 {
		return c.Requests
	}
	return len(c.Target.urls)
}

func readConfigFile(file string) []byte {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("error: %v", err)
	}
	return f
}

func convertYaml(raw []byte) config {
	c := config{}

	err := yaml.Unmarshal(raw, &c)
	if err != nil {
		log.Fatal("error: %v", err)
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

func InitialiseConfig(file string) config {
	c := convertYaml(readConfigFile(file))
	c.Target.LoadUrls()

	if c.Soap {
		tmp, err := loadTemplate(c.SoapFile)
		if err != nil {
			log.Fatal("Config: Unable to load SOAP template %v err: %v", file, err)
		}
		c.soapTemplate = tmp
	}
	return c
}

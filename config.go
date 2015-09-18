package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var configFile = "gravo.yml"

type target struct {
	Host string
	Port string
	Path string
	File string
	urls []string
}

func (t *target) ConstructUrl() string {
	return "http://" + t.Host + ":" + t.Port + "/" + t.Path
}

func (t *target) LoadUrls() {

	// If we're not using a file then just construct the URL
	if t.File == "" {
		t.urls = []string{t.ConstructUrl()}

	}
	// If we've not get any URLs try and get from fie
	if t.urls == nil {
		t.urls, _ = getUrls(t.File)
	}
}
func (t *target) Url(index int) (string, error) {
	if len(t.urls) == 0 {
		return t.ConstructUrl(), nil
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
	Target   target
	Requests int
	Rate     runrate
	Verbose  bool
	Soap     bool
	SoapFile string
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
func InitialiseConfig(file string) config {
	c := convertYaml(readConfigFile(file))
	c.Target.LoadUrls()
	return c
}

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

func (t *target) Urls() ([]string, error) {

	// If we're not using a file then just construct the URL
	if t.File == "" {
		return []string{"http://" + t.Host + ":" + t.Port + "/" + t.Path}, nil

	}
	var err error
	// If we've not get any URLs try and get from fie
	if t.urls == nil {
		t.urls, err = getUrls(t.File)
	}
	return t.urls, err
}
func (t *target) Url(index int) (string, error) {
	if t.File == "" {
		return "http://" + t.Host + ":" + t.Port + "/" + t.Path, nil
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
}

func (c *config) RequestCount() (int, error) {
	if c.Requests != 0 {
		return c.Requests, nil
	}
	urls, err := c.Target.Urls()
	if err != nil {
		return 0, err
	}
	return len(urls), nil
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
func LoadConfig(file string) config {
	return convertYaml(readConfigFile(file))
}

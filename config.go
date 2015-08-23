package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var configFile = "gravo.yml"

type target struct {
	Host string
	Port string
	Path string
}
type config struct {
	Target   target
	Requests int
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

package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var configFile = "gravo.yml"

type config struct {
	Host string
	Port string
	Path string
}

func main() {

	// Read config
	raw, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("error: %v", err)
	}

	c := config{}

	err = yaml.Unmarshal([]byte(raw), &c)
	if err != nil {
		log.Fatal("error: %v", err)
	}

	fmt.Printf("Config: %v\n", c)

	for i := 0; i < 10; i++ {

		t0 := time.Now()
		res, err := http.Get("http://" + c.Host + ":" + c.Port + "/" + c.Path)
		if err != nil {
			log.Fatal(err)
		}
		t1 := time.Now()

		image, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d: Got %d bytes, %d meg in %v\n", i, len(image), len(image)/1024/1024, t1.Sub(t0))

	}
}

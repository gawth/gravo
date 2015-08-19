package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {

	c := LoadConfig("gravo.yml")

	fmt.Printf("Config: %v\n", c)

	for i := 0; i < 10; i++ {

		t0 := time.Now()
		res, err := http.Get("http://" + c.Target.Host + ":" + c.Target.Port + "/" + c.Target.Path)
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

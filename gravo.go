package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var callTarget = func(target string) (resp *http.Response, err error) {
	return http.Get(target)
}

func doStuff(c config) {

	fmt.Printf("Attacking for %d requests\n", c.Requests)
	for i := 0; i < c.Requests; i++ {

		t0 := time.Now()
		res, err := callTarget("http://" + c.Target.Host + ":" + c.Target.Port + "/" + c.Target.Path)
		if err != nil {
			log.Println(err)
			continue
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

func main() {

	c := LoadConfig("gravo.yml")
	fmt.Printf("Config: %v\n", c)

	doStuff(c)

}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
    "time"
)

func main() {
    for i := 0; i < 10; i++ {

        t0 := time.Now()
        res, err := http.Get("http://localhost:8080/img")
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

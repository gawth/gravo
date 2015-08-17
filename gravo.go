package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
    "time"
    "gopkg.in/yaml.v2"
)

var data = `
host: 46.101.29.33
port: 8080
path: img
`

type config struct {
    Host string
    Port string
    Path string
}

func main() {
    c := config{}

    err := yaml.Unmarshal([]byte(data), &c)
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

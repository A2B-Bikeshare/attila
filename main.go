package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	//"io/ioutil"
	"runtime"
	"time"
	"net/url"
	"encoding/json"
	//"strings"
)

var patternsGlobal chan []reqBody
var times chan struct{}

// The way request should have been implemented
type reqBody struct {
	Body string
	Req *http.Request
}

//structs for config
type Config struct {
	Workers int `json:"workers"`
	Patterns []JsonPattern `json:"patterns"`
}

type JsonPattern struct {
	Title string `json:"title"`
	Reqs []JsonReq `json:"requests"`
}

type JsonReq struct {
	URL string `json:"url"`
	Method string `json:"method"`
	Data map[string]string `json:"data"`
	CType string `json:"content-type"`
}

func init() {
	patternsGlobal = make(chan []reqBody)
	times = make(chan struct{}, 1)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	infile, err := os.Open("config.json")
	defer infile.Close()
	if err != nil { panic(err) }
	
	dec := json.NewDecoder(infile)
	
	config := new(Config)
	err = dec.Decode(config)
	if err != nil { panic(err) }

	fmt.Println("%#v\n", config)

	reqs := [][]reqBody{{}}

	for _,pat := range config.Patterns {
		pattern := []reqBody{}
		for _, req := range pat.Reqs {
			hreq, err := http.NewRequest(req.Method, req.URL, nil)
			if err != nil { panic(err) }
			if req.CType != "" {
				hreq.Header.Add("Content-Type", req.CType)
			}
			reqb := reqBody{"", hreq}
			if req.Data != nil {
				values := make(url.Values)
				for k,v := range req.Data {
					values.Add(k, v)
				}
				reqb.Body = values.Encode()
			}
			pattern = append(pattern, reqb)
		}
		reqs = append(reqs, pattern)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	stpchan := make(chan struct{})
	go callFlood(reqs, config.Workers, stpchan)
	_ = <-c
	stpchan <- struct{}{}
	time.Sleep(50 * time.Millisecond)
	fmt.Print("\nExited normally\n")
}

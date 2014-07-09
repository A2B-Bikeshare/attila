package main

import (
	"fmt"
	"net/http"
	"os"
	//"log"
	"os/signal"
	"runtime"
	"time"
	//"bytes"
	"net/url"
	"strings"
)

var reqsGlobal chan http.Request
var times chan struct{}
var v url.Values

func init() {
	reqsGlobal = make(chan http.Request)
	times = make(chan struct{}, 1)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	values := make(url.Values)
	values.Add("device_id", "10101010101")
	values.Add("fields", "{\"device_id\":\"111111111111111\",\"type\":\"location\",\"data\":{\"num_cell\":0,\"speed\":0,\"bearing\":0,\"num_lac\":0,\"num_sat\":7,\"longitude\":-83.75132845,\"latitude\":42.28357511,\"accuracy\":19,\"num_ap\":0}}")
	values.Add("api_key", "4wz9ajxejfih3ai")
	values.Add("timestamp", "1404921667498")
	values.Add("hash", "2a557b3402062b3f2419acbf1059e3ea0ccb292728f2a7ea2b5d211b11733f38")
	v = values
	api, err := http.NewRequest("POST", "http://dev.vcarl.com/api/checkin", strings.NewReader(values.Encode()))
	api.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	
	reqs := []http.Request{*api}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	stpchan := make(chan struct{})
	go callFlood(reqs, 1, stpchan)
	_ = <-c
	stpchan <- struct{}{}
	time.Sleep(50 * time.Millisecond)
	fmt.Println("Exited normally.")
}

package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"
)

var reqsGlobal chan *http.Request
var times chan struct{}

func init() {
	reqsGlobal = make(chan *http.Request)
	times = make(chan struct{}, 1)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	homepage, err := http.NewRequest("GET", "http://google.com", nil)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	other, err := http.NewRequest("GET", "http://facebook.com", nil)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	reqs := []*http.Request{homepage, other}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	stpchan := make(chan struct{})
	go callFlood(reqs, 50, stpchan)
	_ = <-c
	stpchan <- struct{}{}
	time.Sleep(50 * time.Millisecond)
	fmt.Println("Exited normally.")
}

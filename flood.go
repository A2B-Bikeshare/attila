package main

import (
	"github.com/aybabtme/uniplot/spark"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// read from a channel of requests and execute; don't record anything
func flood(sreqs chan *http.Request, times chan struct{}) {
	cl := &http.Client{}
	for req := range sreqs {
		res, err := cl.Do(req)
		res.Body.Close()
		if err != nil {
			log.Print(err)
		}
		times <- struct{}{}
	}
}

// call flood start the flood, then starts plotting
func callFlood(reqs []*http.Request, concurrency int, stchan chan struct{}) {
	for i := 0; i < concurrency; i++ {
		go flood(reqsGlobal, times)
	}
	plotStop := make(chan struct{})
	go plotTimes(plotStop)

	for {
		select {
		case reqsGlobal <- reqs[rand.Intn(len(reqs))]:
		case <-stchan:
			plotStop <- struct{}{}
			close(reqsGlobal)
			return
		}
	}
}

// plot times starts the sparklines
func plotTimes(stchan chan struct{}) {
	sprk := spark.Spark(time.Millisecond * 100)
	sprk.Units = "reqs"
	sprk.Start()
	for {
		select {
		case _ = <-times:
			sprk.Add(1.0)
		case <-stchan:
			sprk.Stop()
			return
		}
	}
}

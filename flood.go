package main

import (
	"github.com/aybabtme/uniplot/spark"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"io/ioutil"
	"time"
)

// read from a channel of requests and execute; don't record anything
func flood(spatterns chan []http.Request, times chan struct{}) {
	cl := &http.Client{}
	for pattern := range spatterns {
		for _,req := range pattern {
			req.Body = ioutil.NopCloser(strings.NewReader(v.Encode()))
			res, err := cl.Do(&req)
			if err != nil {
				log.Print(err)
				continue
			}
			if res.Body != nil {res.Body.Close()}
			times <- struct{}{}
		}
	}
}

// call flood start the flood, then starts plotting
func callFlood(reqs [][]http.Request, concurrency int, stchan chan struct{}) {
	for i := 0; i < concurrency; i++ {
		go flood(patternsGlobal, times)
	}
	plotStop := make(chan struct{})
	go plotTimes(plotStop)

	for {
		select {
		case patternsGlobal <- reqs[rand.Intn(len(reqs))]:
		case <-stchan:
			plotStop <- struct{}{}
			close(patternsGlobal)
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

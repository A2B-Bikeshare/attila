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
func flood(spatterns chan []reqBody, times chan struct{}, fail chan struct{}) {
	cl := &http.Client{}
	for pattern := range spatterns {
		for _,req := range pattern {
			req.Req.Body = ioutil.NopCloser(strings.NewReader(req.Body))
			res, err := cl.Do(req.Req)
			if err != nil {
				log.Print("%#v\n", req.Req)
				log.Print(err)
				continue
			}
			if !(res.StatusCode <= 200 || res.StatusCode < 300) {
				fail <- struct{}{}
			}
			if res.Body != nil {res.Body.Close()}
			times <- struct{}{}
		}
	}
}

// call flood start the flood, then starts plotting
func callFlood(reqs [][]reqBody, concurrency int, stchan chan struct{}) {
	for i := 0; i < concurrency; i++ {
		go flood(patternsGlobal, times, fail)
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
	var requests int
	var failures int
	for {
		select {
		case _ = <-times:
			sprk.Add(1.0)
			requests++
		case _ = <- fail:
			failures++
		case <-stchan:
			sprk.Stop()
			failProp = float64(failures) / float64(requests)
			return 
		}
	}
}

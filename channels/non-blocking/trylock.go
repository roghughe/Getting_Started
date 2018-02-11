package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
This is a sample that replicates the try lock pattern using Go.

In this scenario we have two potential callers for the same resource - as specified by the jobFunc() - but only
one of them can access it at any one time. If the resource is in use, then don't wait, but do nothing.

This sample sychronises the call to jobFunc() using a channel, when the channel contains something, then the resource
can be used, when the chnnel is empty, then the resource is busy.

*/

type sig int

const (
	OKAY sig = iota + 1
)

var (
	signal = make(chan sig)
)

func init() {
	fmt.Println("Initialising the app")
	rand.Seed(42)
}

func main() {

	// You need to initialise the signal channel, but do it in a go func
	go func() { signal <- OKAY }()

	// Set the first func going, calling every second
	go differentCallers()

	// wait to start the second func
	time.Sleep(500 * time.Millisecond)

	// Set the second a caller func going...
	differentCallers()
}

func differentCallers() {

	// This is the main loop that calls the job func
	for {
		fmt.Println("Okay..")
		time.Sleep(time.Second)
		if tryLock() {
			go jobFunc()
		}
	}
}

func tryLock() bool {
	select {
	case msg := <-signal:
		fmt.Println("received message", msg)
		return true
	default:
		fmt.Println("no message received - doing nothing")
		return false
	}
}

func jobFunc() {

	// This mimics doing some work - pick a duration in millis to wait
	dur := time.Duration(rand.Intn(2000))
	fmt.Printf("In job func.. duration of job will be %d\n", dur)

	time.Sleep(dur * time.Millisecond)

	signal <- OKAY
}

package main

import (
	"fmt"
	"time"

	"gsamples/language"
)

func main() {

	fmt.Println("Hello World")
	goRoutine()
	channel1()
	channel2()
	channel3_channel_close()
	channelSelect()
	defaultSelect()
	language.DemoInterfaces()
}

/**
 * Go Routines
 */
func goRoutine() {
	fmt.Println("goRoutine")
	go say("world")
	say("hello")
}

func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(i, " ", s)
	}
}

/**
 * Go routines and channels. By default, sends and receives block until the other side is ready.
 * This allows goroutines to synchronize without explicit locks or condition variables.
 */
func channel1() {
	fmt.Println("channel1")
	// define a slice
	s := []int{7, 2, 8, -9, 4, 0}

	// mak a channel for integers
	c := make(chan int)
	// sum  the first half of the slice
	go sum(s[:len(s)/2], c)
	// sum the second half of the slice
	go sum(s[len(s)/2:], c)
	// recieve what's bee stuffed in the channel
	x, y := <-c, <-c // receive from c writing to x and then y

	fmt.Println("x: ", x, " y: ", y, x+y)
}

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	fmt.Println("New sum: ", sum, " being added to the channel")
	c <- sum // send sum to c
}

/**
* buffered channels
 Channels can be buffered. Provide the buffer length as the second argument to make to initialize a buffered channel:

 ch := make(chan int, 100)

 Sends to a buffered channel block only when the buffer is full. Receives block when the buffer is empty.
*/
func channel2() {

	fmt.Println("channel2")
	ch := make(chan string, 6)
	ch <- "1"
	ch <- "23456"
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}

/**
 * Demonstrates closing a channelto indicate that no more values will be sent
 */
func channel3_channel_close() {
	fmt.Println("channel3_channel_close")
	c := make(chan int, 10)
	go fibonacci(cap(c), c)
	for i := range c {
		fmt.Println(i)
	}
}

func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	fmt.Println("Closing channel")
	close(c)
}

/**
 * The select statement lets a goroutine wait on multiple communication operations.
 * A select blocks until one of its cases can run, then it executes that case.
 * It chooses one at random if multiple are ready.
 */
func channelSelect() {
	fmt.Println("channelSelect")

	c := make(chan int)
	quit := make(chan int)
	fmt.Println("before - go func()")
	go func() {
		fmt.Println("anonymous func")
		for i := 0; i < 10; i++ {
			fmt.Println("anon read: ", <-c)
		}
		quit <- 0
	}()
	fibonacci2(c, quit)
}

func fibonacci2(c, quit chan int) {
	x, y := 0, 1
	fmt.Println("fibonacci2 x: ", x, " y: ", y)
	for {
		select {
		case c <- x:
			fmt.Println("write x into c: ", x)
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func defaultSelect() {

	fmt.Println("defaultSelect")
	// trigger every 100ms
	tick := time.Tick(100 * time.Millisecond)
	// trigger after 500ms
	boom := time.After(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			fmt.Println("tick.")
		case <-boom:
			fmt.Println("BOOM!")
			return
		default:
			fmt.Println("    .")
			// sleep 50ms
			time.Sleep(50 * time.Millisecond)
		}
	}
}

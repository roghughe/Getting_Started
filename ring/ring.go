package main

import (
	"container/ring"
	"fmt"
)

func main() {
	// Create a new ring of size 5
	r := ring.New(5)

	// Get the length of the ring
	n := r.Len() * 2

	// Initialize the ring with some bytes values
	bytes := []byte("These are some bytes")
	for i := 0; i < n; i++ {
		r.Value = bytes
		r = r.Next()
	}

	// Iterate through the ring and print its contents
	r.Do(func(p interface{}) {
		fmt.Println(p.([]byte))
	})

}

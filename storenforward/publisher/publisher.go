package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gsamples/types"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	counter = 0 // The message number
)

const (
	CONTENT = "How do I send a JSON string in a POST request in Go"
	URL     = "http://localhost:7868/message?topic=Bernie"
)

func main() {

	// loop around
	// create a message
	// send it
	// wait

	fmt.Println("Simple publisher, URL:>", URL)

	for {

		// Create a message
		msg := types.Message{
			Id:      counter,
			Content: CONTENT,
			Time:    time.Now(),
		}

		mbytes, err := json.Marshal(msg)

		if err != nil {
			fmt.Printf("JSON marshall error: %+v\n", err)
		}

		// This uses a Request as it gives you more control than a http.Post()
		req, err := http.NewRequest("POST", URL, bytes.NewBuffer(mbytes))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{} // For example you can setuo client values such as time out if necessary
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error in Posting to server: %+v\n", err)
		} else {

			// Now display the response:
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Printf("Response status %s, Headers: %v, Body: %s\n", resp.Status, resp.Header, body)
			resp.Body.Close()
		}

		// this could be a random delay
		time.Sleep(time.Second)

		counter++
	}

}

package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type store struct {
	body []byte // the body of the message
}

// This describes where to reply
type subscriber struct {
	id   int    // a unique id for the subscriber
	port int    // The port number on where to reply
	host string // The host to which we reply
}

const (
	IN_PATTERN  = "/message"
	SUB_PATTERN = "/subscribe"
	OUT_PATTERN = "/forward"
	PORT        = ":7868"
)

// This is the store and forward. To work it needs a list of clients (host names) In this simple example, the client must
// honour the incoming client API (URL Pattern).
// The client subscribes / registers with the storenforward. In this simple example, clients don't unregister.

func main() {

	fmt.Printf("Starting - store and forward - listening on port %s for pattern %s\n", PORT, IN_PATTERN)
	http.HandleFunc(IN_PATTERN, processIncomingMessage)
	http.HandleFunc(SUB_PATTERN, addSubscriber)

	http.ListenAndServe(PORT, nil)
}

// This function does the store and forward
func processIncomingMessage(w http.ResponseWriter, r *http.Request) {

	// This optimisation ignores anything but POSTs
	if r.Method != "POST" {
		fmt.Println("recieved a non-POST request")
		err := errors.New("Unsupported request method")
		http.Error(w, err.Error(), 404)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error reading body: %+v\n", err)
	}

	// At this point send the body back all subscribers
	fmt.Printf("", body)

	io.WriteString(w, "OK")
}

/*
To subscribe, the caller must:
1) provide a listen port (but only so that I can run this on the same machine)
2) a unique id.
3) the return URL pattern

The subscriber must know how to unmarshall the message
*/
func addSubscriber(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		fmt.Println("recieved a non-GET request")
		err := errors.New("Unsupported request method")
		http.Error(w, err.Error(), 404)
		return
	}

}

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type replyMsg struct {
	replyTo *url.URL // When to send the request
	body    *[]byte  // What to send this time
}

// This describes where to reply
type subscriber struct {
	reply url.URL       // key = id, value = subscriber information. Added for clarity only.
	ch    chan replyMsg // Where to send the replay. Use of a channel is more complex, but will maintain message order.
}

type subscribers map[int]subscriber // All the subscribers, by id, for a topic

const (
	IN_PATTERN  = "/message"
	SUB_PATTERN = "/subscribe"
	PORT        = ":7868"
	BAD_REQUEST = 400
)

var (
	submap = make(map[string]subscribers) // record the subscribers for a topic by id
	strmap = make(map[string][]byte)      // record the data we're storing

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

	// Store the data
	topic := r.URL.Query().Get("topic")
	if topic == "" {
		fmt.Println("Missing Topic")
		http.Error(w, "Missing Topic", BAD_REQUEST)
		return
	}

	strmap[topic] = body

	// Now send the data to any clients.
	subscribers := submap[topic]

	for id, subs := range subscribers {

		msg := replyMsg{
			replyTo: &subs.reply,
			body:    &body,
		}
		fmt.Printf("forwarding to : %d\n", id)
		subs.ch <- msg
	}

	io.WriteString(w, "OK")
}

/*
To subscribe, the caller must:
1) Somewhere to reply to..
2) a unique id.
3) the topic to which its subscribing.

The subscriber must know how to unmarshall the message
*/
func addSubscriber(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		fmt.Println("recieved a non-GET request")
		err := errors.New("Unsupported request method")
		http.Error(w, err.Error(), 404)
		return
	}

	ids := r.URL.Query().Get("id")

	id, err := strconv.Atoi(ids)
	if err != nil {
		fmt.Printf("Cannot decode id: %+v\n", err)
		http.Error(w, "Invalid id - "+err.Error(), BAD_REQUEST)
		return
	}

	rep := r.URL.Query().Get("replyto")
	if rep == "" {
		fmt.Println("Missing ReplyTo")
		http.Error(w, "Missing replyTo", BAD_REQUEST)
		return
	}

	reply, err := url.Parse(rep)
	if err != nil {
		fmt.Println("Missing ReplyTo")
		http.Error(w, "Cannot Parse replyTo", BAD_REQUEST)
		return
	}

	topic := r.URL.Query().Get("topic")
	if topic == "" {
		fmt.Println("Missing Topic")
		http.Error(w, "Missing Topic", BAD_REQUEST)
		return
	}

	subs, ok := submap[topic]
	if !ok {
		subs = make(subscribers)
		submap[topic] = subs
	}

	s := subscriber{
		reply: *reply,              // Save the reply
		ch:    make(chan replyMsg), // Create somewhere to send the data message
	}

	subs[id] = s
	// start listening for messages
	go s.forward()

	// Okay, so now we're subscribed....

	// Send what's in the store
	if val, ok := strmap[topic]; ok {
		sendData(val, &s)
	} else {
		fmt.Println("No existing data... for topic")
	}
}

// Read the messages from the channel and send on via HTTP.
func (s *subscriber) forward() {

	fmt.Println("listening for messages...")
	for {

		msg := <-s.ch
		fmt.Printf("received message. Replying to: %s, data: %s\n", msg.replyTo.String(), string(*msg.body))

		// This uses a Request as it gives you more control than a http.Post()
		req, err := http.NewRequest("POST", s.reply.String(), bytes.NewBuffer(*msg.body))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{} // For example you can setuo client values such as time out if necessary
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error in Posting to subscriber: %+v\n", err)
		} else {
			// Now display the response:
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Printf("Response status %s, Headers: %v, Body: %s\n", resp.Status, resp.Header, body)
			resp.Body.Close()
		}

	}
}

// Send the message to the subscriber.
func sendData(body []byte, s *subscriber) {

	rm := replyMsg{
		body:    &body,
		replyTo: &s.reply,
	}

	s.ch <- rm // Put the message in the channel to be picked up bu forward().
}

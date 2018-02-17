package main

import (
	"bytes"
	"container/ring"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
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

// defines a ring (circular list) and a mutex to lock access to it
type ringBuf struct {
	buf *ring.Ring  // The ring buffer
	mut *sync.Mutex // Something to lock access to the buffer...
}

const (
	IN_PATTERN  = "/message"   // URL used to receive the data
	SUB_PATTERN = "/subscribe" // URL used to subscribe to the data
	PORT        = ":7868"
	BAD_REQUEST = 400 // Simple HTTP status code
	BUFF_SIZE   = 40  // Allows us to keep this many messages in memory.
)

var (
	submap = make(map[string]subscribers) // record the subscribers for a topic by id

	strmap = make(map[string]ringBuf) // record the data we're storing
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

// This function does the store and forward.
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

	addToStore(body, topic)
	updateSubscribers(body, topic)
	io.WriteString(w, "OK")
}

// Add the latest data to the store, lazily initialising the Ring ptr.
func addToStore(body []byte, topic string) {

	fmt.Printf("About to store bytes len: %d,  --- %v\n", len(body), body)

	rb, ok := strmap[topic]
	if !ok {
		rb = ringBuf{
			buf: ring.New(BUFF_SIZE),
			mut: &sync.Mutex{},
		}
	}

	rb.mut.Lock()
	rb.buf.Value = body
	r := rb.buf.Next()
	rb.buf = r
	strmap[topic] = rb
	rb.mut.Unlock()
}

// Now send the data to any clients.
func updateSubscribers(body []byte, topic string) {
	subscribers := submap[topic]

	for id, subs := range subscribers {

		msg := replyMsg{
			replyTo: &subs.reply,
			body:    &body,
		}
		fmt.Printf("forwarding to : %d\n", id)
		subs.ch <- msg
	}
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

	// Send back existing data...
	go updateWithExisting(topic, &s)

	// Okay, so now we're subscribed....

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

// takes what we have and sends it to any new subscribers
func updateWithExisting(topic string, s *subscriber) {
	// Send what's in the store
	if rb, ok := strmap[topic]; ok {
		rb.mut.Lock()
		defer rb.mut.Unlock()
		if rb.buf.Len() > 0 {

			// Iterate through the ring and send its contents
			rb.buf.Do(func(val interface{}) {
				if val != nil {
					sendData(val.([]byte), s)
				}
			})
		} else {
			fmt.Println("Topic with an empty buffer???")
		}
	} else {
		fmt.Println("No existing data... for topic")
	}
}

// Send the message to the subscriber. The data is copied for thread safety.
func sendData(body []byte, s *subscriber) {

	rm := replyMsg{
		body:    &body,
		replyTo: &s.reply,
	}

	s.ch <- rm // Put the message in the channel to be picked up bu forward().
}

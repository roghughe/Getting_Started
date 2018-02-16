package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gsamples/types"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	PATTERN = "/forward"
	PORT    = 7868
	TOPIC   = "Bernie"
)

var (
	randomSeq *rand.Rand
	port      int
)

// The default number generator is deterministic, so it'll
// produce the same sequence of numbers each time by default.
// To produce varying sequences, give it a seed that changes.
// The random number is to create a unique id
func init() {
	// If we run lots of these,they could clash, but thisi s only a sample
	seed := rand.NewSource(time.Now().UnixNano())
	randomSeq = rand.New(seed)

	port = PORT + randomSeq.Intn(100)

}

func main() {
	fmt.Printf("Starting - subscriber - listening on port %d for pattern %s\n", port, PATTERN)

	if !subscribe() {
		os.Exit(1)
	}

	http.HandleFunc(PATTERN, processMessage)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

/*
To subscribe, the caller must:
1) provide a listen port (but only so that I can run this on the same machine)
2) a unique id.
3) the return URL pattern

The subscriber must know how to unmarshall the message
*/
func subscribe() bool {

	u, err := createURL()
	if err != nil {
		return false
	}

	fmt.Printf("Subscribing with %s\n", u)

	response, err := http.Get(u)
	if err != nil {
		fmt.Printf("Subscribe error: %s - exiting.", err.Error())
		return false
	}

	if response.StatusCode != 200 {
		fmt.Printf("Storenfor replied with invalid status code: %d - exiting.", response.StatusCode)
		return false
	}

	fmt.Println("Subscribed okay")
	response.Body.Close()
	return true
}

func createURL() (string, error) {

	replyTo := "http://localhost:" + strconv.Itoa(port) + PATTERN + "?topic=" + TOPIC

	sendTo, err := url.Parse("http://localhost:" + strconv.Itoa(PORT) + "/subscribe")
	if err != nil {
		fmt.Printf("Error parsing sendTo: %+v", err)
		return "", err
	}

	parameters := url.Values{}
	id := randomSeq.Intn(100)
	parameters.Add("id", strconv.Itoa(id))
	parameters.Add("topic", TOPIC)
	parameters.Add("replyto", replyTo)
	sendTo.RawQuery = parameters.Encode()

	fmt.Printf("Encoded URL is %q\n", sendTo.String())

	return sendTo.String(), nil
}

// Once subscribed, then the storenforward will send messages here
func processMessage(w http.ResponseWriter, r *http.Request) {

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

	var msg types.Message
	err = json.Unmarshal(body, &msg)
	if err != nil {
		fmt.Printf("Error unmarshalling body: %+v\n", err)
	}

	fmt.Printf("Message is: %+v\n", msg)

	io.WriteString(w, "OK")
}

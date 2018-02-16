package types

import "time"

type Message struct {
	Topic   string    // Some name for the message
	Id      int       // The message number
	Content string    // The contents... could be anything
	Time    time.Time // The time that the message was sent.
}

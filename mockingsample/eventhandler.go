/*
This package models any old event handler.


*/
package mockingsample

import (
	"fmt"
	"gsamples/mockingsample/dbaccess"
)

// The event handler struct. Models event attributes
type EventHandler struct {
	name  string                          // The name of the event
	actor dbaccess.SomeFunctionalityGroup // The interface for our dbaccess fucntions
}

// This creates a event handler instance, using whatever name an actor  are passed in.
func NewEventHandler(actor dbaccess.SomeFunctionalityGroup, name string) EventHandler {

	return EventHandler{
		name:  name,
		actor: actor,
	}
}

// This is a sample event handler - it reads from the DB does some imaginary business logic and writes the results back
// to the DB.
func (eh *EventHandler) HandleSomeEvent(action string) error {

	fmt.Printf("Handling event: %s\n", action)
	value, err := eh.actor.ReadSomething(action, "arg1")
	if err != nil {
		fmt.Printf("Use the logger to log your error here. The read error is: %+v\n", err)
		return err
	}

	// Do some business logic here
	if len(value) == 2 && value[0] == "Hello" {
		value[1] = "World"
	}

	// Now write the result back to the database
	err = eh.actor.WriteSomething(value)

	if err != nil {
		fmt.Printf("Use the logger to log your error here. The write error is: %+v\n", err)
	}

	return err
}

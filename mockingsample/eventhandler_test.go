package mockingsample

import (
	"errors"
	"fmt"
	"testing"
)

// A DummyConnection mimics a real database connection - but allows us to mock the connection and follow happy and fail code paths
type DummyConnection struct {

	fail bool // Set this to true to mimic a DB read / write failure
}

var dummyError = errors.New("This is my test error")

/* This section contains the mock database access functions */

// This is the mock database read function
func (r *DummyConnection) ReadSomething(arg0, arg1 string) ([]string, error)  {
	fmt.Printf("This is the MOCK database driver - read args: %s -- %s",arg0,arg1)

	if r.fail {
		fmt.Println("Whoops - there's been a database write error")
		return []string{}, dummyError
	}

	return []string{"Hello", ""}, nil
}

// This is mock database write function
func (r * DummyConnection) WriteSomething(arg0 []string) error {
	fmt.Printf("This is the MOCK database driver - write args: %v",arg0)

	if r.fail {
		fmt.Println("Whoops - there's been a database write error")
		return dummyError
	}

	return nil
}


/* Now create the tests */

// Test calling the event handler with a dummy database connection.
func TestEventHandlerDB_happy_flow(t *testing.T) {

	testCon := DummyConnection{
		fail: false,
	}

	eh := NewEventHandler(&testCon,"Happy")


	err := eh.HandleSomeEvent("Action")

	if err != nil {
		t.Errorf("Failed - with error: %+v\n", err)
	}
}


// Test calling the event handler with a dummy database connection, for the failure flow.
func TestEventHandlerDB_fail_flow(t *testing.T) {

	testCon := DummyConnection{
		fail: true,
	}

	eh := NewEventHandler(&testCon,"Fail")


	err := eh.HandleSomeEvent("Action 2")

	if err == nil {
		t.Errorf("Failed - with error: %+v\n", err)
	}
}
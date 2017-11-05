package mockingsample

import (
	"errors"
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

// A MockConnection mimics a real database connection - but allows us to mock the connection and follow happy and fail code paths
type MockConnection struct {
	fail bool // Set this to true to mimic a DB read / write failure
}

var dummyError = errors.New("this is my test error")

/* This section contains the mock database access functions */

// This is the mock database read function
func (r *MockConnection) ReadSomething(arg0, arg1 string) ([]string, error) {
	fmt.Printf("This is the MOCK database driver - read args: %s -- %s\n", arg0, arg1)

	if r.fail {
		fmt.Println("Whoops - there's been a database write error")
		return []string{}, dummyError
	}

	return []string{"Hello", ""}, nil
}

// This is mock database write function
func (r *MockConnection) WriteSomething(arg0 []string) error {
	fmt.Printf("This is the MOCK database driver - write args: %v\n", arg0)

	if r.fail {
		fmt.Println("Whoops - there's been a database write error")
		return dummyError
	}

	return nil
}

/* Now create the tests */

// Test calling the event handler with a dummy database connection.
func TestEventHandlerDB_happy_flow(t *testing.T) {

	testCon := MockConnection{
		fail: false,
	}

	eh := NewEventHandler(&testCon, "Happy")

	err := eh.HandleSomeEvent("Action")

	assert.NotNil(t,"Failed - with error: %+v\n", err)
}

// Test calling the event handler with a dummy database connection, for the failure flow.
func TestEventHandlerDB_fail_flow(t *testing.T) {

	testCon := MockConnection{
		fail: true,
	}

	eh := NewEventHandler(&testCon, "Fail")

	err := eh.HandleSomeEvent("Action 2")

	assert.Nil(t,err)
}

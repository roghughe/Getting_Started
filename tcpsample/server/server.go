package main

/*
 Sample server that reads messages where the frist four bytes are the length of data to read (BigEndian). Reads the message
 and displays it.
 */


import (
	"fmt"
	"net"
	"os"
	"encoding/binary"
)

// Some defaults
const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func main() {

	fmt.Println("Sample Server Starting...")
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data len.
	// which is the first four bytes
	buf := make([]byte, 4)

	for {
		// Read the incoming header into the buffer.
		reqLen, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading header bytes:", err.Error())
			break;
		}

		if reqLen != 4 {
			fmt.Printf("Only read %d bytes as header - error", reqLen)
			continue   // This is probably an exit type error here...
		}

		// This will be the amount of data to read
		dataLen := binary.BigEndian.Uint32(buf)
		dataBuf := make([]byte,dataLen)

		// read the actual message
		reqLen, err  = conn.Read(dataBuf)

		if reqLen != int(dataLen) || err != nil {
			fmt.Printf("Cannot read bytes. Num read %d, expected %d, err %+v\n", dataLen, reqLen, err)
		} else if reqLen == 0 {
			fmt.Println("No data to read - header lenght says 0")
		} else {
			// everything went okay

			fmt.Printf("Read %d bytes\n", reqLen)
			s := string(dataBuf[:reqLen])
			fmt.Printf("Message was: %s\n", s)
			if s == "bye" {
				break;
			}
		}
	}
	// Close the connection when you're done with it.
	conn.Close()
}

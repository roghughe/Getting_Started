package main

/*
 Sample client that sends a message - either from the command line or the console - to a server. The message length
 is sent first as four bytes (BigEndian) to tell the server how much more is coming...
 */

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"encoding/binary"
)

const (
	HOST  = "localhost:3333"
	STDIN = "CON"
	BYE = "bye"
)


func main() {

	hostPtr := flag.String("host", HOST, "the host to connect to in \"host:port\" format")
	msgPtr := flag.String("message", STDIN, "the message to send. Will default to console input.")
	flag.Parse()

	fmt.Printf("Sample TCP client - running with args: %s and %s\n",*hostPtr,*msgPtr)

	// Connect and exit on error
	conn, err := net.Dial("tcp", *hostPtr)
	checkError(err)

	defer conn.Close()

	if *msgPtr == STDIN {
		fmt.Println("Running the client in console mode.....")
		for {
			reader := bufio.NewReader(os.Stdin)
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				break
			}

			line = strings.TrimRight(line, " \t\r\n")

			if !writeString(line, conn) {
				break
			}

			if line == BYE {
				// Maybe this shouldn't be sent... just use it as a way out of the loop.
				fmt.Println("Exiting...")
				break
			}

		}
	} else {
		writeString(*msgPtr,conn)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func writeString(line string, conn net.Conn) bool {

	bytes := []byte(line)
	l := len(bytes)
	ul := uint32(l)
	head := make([]byte,4)

	binary.BigEndian.PutUint32(head,ul)

	if !writeBytes(head,conn) {
		return false
	}

	writeBytes(bytes,conn)
	return true
}

func 	writeBytes(bytes []byte, conn net.Conn) bool {
	_, err := conn.Write(bytes)
	if err != nil {
		fmt.Println(err)
		return false

	}
	return true
}
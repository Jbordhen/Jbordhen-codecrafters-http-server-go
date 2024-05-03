package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")

	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	connection, err := l.Accept()

	handleHTTPResponse(connection)

	defer connection.Close()

}

func handleHTTPResponse(connection net.Conn) {
	buffer := make([]byte, 255)

	_, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
		return
	}

	input := string(buffer)

	lines := strings.Split(input, "\r\n")

	path := strings.Split(lines[0], " ")[1]

	var response []byte

	if path == "/" {
		response = []byte("HTTP/1.1 200 OK\r\n\r\n")
	} else {
		response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}

	_, err = connection.Write(response)

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
	}
}

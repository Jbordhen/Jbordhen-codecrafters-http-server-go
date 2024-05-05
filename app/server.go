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

	for {
		connection, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go handleHTTPResponse(connection)
	}
}

func handleHTTPResponse(connection net.Conn) {
	defer connection.Close()
	buffer := make([]byte, 255)

	_, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
		return
	}

	var response []byte = generateResponse(buffer)

	_, err = connection.Write(response)

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
		return
	}
}

func generateResponse(buffer []byte) []byte {
	input := string(buffer)

	lines := strings.Split(input, "\r\n")

	httpInfo := strings.Split(lines[0], " ")

	httpMethod := httpInfo[0]

	var body string

	if httpMethod == "POST" {
		body = string(lines[2])
	}

	path := httpInfo[1]

	var response []byte

	if path == "/" {
		response = []byte("HTTP/1.1 200 OK\r\n\r\n")
		return response
	}

	if strings.Contains(path, "/echo") {
		content := strings.Split(path, "/echo/")[1]

		response = []byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(content), content))
		return response
	}

	if path == "/user-agent" {
		userAgent := strings.Split(lines[2], ": ")[1]
		response = []byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent))
		return response
	}

	if strings.HasPrefix(path, "/files/") {
		if httpMethod == "GET" {
			content, found := strings.CutPrefix(path, "/files/")

			if !found || len(os.Args) < 3 {
				response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
				return response
			}

			directory := os.Args[2]

			file, err := os.ReadFile(fmt.Sprintf("%s/%s", directory, content))

			if err != nil {
				response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
				return response
			}

			response = []byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(file), file))
			return response
		}

		if httpMethod == "POST" {
			fileName, found := strings.CutPrefix(path, "/files/")

			if !found || len(os.Args) < 3 {
				response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
				return response
			}

			if len(lines) > 6 {
				body = (strings.Trim(lines[6], "\r\n"))
			}

			directory := os.Args[2]
			body := strings.Replace(body, "\x00", "", -1)

			err := os.WriteFile(fmt.Sprintf("%s/%s", directory, fileName), []byte(body), 0644)

			if err != nil {
				response = []byte("HTTP/1.1 500 Error creating file\r\n\r\n")
			}

			response = []byte(fmt.Sprintf("HTTP/1.1 201 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(body), body))
			return response
		}
	}

	response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")

	return response
}

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	network = "tcp"
	port    = ":8080"
)

func StartClient() {
	conn, err := buildClientConnection()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	messageReader := bufio.NewReader(os.Stdin)
	connReader := bufio.NewReader(conn)

	errors := make(chan error)
	messages := make(chan string)

	go readServerMessage(connReader, errors, messages)
	go writeToServerMessage(messageReader, conn, errors)

	for {
		select {
		case msg := <-messages:
			fmt.Println("server:", strings.TrimSpace(msg))

		case err := <-errors:
			fmt.Println("error:", err)
			return
		}
	}
}

func buildClientConnection() (net.Conn, error) {
	tcpAddr, err := net.ResolveTCPAddr(network, port)
	if err != nil {
		return nil, err
	}

	return net.DialTCP(network, nil, tcpAddr)
}

func readServerMessage(connReader *bufio.Reader, errors chan<- error, messages chan<- string) {
	for {
		message, err := connReader.ReadString('\n')
		if err != nil {
			errors <- err
			return
		}
		messages <- message
	}
}

func writeToServerMessage(messageReader *bufio.Reader, conn net.Conn, errors chan<- error) {
	for {
		data, err := messageReader.ReadString('\n')
		if err != nil {
			errors <- err
		}

		message := strings.TrimSpace(data)
		message += "\n"
		conn.Write([]byte(message))
	}
}

func main() {
	StartClient()
}

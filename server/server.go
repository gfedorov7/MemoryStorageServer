package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync/atomic"
)

var hasConnection atomic.Bool

const (
	network = "tcp"
	port    = ":8080"
)

func StartServer() {
	listener, err := buildServerConnection()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = acceptConnection(listener)
	if err != nil {
		fmt.Println(err)
	}
}

func buildServerConnection() (*net.TCPListener, error) {
	tcpAddr, err := net.ResolveTCPAddr(network, port)
	if err != nil {
		return nil, err
	}

	return net.ListenTCP(network, tcpAddr)
}

func acceptConnection(listener *net.TCPListener) error {
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		if !hasConnection.CompareAndSwap(false, true) {
			conn.Write([]byte("server already has active connection\n"))
			conn.Close()
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	defer hasConnection.Store(false)

	connReader := bufio.NewReader(conn)
	for {
		data, err := connReader.ReadString('\n')
		if err != nil {
			fmt.Println("client disconnected:", err)
			return
		}

		split := strings.Split(strings.TrimSpace(data), " ")
		if len(split) <= 1 {
			conn.Write([]byte("invalid command\n"))
			continue
		}

		commandAnswer := commandHandler(split[0]) + "\n"
		conn.Write([]byte(commandAnswer))
	}
}

func commandHandler(command string) string {
	switch command {
	case "GET":
		return "receive command get"
	case "SET":
		return "receive command set"
	default:
		return "unknow command"
	}
}

func main() {
	StartServer()
}

package main

import (
	"MemoryStorageServer/cmd"
	"MemoryStorageServer/collection"
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync/atomic"
	"time"
)

var hasConnection atomic.Bool

const (
	network = "tcp"
	port    = ":8080"
)

func StartServer(storage collection.AsyncCollectionInterface) {
	listener, err := buildServerConnection()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = acceptConnection(listener, storage)
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

func acceptConnection(listener *net.TCPListener, storage collection.AsyncCollectionInterface) error {
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

		go handleConnection(conn, storage)
	}
}

func handleConnection(conn net.Conn, storage collection.AsyncCollectionInterface) {
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

		commandHandler(storage, split, conn)
	}
}

func commandHandler(storage collection.AsyncCollectionInterface, split []string, conn net.Conn) {
	command := split[0]
	commandAnswer, err := cmd.CommandHandler(storage, command, split[1:])
	if err != nil {
		conn.Write([]byte(err.Error() + "\n"))
	} else if commandAnswer != nil {
		conn.Write([]byte(commandAnswer.String() + "\n"))
	} else {
		conn.Write([]byte(command + " success" + "\n"))
	}
}

func main() {
	storage := collection.NewAsyncCollection()
	go storage.StartJanitor(time.Duration(100) * time.Second)

	StartServer(storage)
}

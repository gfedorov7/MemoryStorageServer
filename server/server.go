package main

import (
	"MemoryStorageServer/collection"
	"bufio"
	"fmt"
	"net"
	"strconv"
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
		if len(split) <= 1 {
			conn.Write([]byte("invalid command\n"))
			continue
		}

		command := split[0]
		commandAnswer, err := commandHandler(storage, command, split[1:])
		if err != nil {
			conn.Write([]byte(err.Error() + "\n"))
		} else if commandAnswer != nil {
			conn.Write([]byte(commandAnswer.String() + "\n"))
		} else {
			conn.Write([]byte(command + " success" + "\n"))
		}
	}
}

func commandHandler(storage collection.AsyncCollectionInterface, command string,
	args []string) (*collection.MemoryCollection, error) {
	switch command {
	case "GET":
		if len(args) < 1 {
			return nil, fmt.Errorf("GET command wait 1 arg")
		}
		val, err := storage.Get(args[0])
		if err != nil {
			return nil, err
		}
		return &val, nil
	case "SET":
		if len(args) < 3 {
			return nil, fmt.Errorf("SET command wait 2 arg")
		}
		num, err := strconv.Atoi(args[2])
		if err != nil {
			return nil, err
		}
		memoryCollection, err := collection.Create(args[1], time.Duration(num)*time.Second, time.Now())
		if err != nil {
			return nil, err
		}
		storage.Set(args[0], memoryCollection)
		return nil, nil
	default:
		return nil, fmt.Errorf("unknow command")
	}
}

func main() {
	storage := collection.NewAsyncCollection()
	go storage.StartJanitor(time.Duration(100) * time.Second)

	StartServer(storage)
}

package parser

import (
	"flag"
)

const (
	network = "tcp"
	address = "127.0.0.1:8080"
)

type ConsoleFlag struct {
	Network string
	Address string
}

func ParseFlag(consoleFlag *ConsoleFlag) {
	flag.StringVar(&consoleFlag.Network, "network", network, "network protocol")
	flag.StringVar(&consoleFlag.Address, "address", address, "address")

	flag.Parse()
}

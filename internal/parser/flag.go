package parser

import (
	"flag"
)

type ConsoleFlag struct {
	Network string
	Address string
}

func ParseFlag(consoleFlag *ConsoleFlag) {
	flag.StringVar(&consoleFlag.Network, "network", "tcp", "network protocol")
	flag.StringVar(&consoleFlag.Address, "address", ":8080", "address")

	flag.Parse()
}

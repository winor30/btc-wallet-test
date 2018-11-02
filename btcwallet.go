package main

import (
	"fmt"
	"github.com/winor30/btc-wallet-test/service"
	"log"
	"os"
	"strconv"
)

func initialize() {
	if len(os.Args) != 3 {
		log.Fatalln("init command is invalid!")
		return
	}
	net := os.Args[2]
	service.Initialize(net)
}

func account() {
	if len(os.Args) < 3 {
		log.Fatalln("account command is invalid!")
	}

	switch os.Args[2] {
	case "add":
		name := os.Args[3]
		net := os.Args[4]
		service.AddAccount(name, net)
	}
}

func address() {
	if len(os.Args) < 3 {
		log.Fatalln("account command is invalid!")
	}

	switch os.Args[2] {
	case "receive":
		name := os.Args[3]
		net := os.Args[4]
		i, _ := strconv.Atoi(os.Args[5])
		index6 := uint32(i)
		address, _ := service.Receive(name, net, index6)
		fmt.Println(address)
	case "change":
		name := os.Args[3]
		net := os.Args[4]
		i, _ := strconv.Atoi(os.Args[5])
		index6 := uint32(i)
		address, _ := service.Change(name, net, index6)
		fmt.Println(address)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("command is invalid!")
		return
	}
	service.Setup()
	cmd := os.Args[1]
	switch cmd {
	case "init":
		initialize()
	case "account":
		account()
	case "address":
		address()
	}
}

package kademlia_node

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var me *Contact
var isRunning bool

func CliInit(node *Node) {
	me = node.Me
	isRunning = true
	CliMain()
}

func CliParsing() []string {
	var s string
	var err error

	r := bufio.NewReader(os.Stdin)
	for s != "" {
		fmt.Println(">>> ")
		s, err = r.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
	}

	return strings.Fields(s)
}

func CliMain() {
	for isRunning {
		args := CliParsing()

		switch args[0] {
		case "ping":
			Ping(args)
		case "exit":
			Exit(args)
		case "put":
			Put(args)
		case "get":
			Get(args)

		}
	}
}

func Ping(args []string) {
	ip := args[1]
	port, _ := strconv.Atoi(args[2])
	target := NewContact(NewRandomKademliaID(), ip, port) // FIXME
	networkInstance.MessageHandler.SendPingRequest(me, target)
}

func Exit(args []string) {
	isRunning = false
	// TODO
}

func Put(args []string) {
	// TODO
}

func Get(args []string) {
	// TODO
}

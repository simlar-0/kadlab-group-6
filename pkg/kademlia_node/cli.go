package kademlia_node

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var me *Node
var isRunning bool

func CliInit(node *Node) {
	me = node
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
		input := CliParsing()
		args := input[1:]
		cmd := input[0]

		switch cmd {
		case "put":
			r := Put(args)
			fmt.Println(r) // debug only
		case "get":
			r, n := Get(args)
			fmt.Println(r, n) // debug only
		case "exit":
			Exit(args)
		}
	}
}

// takes a single argument, the contents of the file you are uploading, and outputs the
// hash of the object, if it can be uploaded successfully.
func Put(args []string) *KademliaID {

	if len(args) > 1 {
		fmt.Printf("Put command only take 1 argument: put <file content>")
		return nil
	}

	content := []byte(args[0])
	key, err := me.Store(content)

	if err != nil {
		fmt.Printf("Couldn't store data")
		return nil
	}

	return key
}

// takes a hash as its only argument, and outputs the contents of the object and the
// node it was retrieved from, if it could be downloaded successfully.
func Get(args []string) ([]byte, *Node) {
	if len(args) > 1 {
		panic("Get command only take 1 argument: get <file hash>")
	}

	hash := args[0]
	data, node, err := me.LookupData(hash)

	if err != nil {
		fmt.Printf("Couldn't find data")
		return nil, nil
	}

	return data, node
}

// terminates the node
func Exit(args []string) {
	isRunning = false
	// TODO
}

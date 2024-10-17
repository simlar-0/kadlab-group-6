package kademlia_node

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Cli struct {
	me        NodeInterface
	isRunning bool
}

func CliInit(node NodeInterface) *Cli {
	c := &Cli{
		me:        node,
		isRunning: true,
	}
	return c
}

func (cli *Cli) Main() {
	fmt.Println("\033[0;31mCli is running \033[0m")

	var s string
	var err error

	r := bufio.NewReader(os.Stdin)
	for cli.isRunning {

		fmt.Printf(">>> ")
		s, err = r.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		} else {
			input := strings.Fields(s)
			cmd := input[0]
			args := input[1:]

			cli.Run(cmd, args)
		}
	}
}

func (cli *Cli) Run(cmd string, args []string) {
	switch cmd {
	case "help":
		cli.Help()
	case "put":
		r := cli.Put(args)
		fmt.Println(r)
	case "get":
		r, n := cli.Get(args)
		fmt.Println(string(r), n)
	case "exit":
		cli.Exit(args)
	case "ping":
		fmt.Println("pong")
	}

}

// takes a single argument, the contents of the file you are uploading, and outputs the
// hash of the object, if it can be uploaded successfully.
func (cli *Cli) Put(args []string) *KademliaID {
	if len(args) < 1 {
		fmt.Println("Put command only take 1 argument: put <file content>")
		return nil
	}

	content := []byte(args[0])
	key, err := cli.me.Store(content)

	if err != nil {
		fmt.Printf("Couldn't store data")
		return nil
	}

	fmt.Println("Data stored, here is your key")
	return key
}

// takes a hash as its only argument, and outputs the contents of the object and the
// node it was retrieved from, if it could be downloaded successfully.
func (cli *Cli) Get(args []string) ([]byte, *Contact) {
	if len(args) < 1 {
		fmt.Println("Get command only take 1 argument: get <file hash>")
	}

	hash := args[0]
	data, contact, err := cli.me.LookupData(hash)

	if err != nil {
		fmt.Printf("Couldn't find data:")
		return nil, nil
	}

	fmt.Println("Data retrieved, here is your content:")
	return data, contact
}

// terminates the node
func (cli *Cli) Exit(args []string) {
	cli.isRunning = false
}

func (cli *Cli) Help() {
	fmt.Println("- store data: put <file content>")
	fmt.Println("- retrieve data: get <file hash>")
	fmt.Println("- ping the current node: ping")
	fmt.Println("- exit: exit")
}

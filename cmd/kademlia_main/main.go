// TODO: Add package documentation for `main`, like this:
// Package main something something...
package main

import (
	"fmt"
	kademlia "kadlab-group-6/pkg/kademlia_node"
	"os"
	"strconv"
	"sync"

	getopt "github.com/pborman/getopt/v2"
)

var (
	// Bootstrap node configuration
	isBootstrapNode, _   = strconv.ParseBool(os.Getenv("IS_BOOTSTRAP_NODE"))
	BootstrapNodeAddress = os.Getenv("BOOTSTRAP_IP")
	BootstrapNodePort, _ = strconv.Atoi(os.Getenv("BOOTSTRAP_PORT"))
	BootstrapNodeId      = os.Getenv("BOOTSTRAP_ID")
)

func main() {

	fmt.Println("Bootstrap informations:", BootstrapNodeAddress, BootstrapNodeId)

	optCli := getopt.Bool('c', "Cli")
	getopt.Parse()

	// Create a new node
	fmt.Println("Creating a new node")
	var wg sync.WaitGroup
	wg.Add(2)

	// Create a bootstrap node and join the network
	if !isBootstrapNode {
		bootstrapNode := kademlia.NewContact(
			kademlia.NewKademliaID(BootstrapNodeId),
			BootstrapNodeAddress,
			BootstrapNodePort)
		node := kademlia.NewNode(kademlia.NewRandomKademliaID())
		go node.Network.Listen()
		go node.Join(bootstrapNode)

		if *optCli {
			c := kademlia.CliInit(node)
			go c.Main()
		}

		fmt.Println("Node id: ", node.Me.Id)
	} else {
		node := kademlia.NewNode(kademlia.NewKademliaID(BootstrapNodeId))
		go node.Network.Listen()
		fmt.Println("Node id: ", node.Me.Id)
	}
	wg.Wait() // Wait indefinitely
}

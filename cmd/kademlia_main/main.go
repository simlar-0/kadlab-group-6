// TODO: Add package documentation for `main`, like this:
// Package main something something...
package main

import (
	"fmt"
	kademlia "kadlab-group-6/pkg/kademlia_node"
	"os"
	"strconv"
)

var (
	// Bootstrap node configuration
	isBootstrapNode, _   = strconv.ParseBool(os.Getenv("IS_BOOTSTRAP_NODE"))
	BootstrapNodeAddress = os.Getenv("BOOTSTRAP_IP")
	BootstrapNodePort, _ = strconv.Atoi(os.Getenv("BOOTSTRAP_PORT"))
	BootstrapNodeId      = os.Getenv("BOOTSTRAP_ID")
)

func main() {
	// Create a new node
	fmt.Println("Creating a new node")
	var id *kademlia.KademliaID
	if !isBootstrapNode {
		id = kademlia.NewRandomKademliaID()
	} else {
		id = kademlia.NewKademliaID(BootstrapNodeId)
	}
	node := kademlia.NewNode(id)

	// Create a bootstrap node and join the network
	if !isBootstrapNode {
		fmt.Println("Creating a bootstrap node")
		bootstrapNode := kademlia.NewContact(
			kademlia.NewKademliaID(BootstrapNodeId),
			BootstrapNodeAddress,
			BootstrapNodePort)

		node.Join(bootstrapNode)
		fmt.Println("Joined the bootstrap network: ")
		fmt.Println(bootstrapNode)
	} else {
		fmt.Println("This node is a bootstrap node")
	}
	node.Network.Listen()
}

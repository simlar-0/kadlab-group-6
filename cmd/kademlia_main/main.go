// TODO: Add package documentation for `main`, like this:
// Package main something something...
package main

import (
	"fmt"
	kademlia "kadlab-group-6/pkg/kademlia_node"
)

func main() {
	// Create a new node
	fmt.Println("Creating a new node")
	node := kademlia.NewNode()
	node.Network.Listen()
}

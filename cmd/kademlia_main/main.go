// TODO: Add package documentation for `main`, like this:
// Package main something something...
package main

import (
	"fmt"
	node "kadlab-group-6/pkg/kademlia_node"
)

func main() {
	fmt.Println("Pretending to run the kademlia app...")
	// Using stuff from the kademlia package here. Something like...
	id := node.NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	contact := node.NewContact(id, "localhost:8000")
	fmt.Println(contact.String())
	fmt.Printf("%v\n", contact)
}

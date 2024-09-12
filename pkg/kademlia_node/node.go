package kademlia_node

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

type Node struct {
	// The contact of the node
	Me *Contact
	// The routing table of the node
	RoutingTable *RoutingTable
	// The network of the node
	Network *Network
	//
}

// NewNode returns a new instance of a Node
func NewNode() *Node {
	ip := GetLocalIp("eth0")
	port := GetPort()

	me := NewContact(NewRandomKademliaID(), ip, port)
	return &Node{
		Me:           me,
		RoutingTable: NewRoutingTable(me),
		Network:      InitNetwork(me)}
}

// Get random port between 1024 and 65535
func GetPort() int {
	source := rand.NewSource(time.Now().UnixNano())
	randomgen := rand.New(source)

	return randomgen.Intn(65535-1024) + 1024
}

// GetLocalIp returns the local ip address of the interface
// with the given name
func GetLocalIp(interfaceName string) string {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		fmt.Println("Error getting interface:", err)
		return "127.0.0.1"
	}

	addrs, err := iface.Addrs()
	if err != nil {
		fmt.Println("Error getting addresses:", err)
		return "127.0.0.1"
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}

	return "127.0.0.1"
}

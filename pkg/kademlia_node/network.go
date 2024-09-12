package kademlia_node

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type Network struct {
	Me *Contact
	// Channels for messages
	Outgoing chan RPC
	Incoming chan RPC
	Wg       sync.WaitGroup
}

func InitNetwork(me *Contact) *Network {
	return &Network{
		Outgoing: make(chan RPC),
		Incoming: make(chan RPC),
		Me:       me}
}

func (network *Network) Listen() {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", network.Me.Ip, network.Me.Port))
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}

	listener, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error starting UDP listener:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Listening on %s:%d\n", network.Me.Ip, network.Me.Port)

	// WaitGroup to keep the goroutines alive
	network.Wg.Add(3)
	// Start the goroutines for handling incoming and outgoing connections
	go network.ContinouslyReadUDP(listener)
	go network.ContinouslyWriteUDP(listener)
	go network.HandleIncomingChannel()
	network.Wg.Wait()
}

// ContinouslyReadUDP reads from the UDP connection
// and sends the message to the Incoming channel
func (network *Network) ContinouslyReadUDP(listener *net.UDPConn) {
	defer network.Wg.Done()
	buf := make([]byte, 1024)

	for {
		fmt.Println("Waiting for message")

		n, _, err := listener.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP connection:", err)
			continue
		}

		fmt.Println("Received message")

		rpc, err := network.DeserializeMessage(buf[:n])
		if err != nil {
			fmt.Println("Error deserializing message:", err)
			continue
		}

		// Send the message to the Incoming channel
		network.Incoming <- rpc
	}
}

// HandleIncomingChannel handles the incoming messages from the Incoming channel
// and sends the response to the Outgoing channel
func (network *Network) HandleIncomingChannel() {
	defer network.Wg.Done()

	for rpc := range network.Incoming {
		rpc, _ := network.HandelIncomingRPC(rpc)

		network.Outgoing <- rpc
	}
}

// ContinouslyWriteUDP writes to the UDP connection from the Outgoing channel
func (network *Network) ContinouslyWriteUDP(listener *net.UDPConn) {
	defer network.Wg.Done()

	for rpc := range network.Outgoing {
		// Serialize the message
		serializedMessage, err := network.SerializeMessage(rpc)
		if err != nil {
			fmt.Println("Error serializing message:", err)
			continue
		}

		// Get IP and port of the destination
		addrPort, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", rpc.Destination.Ip, rpc.Destination.Port))
		if err != nil {
			fmt.Println("Error parsing address and port:", err)
			continue
		}
		// Send the message
		_, err = listener.WriteToUDP(serializedMessage, addrPort)
		if err != nil {
			fmt.Println("Error writing to UDP connection:", err)
		}
	}
}

// HandelIncomingRPC handles the incoming RPC message and returns the response
func (network *Network) HandelIncomingRPC(rpc RPC) (RPC, error) {
	if !ValidateRPC(rpc) {
		return rpc, fmt.Errorf("invalid RPC")
	}

	rpc.Sender = network.Me
	// Switch on the type of the RPC
	switch rpc.Type {
	case PingRPC:
		//TODO: Implement this function
	case StoreRPC:
		//TODO: Implement this function
	case FindNodeRPC:
		//TODO: Implement this function
	case FindValueRPC:
		//TODO: Implement this function
	}
	return rpc, nil
}

// Serializes the RPC message to a byte array
func (network *Network) SerializeMessage(rpc RPC) (data []byte, err error) {
	return json.Marshal(rpc)
}

// Deserializes the byte array to an RPC message
func (network *Network) DeserializeMessage(data []byte) (RPC, error) {
	var rpc RPC
	err := json.Unmarshal(data, &rpc)
	return rpc, err
}

func (network *Network) SendPing(rpc RPC) {
}

func (network *Network) SendPingResponse(rpc RPC) {
}

func (network *Network) SendFindContact(rpc RPC) {
}

func (network *Network) SendFindContactResponse(rpc RPC) {
}

func (network *Network) SendFindData(rpc RPC) {
}

func (network *Network) SendFindDataResponse(rpc RPC) {
}

func (network *Network) SendStore(rpc RPC) {
}

func (network *Network) SendStoreResponse(rpc RPC) {
}

func (network *Network) SendMessage(rpc RPC) {
	network.Outgoing <- rpc
}

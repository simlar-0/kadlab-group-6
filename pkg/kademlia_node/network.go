package kademlia_node

// TODO:

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type Network struct {
	me *Contact
	// Channels for messages
	outgoing chan RPC
	incoming chan RPC
	wg       sync.WaitGroup
	// Map to keep track of sent RPCs
	sentRPCs map[string]*RPC
	mu       sync.Mutex
}

func InitNetwork(me *Contact) *Network {
	return &Network{
		outgoing: make(chan RPC),
		incoming: make(chan RPC),
		me:       me,
		sentRPCs: make(map[string]*RPC),
		mu:       sync.Mutex{}}
}

func (network *Network) Listen() {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", network.me.ip, network.me.port))
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

	fmt.Printf("Listening on %s:%d\n", network.me.ip, network.me.port)

	// WaitGroup to keep the goroutines alive
	network.wg.Add(3)
	// Start the goroutines for handling incoming and outgoing connections
	go network.ContinouslyReadUDP(listener)
	go network.ContinouslyWriteUDP(listener)
	go network.HandleIncomingChannel()
	network.wg.Wait()
}

// ContinouslyReadUDP reads from the UDP connection
// and sends the message to the Incoming channel
func (network *Network) ContinouslyReadUDP(listener *net.UDPConn) {
	defer network.wg.Done()
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
		network.incoming <- rpc
	}
}

// HandleIncomingChannel handles the incoming messages from the Incoming channel
// and sends the response to the Outgoing channel
func (network *Network) HandleIncomingChannel() {
	defer network.wg.Done()

	for rpc := range network.incoming {
		if rpc.IsResponse {
			network.mu.Lock()
			_, exists := network.sentRPCs[rpc.ID.String()]
			if exists {
				// Handle valid response
				delete(network.sentRPCs, rpc.ID.String())
			} else {
				// Handle invalid response
			}
			network.mu.Unlock()
		} else {
			// Handle request
			rpc, _ := network.HandelIncomingRPC(rpc)
			network.outgoing <- rpc
		}
	}
}

// ContinouslyWriteUDP writes to the UDP connection from the Outgoing channel
func (network *Network) ContinouslyWriteUDP(listener *net.UDPConn) {
	defer network.wg.Done()

	for rpc := range network.outgoing {
		// Serialize the message
		serializedMessage, err := network.SerializeMessage(rpc)
		if err != nil {
			fmt.Println("Error serializing message:", err)
			continue
		}

		// Get IP and port of the destination
		addrPort, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", rpc.Destination.ip, rpc.Destination.port))
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

	rpc.Sender = network.me
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
	network.mu.Lock()
	network.sentRPCs[rpc.ID.String()] = &rpc
	network.mu.Unlock()

	network.outgoing <- rpc
}

func (network *Network) SendResponse(rpc RPC) {
	network.outgoing <- rpc
}

package kademlia_node

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

// TODO move this to a config file or env variable
const (
	Timeout = 10 * time.Second
	Alpha   = 3
	Buffer  = 10
)

type Network struct {
	me       *Contact
	outgoing chan RPC
	wg       sync.WaitGroup
}

func NewNetwork(me *Contact) *Network {
	return &Network{
		me:       me,
		outgoing: make(chan RPC, Buffer),
		wg:       sync.WaitGroup{}}
}

// Listen starts a UDP listener on the specified IP and port of the network node.
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
	network.wg.Add(2)
	// Start the goroutines for handling incoming and outgoing connections
	go network.read(listener)

	// Create Alpha number of writer goroutines
	for i := 0; i < Alpha; i++ {
		go network.writer(listener)
	}
	network.wg.Wait()
}

// reads from the UDP connection and handles the incoming messages
func (network *Network) read(listener *net.UDPConn) {
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

		rpc, err := DeserializeMessage(buf[:n])
		if err != nil {
			fmt.Println("Error deserializing message:", err)
			continue
		}

		go RequestHandler(rpc, network)
	}
}

// Continuously writes to the UDP connection from the outgoing channel
func (network *Network) writer(listener *net.UDPConn) {
	defer network.wg.Done()

	for rpc := range network.outgoing {
		// Serialize the message
		serializedMessage, err := SerializeMessage(rpc)
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

func SerializeMessage(rpc RPC) (data []byte, err error) {
	return json.Marshal(rpc)
}

func DeserializeMessage(data []byte) (RPC, error) {
	var rpc RPC
	err := json.Unmarshal(data, &rpc)
	return rpc, err
}

func SendResponse(rpc RPC, network *Network) {
	network.outgoing <- rpc
}

func SendRequest(rpc RPC, network *Network) bool {
	network.outgoing <- rpc

	// TODO: Implement waiting for response and return false if timeout
	// somewhere else?
/* 	response, err := WaitForResponse(rpc, Timeout)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	fmt.Println("Received response:", response)
	*/
	return true 
}

func WaitForResponse(rpc RPC, timeout time.Duration) (RPC, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", rpc.Destination.ip, rpc.Destination.port))
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return RPC{}, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Error dialing UDP address:", err)
		return RPC{}, err
	}
	defer conn.Close()

	// Set read deadline
	conn.SetReadDeadline(time.Now().Add(timeout))

	buf := make([]byte, 1024)
	_, _, err = conn.ReadFromUDP(buf)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			fmt.Println("Timeout waiting for response")
		} else {
			fmt.Println("Error reading from UDP connection:", err)
		}
		return RPC{}, err
	}

	rpc, _ = DeserializeMessage(buf)
	return rpc, nil
}

// Get random port between 1024 and 65535
func GetRandomPort() int {
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

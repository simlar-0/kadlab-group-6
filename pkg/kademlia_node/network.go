package kademlia_node

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	Timeout         = 3 * time.Second // Timeout for waiting for a response
	Buffer          = 50              // Buffer size for the response queue
	NumberOfWorkers = 10              // Number of response workers
)

// NetworkInterface is an interface for sending and receiving messages
type NetworkInterface interface {
	SendRequest(rpc *RPC) (*RPC, error)
	SendResponse(rpc *RPC)
	Listen()
	Write(listener *net.UDPConn, serializedMessage []byte, addrPort *net.UDPAddr)
}

type Network struct {
	Node          *Node
	Wg            sync.WaitGroup
	ResponseQueue chan *RPC
	SentRequests  map[string]chan *RPC
	MutexRequest  sync.RWMutex
	MutexWrite    sync.RWMutex
}

var (
	networkInstance  *Network
	networkSingleton sync.Once
)

func NewNetwork(node *Node) *Network {
	networkSingleton.Do(func() {
		networkInstance = &Network{
			ResponseQueue: make(chan *RPC, Buffer),
			SentRequests:  make(map[string]chan *RPC),
			Wg:            sync.WaitGroup{},
			Node:          node}
	})
	return networkInstance
}

// Listen starts a UDP listener on the specified IP and port of the network node.
func (network *Network) Listen() {
	fmt.Println(network.Node.Me.Ip)
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", network.Node.Me.Ip, network.Node.Me.Port))
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

	fmt.Printf("Listening on %s:%d\n", network.Node.Me.Ip, network.Node.Me.Port)

	// Start the goroutines for handling incoming and outgoing connections
	network.Wg.Add(1)
	go network.read(listener)

	// Create Alpha number of response goroutines
	for i := 0; i < NumberOfWorkers; i++ {
		network.Wg.Add(1)
		go network.responseWorker(listener)
	}
	// WaitGroup to keep the goroutines alive
	network.Wg.Wait()
}

// reads from the UDP connection and handles the incoming messages
func (network *Network) read(listener *net.UDPConn) {
	defer network.Wg.Done()
	buf := make([]byte, 16384)

	for {
		fmt.Println("Waiting for message")

		n, _, err := listener.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP connection:", err)
			continue
		}

		rpc, err := network.Node.MessageHandler.DeserializeMessage(buf[:n])
		if err != nil {
			fmt.Println("Error deserializing message:", err)
			continue
		}
		fmt.Println("Received message:", rpc)

		// Check if the message is a response to a request
		reqID := rpc.ID.String()
		network.MutexRequest.RLock()
		recievedResponse, exists := network.SentRequests[reqID]
		network.MutexRequest.RUnlock()

		if !exists {
			// Create a goroutine to handle the incoming requests
			go network.Node.MessageHandler.ProcessRequest(rpc)
		} else {
			recievedResponse <- rpc
		}
	}
}

// responseWorker sends responses to the destination nodes
func (network *Network) responseWorker(listener *net.UDPConn) {
	defer network.Wg.Done()
	for rpc := range network.ResponseQueue {
		// Serialize the message
		serializedMessage, err := network.Node.MessageHandler.SerializeMessage(rpc)
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

		network.Write(listener, serializedMessage, addrPort)

		fmt.Println("Sent response with RPC ID: ", rpc.ID)
	}
}

// Write the response to the response channel
func (network *Network) Write(listener *net.UDPConn, serializedMessage []byte, addrPort *net.UDPAddr) {
	network.MutexWrite.Lock()
	defer network.MutexWrite.Unlock()
	// Send the message
	_, err := listener.WriteToUDP(serializedMessage, addrPort)
	if err != nil {
		fmt.Println("Error writing to UDP connection:", err)
	}
}

func (network *Network) SendResponse(rpc *RPC) {
	fmt.Println("Adding response to channel with RPC ID: ", rpc.ID)
	network.ResponseQueue <- rpc
}

// SendRequest sends an RPC to the destination node and waits for a response
// for a certain amount of time before timing out.
func (network *Network) SendRequest(rpc *RPC) (*RPC, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", rpc.Destination.Ip, rpc.Destination.Port))
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return &RPC{}, err
	}

	// Create a UDP connection to the destination
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Error dialing UDP address:", err)
		return &RPC{}, err
	}
	defer conn.Close()

	// Serialize the message
	serializedMessage, err := network.Node.MessageHandler.SerializeMessage(rpc)
	if err != nil {
		fmt.Println("Error serializing message:", err)
		return &RPC{}, err
	}

	// Create a unique request ID and response channel
	recievedResponse := make(chan *RPC)
	reqID := rpc.ID.String()
	network.MutexRequest.Lock()
	network.SentRequests[reqID] = recievedResponse
	network.MutexRequest.Unlock()

	// Defer cleanup
	defer func() {
		network.MutexRequest.Lock()
		delete(network.SentRequests, reqID)
		network.MutexRequest.Unlock()
		close(recievedResponse)
	}()

	// Send the message
	_, err = conn.Write(serializedMessage)
	if err != nil {
		fmt.Println("Error writing to UDP connection:", err)
		return &RPC{}, err
	}

	fmt.Println("Sent Request: ", rpc)

	fmt.Println("Waiting for response to RPC ID: ", rpc.ID)
	// Wait for response or timeout
	select {
	case response := <-recievedResponse:
		fmt.Println("Received Response: ", response)
		return response, nil
	case <-time.After(Timeout):
		fmt.Println("Timeout waiting for response to RPC ID: ", rpc.ID)
		return &RPC{}, fmt.Errorf("timeout waiting for response to RPC ID: %s", rpc.ID)
	}
}

// Get random port between 1024 and 65535
func GetRandomPortOrDefault() int {
	// If the node is a bootstrap node, use the port specified in the environment
	isBootstrapNode, _ := strconv.ParseBool(os.Getenv("IS_BOOTSTRAP_NODE"))
	if isBootstrapNode {
		portStr := os.Getenv("BOOTSTRAP_PORT")
		if portStr != "" {
			port, _ := strconv.Atoi(portStr)
			return port
		}
	}

	// If not, generate a random port
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

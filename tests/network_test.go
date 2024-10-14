package tests

import (
	kademlia "kadlab-group-6/pkg/kademlia_node"
	mocks "kadlab-group-6/pkg/mocks"
	"net"
	"os"
	"testing"
	"time"
)

func initNodeNetwork() *kademlia.Node {
	nodeID := kademlia.NewKademliaID("0000000000000000000000000000000000000000")
	me := kademlia.NewContact(nodeID, "127.0.0.1", 8000)
	node := &kademlia.Node{
		K:  20,
		Me: me,
	}
	node.RoutingTable = kademlia.NewRoutingTable(node)
	node.MessageHandler = mocks.NewMockMessageHandler(node)
	node.Network = kademlia.NewNetwork(node)

	return node
}

func TestNewNetwork(t *testing.T) {
	node := initNodeNetwork()

	if node.Network == nil {
		t.Fatal("Expected network instance, got nil")
	}

	// Ensure singleton pattern
	network2 := kademlia.NewNetwork(node)
	if node.Network != network2 {
		t.Error("Expected the same network instance, got different instances")
	}
}

func TestSendResponse(t *testing.T) {
	node := initNodeNetwork()
	network := kademlia.NewNetwork(node)

	rpc := &kademlia.RPC{ID: kademlia.NewRandomKademliaID()}
	node.Network.SendResponse(rpc)

	select {
	case res := <-network.ResponseQueue:
		if !res.ID.Equals(rpc.ID) {
			t.Errorf("Expected response ID %v, got %v", rpc.ID, res.ID)
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Expected response in queue, got timeout")
	}
}

func TestWrite(t *testing.T) {
	node := initNodeNetwork()

	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:8002")
	conn, _ := net.ListenUDP("udp", addr)
	defer conn.Close()

	message := []byte("test message")
	go node.Network.Write(conn, message, addr)

	buf := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		t.Fatalf("Expected to read from UDP, got error: %v", err)
	}

	if string(buf[:n]) != string(message) {
		t.Errorf("Expected message %s, got %s", message, buf[:n])
	}
}

func TestGetRandomPortOrDefault(t *testing.T) {
	port := kademlia.GetRandomPortOrDefault()
	if port < 1024 || port > 65535 {
		t.Errorf("Expected port between 1024 and 65535, got %d", port)
	}

	os.Setenv("IS_BOOTSTRAP_NODE", "true")
	os.Setenv("BOOTSTRAP_PORT", "8080")
	port = kademlia.GetRandomPortOrDefault()
	if port != 8080 {
		t.Errorf("Expected bootstrap port 8080, got %d", port)
	}
}

func TestResponseWorker(t *testing.T) {
	node := initNodeNetwork()
	network := kademlia.NewNetwork(node)

	// Create a mock UDP listener
	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	conn, _ := net.ListenUDP("udp", addr)
	defer conn.Close()

	// Create an RPC message and add it to the ResponseQueue
	rpc := &kademlia.RPC{
		ID:          kademlia.NewRandomKademliaID(),
		Destination: kademlia.NewContact(kademlia.NewRandomKademliaID(), "127.0.0.1", conn.LocalAddr().(*net.UDPAddr).Port),
	}
	network.ResponseQueue <- rpc

	// Run the responseWorker in a separate goroutine
	go network.ResponseWorker(conn)

	// Wait for the responseWorker to process the queue
	time.Sleep(1 * time.Second)

	// Verify the results
	buf := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		t.Fatalf("Expected to read from UDP, got error: %v", err)
	}

	// Deserialize the received message
	receivedRpc, err := network.Node.MessageHandler.DeserializeMessage(buf[:n])
	if err != nil {
		t.Fatalf("Expected to deserialize message, got error: %v", err)
	}

	if !receivedRpc.ID.Equals(rpc.ID) {
		t.Errorf("Expected RPC ID %v, got %v", rpc.ID, receivedRpc.ID)
	}
}

func TestListen(t *testing.T) {
	node := initNodeNetwork()
	go node.Network.Listen()

	time.Sleep(1 * time.Second) // Give some time for the listener to start

	// Check if the listener is running
	conn, err := net.Dial("udp", "127.0.0.1:8000")
	if err != nil {
		t.Fatalf("Expected to connect to listener, got error: %v", err)
	}
	conn.Close()
}

func TestSendRequest(t *testing.T) {
	node := initNodeNetwork()

	requestRpc := kademlia.NewRPC(kademlia.PingRequest, false, kademlia.NewRandomKademliaID(), nil, node.Me, node.Me)

	response, err := node.Network.SendRequest(requestRpc)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !response.ID.Equals(requestRpc.ID) {
		t.Errorf("Expected response ID %v, got %v", requestRpc.ID, response.ID)
	}
}

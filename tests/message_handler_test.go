package tests

import (
	"encoding/json"
	kademlia "kadlab-group-6/pkg/kademlia_node"
	mocks "kadlab-group-6/pkg/mocks"
	"testing"
)

func initNode() *kademlia.Node {
	nodeID := kademlia.NewKademliaID("0000000000000000000000000000000000000000")
	me := kademlia.NewContact(nodeID, "", 0)
	node := &kademlia.Node{
		K:  20,
		Me: me,
	}
	node.RoutingTable = kademlia.NewRoutingTable(node)
	node.MessageHandler = kademlia.NewMessageHandler(node)
	node.Network = mocks.NewMockNetwork(node)

	return node
}

func TestNewMessageHandler(t *testing.T) {
	node := initNode()
	handler := kademlia.NewMessageHandler(node)

	if handler.Node != node {
		t.Errorf("Expected Node %v, got %v", node, handler.Node)
	}
}

func TestProcessRequest(t *testing.T) {
	node := initNode()
	rpcID := kademlia.NewRandomKademliaID()
	source := kademlia.NewContact(kademlia.NewRandomKademliaID(), "", 0)
	destination := kademlia.NewContact(kademlia.NewRandomKademliaID(), "", 0)

	// Test ping request
	requestRPC := kademlia.NewRPC(kademlia.PingRequest, false, rpcID, nil, source, destination)
	expectedResponse := kademlia.NewRPC(kademlia.PingResponse, true, rpcID, nil, destination, source)
	responseRPC, err := node.MessageHandler.ProcessRequest(requestRPC)

	if expectedResponse.String() != responseRPC.String() {
		t.Errorf("Expected response %v, got %v", expectedResponse, responseRPC)
	}
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test Store request
	// TODO

	// Test FindNode request
	payload := kademlia.NewPayload(kademlia.NewRandomKademliaID(), nil, nil)
	requestRPC = kademlia.NewRPC(kademlia.FindNodeRequest, false, rpcID, payload, source, destination)
	contacts := node.RoutingTable.FindClosestContacts(requestRPC.Payload.Key)
	payload = kademlia.NewPayload(nil, nil, contacts)
	expectedResponse = kademlia.NewRPC(kademlia.FindNodeResponse, true, rpcID, payload, destination, source)
	responseRPC, err = node.MessageHandler.ProcessRequest(requestRPC)

	if expectedResponse.String() != responseRPC.String() {
		t.Errorf("Expected response %v, got %v", expectedResponse, responseRPC)
	}
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	// Test FindValue request
	// TODO

	// Test invalid RPC
	requestRPC = kademlia.NewRPC("INVALID_TYPE", false, rpcID, payload, source, destination)
	_, err = node.MessageHandler.ProcessRequest(requestRPC)
	expectedError := "invalid RPC"

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if err.Error() != expectedError {
		t.Errorf("Expected error message 'invalid RPC', got %v", err)
	}
}

func TestSerializeMessage(t *testing.T) {
	node := initNode()

	rpcID := kademlia.NewRandomKademliaID()
	source := kademlia.NewContact(kademlia.NewRandomKademliaID(), "", 0)
	destination := kademlia.NewContact(kademlia.NewRandomKademliaID(), "", 0)
	rpc := kademlia.NewRPC(kademlia.PingRequest, false, rpcID, nil, source, destination)

	data, err := node.MessageHandler.SerializeMessage(rpc)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	var deserializedRPC kademlia.RPC
	err = json.Unmarshal(data, &deserializedRPC)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if deserializedRPC.ID.String() != rpc.ID.String() {
		t.Errorf("Expected ID %s, got %s", rpc.ID.String(), deserializedRPC.ID.String())
	}
}

func TestDeserializeMessage(t *testing.T) {
	node := initNode()

	rpcID := kademlia.NewRandomKademliaID()
	source := kademlia.NewContact(kademlia.NewRandomKademliaID(), "", 0)
	destination := kademlia.NewContact(kademlia.NewRandomKademliaID(), "", 0)
	rpc := kademlia.NewRPC(kademlia.PingRequest, false, rpcID, nil, source, destination)

	data, err := json.Marshal(rpc)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	deserializedRPC, err := node.MessageHandler.DeserializeMessage(data)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if deserializedRPC.ID.String() != rpc.ID.String() {
		t.Errorf("Expected ID %s, got %s", rpc.ID.String(), deserializedRPC.ID.String())
	}
}

func TestSendPingRequest(t *testing.T) {
	node := initNode()

	source := kademlia.NewContact(kademlia.NewRandomKademliaID(), "", 0)
	destination := kademlia.NewContact(kademlia.NewRandomKademliaID(), "", 0)

	_, err := node.MessageHandler.SendPingRequest(source, destination)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSendFindNodeRequest(t *testing.T) {
	node := initNode()

	source := kademlia.NewContact(kademlia.NewRandomKademliaID(), "", 0)
	destination := kademlia.NewContact(kademlia.NewRandomKademliaID(), "", 0)
	key := kademlia.NewRandomKademliaID()

	_, err := node.MessageHandler.SendFindNodeRequest(source, destination, key)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSendPingResponse(t *testing.T) {
	node := initNode()

	rpcID := kademlia.NewRandomKademliaID()
	source := kademlia.NewContact(kademlia.NewRandomKademliaID(), "", 0)
	destination := kademlia.NewContact(kademlia.NewRandomKademliaID(), "", 0)
	requestRPC := kademlia.NewRPC(kademlia.PingRequest, false, rpcID, nil, source, destination)
	expectedRPC := kademlia.NewRPC(kademlia.PingResponse, true, rpcID, nil, destination, source)

	responseRPC := node.MessageHandler.SendPingResponse(requestRPC)
	if responseRPC.String() != expectedRPC.String() {
		t.Errorf("Expected RPC: %s, got %s", expectedRPC, requestRPC)
	}
}

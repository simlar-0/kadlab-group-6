package tests

import (
	//"encoding/json"

	kademlia "kadlab-group-6/pkg/kademlia_node"
	"testing"
)

func TestNewMessageHandler(t *testing.T) {
	n := &kademlia.Node{}
	handler := kademlia.NewMessageHandler(n)

	if handler.Node != n {
		t.Errorf("Expected Node %v, got %v", n, handler.Node)
	}
}

func TestProcessRequest(t *testing.T) {
	nodeID := kademlia.NewKademliaID("0000000000000000000000000000000000000000")
	me := kademlia.NewContact(nodeID, "", 0)
	node := &kademlia.Node{
		K:  20,
		Me: me,
	}
	node.RoutingTable = kademlia.NewRoutingTable(node)
	node.MessageHandler = kademlia.NewMessageHandler(node)
	node.Network = NewMockNetwork()
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

/*
func TestSerializeMessage(t *testing.T) {
	n := &kademlia.Node{}
	handler := kademlia.NewMessageHandler(n)

	id := kademlia.NewKademliaID("0000000000000000000000000000000000000000")
	source := kademlia.NewContact(kademlia.NewKademliaID("1000000000000000000000000000000000000000"), "1.2.3.4", 1234)
	destination := kademlia.NewContact(kademlia.NewKademliaID("2000000000000000000000000000000000000000"), "5.6.7.8", 5678)
	payload := kademlia.NewPayload(id, []byte("test data"), nil)
	rpc := kademlia.NewRPC(kademlia.PingRequest, false, id, payload, source, destination)

	data, err := handler.SerializeMessage(rpc)
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
	n := &kademlia.Node{}
	handler := kademlia.NewMessageHandler(n)

	id := kademlia.NewKademliaID("0000000000000000000000000000000000000000")
	source := kademlia.NewContact(kademlia.NewKademliaID("1000000000000000000000000000000000000000"), "1.2.3.4", 1234)
	destination := kademlia.NewContact(kademlia.NewKademliaID("2000000000000000000000000000000000000000"), "5.6.7.8", 5678)
	payload := kademlia.NewPayload(id, []byte("test data"), nil)
	rpc := kademlia.NewRPC(kademlia.PingRequest, false, id, payload, source, destination)

	data, err := handler.SerializeMessage(rpc)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	deserializedRPC, err := handler.DeserializeMessage(data)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if deserializedRPC.ID.String() != rpc.ID.String() {
		t.Errorf("Expected ID %s, got %s", rpc.ID.String(), deserializedRPC.ID.String())
	}
}

func TestSendPingRequest(t *testing.T) {
	n := &kademlia.Node{
		Network: &kademlia.Network{},
	}
	handler := kademlia.NewMessageHandler(n)

	source := kademlia.NewContact(kademlia.NewKademliaID("1000000000000000000000000000000000000000"), "1.2.3.4", 1234)
	destination := kademlia.NewContact(kademlia.NewKademliaID("2000000000000000000000000000000000000000"), "5.6.7.8", 5678)

	_, err := handler.SendPingRequest(source, destination)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSendPingResponse(t *testing.T) {
	n := &kademlia.Node{
		Network: &kademlia.Network{},
	}
	handler := kademlia.NewMessageHandler(n)

	id := kademlia.NewKademliaID("0000000000000000000000000000000000000000")
	source := kademlia.NewContact(kademlia.NewKademliaID("1000000000000000000000000000000000000000"), "1.2.3.4", 1234)
	destination := kademlia.NewContact(kademlia.NewKademliaID("2000000000000000000000000000000000000000"), "5.6.7.8", 5678)
	payload := kademlia.NewPayload(id, []byte("test data"), nil)
	requestRPC := kademlia.NewRPC(kademlia.PingRequest, false, id, payload, source, destination)

	responseRPC := handler.SendPingResponse(requestRPC)
	if responseRPC.Type != kademlia.PingResponse {
		t.Errorf("Expected PingResponse, got %s", responseRPC.Type)
	}
	if responseRPC.IsResponse != true {
		t.Errorf("Expected IsResponse true, got %t", responseRPC.IsResponse)
	}
	if responseRPC.Source.String() != destination.String() {
		t.Errorf("Expected Source %s, got %s", destination.String(), responseRPC.Source.String())
	}
	if responseRPC.Destination.String() != source.String() {
		t.Errorf("Expected Destination %s, got %s", source.String(), responseRPC.Destination.String())
	}
}
*/

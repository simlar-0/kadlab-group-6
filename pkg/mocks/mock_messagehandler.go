package mocks

import (
	"encoding/json"
	"fmt"
	kademlia "kadlab-group-6/pkg/kademlia_node"
)

type MockMessageHandler struct {
	Node *kademlia.Node
}

func NewMockMessageHandler(node *kademlia.Node) *MockMessageHandler {
	return &MockMessageHandler{Node: node}
}

func (handler *MockMessageHandler) ProcessRequest(rpc *kademlia.RPC) (*kademlia.RPC, error) {
	if !kademlia.ValidateRPC(rpc) || rpc.IsResponse {
		return nil, fmt.Errorf("invalid RPC")
	}

	fmt.Println("RPC: ", rpc)
	// Add the source to the routing table or update it
	handler.Node.RoutingTable.AddContact(rpc.Source)
	fmt.Println("Added contact to routing table")
	return nil, nil
}

func (handler *MockMessageHandler) SerializeMessage(rpc *kademlia.RPC) (data []byte, err error) {
	data, err = json.Marshal(rpc)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (handler *MockMessageHandler) DeserializeMessage(data []byte) (*kademlia.RPC, error) {
	var rpc kademlia.RPC
	err := json.Unmarshal(data, &rpc)
	if err != nil {
		return nil, err
	}
	return &rpc, nil
}

func (handler *MockMessageHandler) SendPingRequest(source *kademlia.Contact, destination *kademlia.Contact) (*kademlia.RPC, error) {
	return &kademlia.RPC{}, nil
}

func (handler *MockMessageHandler) SendPingResponse(requestRPC *kademlia.RPC) *kademlia.RPC {
	return nil
}

func (handler *MockMessageHandler) SendStoreRequest(source *kademlia.Contact, destination *kademlia.Contact, data []byte) (*kademlia.RPC, error) {
	return nil, nil
}

func (handler *MockMessageHandler) SendStoreResponse(requestRPC *kademlia.RPC) *kademlia.RPC {
	return nil
}

func (handler *MockMessageHandler) SendFindNodeRequest(source *kademlia.Contact, destination *kademlia.Contact, target *kademlia.KademliaID) (*kademlia.RPC, error) {
	return nil, nil
}

func (handler *MockMessageHandler) SendFindNodeResponse(requestRPC *kademlia.RPC) *kademlia.RPC {
	return nil
}

func (handler *MockMessageHandler) SendFindValueRequest(source *kademlia.Contact, destination *kademlia.Contact, key *kademlia.KademliaID) (*kademlia.RPC, error) {
	return nil, nil
}

func (handler *MockMessageHandler) SendFindValueResponse(requestRPC *kademlia.RPC) *kademlia.RPC {
	return nil
}

type MockMessageHandlerError struct {
	Node *kademlia.Node
}

func NewMockMessageHandlerError(node *kademlia.Node) *MockMessageHandlerError {
	return &MockMessageHandlerError{Node: node}
}

func (handler *MockMessageHandlerError) ProcessRequest(rpc *kademlia.RPC) (*kademlia.RPC, error) {
	return nil, fmt.Errorf("error")
}

func (handler *MockMessageHandlerError) SerializeMessage(rpc *kademlia.RPC) ([]byte, error) {
	return nil, fmt.Errorf("error")
}

func (handler *MockMessageHandlerError) DeserializeMessage(data []byte) (*kademlia.RPC, error) {
	return nil, fmt.Errorf("error")
}

func (handler *MockMessageHandlerError) SendPingRequest(source *kademlia.Contact, destination *kademlia.Contact) (*kademlia.RPC, error) {
	return nil, fmt.Errorf("error")
}

func (handler *MockMessageHandlerError) SendPingResponse(requestRPC *kademlia.RPC) *kademlia.RPC {
	return nil
}

func (handler *MockMessageHandlerError) SendStoreRequest(source *kademlia.Contact, destination *kademlia.Contact, data []byte) (*kademlia.RPC, error) {
	return nil, fmt.Errorf("error")
}

func (handler *MockMessageHandlerError) SendStoreResponse(requestRPC *kademlia.RPC) *kademlia.RPC {
	return nil
}

func (handler *MockMessageHandlerError) SendFindNodeRequest(source *kademlia.Contact, destination *kademlia.Contact, target *kademlia.KademliaID) (*kademlia.RPC, error) {
	return nil, fmt.Errorf("error")
}

func (handler *MockMessageHandlerError) SendFindNodeResponse(requestRPC *kademlia.RPC) *kademlia.RPC {
	return nil
}

func (handler *MockMessageHandlerError) SendFindValueRequest(source *kademlia.Contact, destination *kademlia.Contact, key *kademlia.KademliaID) (*kademlia.RPC, error) {
	return nil, fmt.Errorf("error")
}

func (handler *MockMessageHandlerError) SendFindValueResponse(requestRPC *kademlia.RPC) *kademlia.RPC {
	return nil
}

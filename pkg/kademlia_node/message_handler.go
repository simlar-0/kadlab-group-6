package kademlia_node

import (
	"encoding/json"
	"fmt"
)

type MessageHandler struct {
	Node *Node
	Network *Network
}

func NewMessageHandler(node *Node, network *Network) *MessageHandler {
	handler := &MessageHandler{
		Node: node,
		Network: network,}
	return handler
}

func (handler *MessageHandler) ProcessRequest(rpc *RPC) {
	if !ValidateRPC(rpc) || rpc.IsResponse {
		fmt.Errorf("invalid RPC")
		return
	}

	fmt.Println("RPC: ", rpc)
	// Add the source to the routing table or update it
	handler.Node.RoutingTable.AddContact(rpc.Source)
	fmt.Println("Added contact to routing table")

	switch rpc.Type {
	case PingRequest:
		fmt.Println("Received PingRequest")
		handler.SendPingResponse(rpc)
	case StoreRequest:
		// TODO: Store the data
		handler.SendStoreResponse(rpc)
	case FindNodeRequest:
		handler.Node.RoutingTable.FindClosestContacts(rpc.Payload.Key)
		handler.SendFindNodeResponse(rpc)
	case FindValueRequest:
		// TODO: Find the value
		handler.SendFindValueResponse(rpc)
	}
}

func (handler *MessageHandler) SerializeMessage(rpc *RPC) (data []byte, err error) {
	var serializedRPC []byte
	serializedRPC, err = json.Marshal(rpc)
	if err != nil {
		return nil, err
	}
	return serializedRPC, nil
}

func (handler *MessageHandler) DeserializeMessage(data []byte) (*RPC, error) {
	var rpc RPC
	err := json.Unmarshal(data, &rpc)
	if err != nil {
		return nil, err
	}
	return &rpc, nil
}

func (handler *MessageHandler) SendPingRequest(source *Contact, destination *Contact) (*RPC, error) {
	//TODO: implement
	RPC := NewRPC(PingRequest, false, NewRandomKademliaID(), nil, source, destination)
	response, err := handler.Node.Network.SendRequest(RPC)
	return response, err
}

func (handler *MessageHandler) SendPingResponse(requestRPC *RPC) *RPC {
	//TODO: implement
	rpc := NewRPC(PingResponse, true, requestRPC.ID, nil, requestRPC.Destination, requestRPC.Source)
	handler.Node.Network.SendResponse(rpc)
	return rpc
}

func (handler *MessageHandler) SendStoreRequest(source *Contact, destination *Contact, data []byte) (*RPC, error) {
	//TODO: implement
	// Create a new RPC for the StoreRequest
	RPC := NewRPC(StoreRequest, false, NewRandomKademliaID(), NewPayload(NewRandomKademliaID(), data, nil), source, destination)
	response, err := handler.Node.Network.SendRequest(RPC)
	return response, err
}

func (handler *MessageHandler) SendStoreResponse(requestRPC *RPC) *RPC {
	//TODO: implement
	// Create a new RPC for the StoreResponse
	responseRPC := NewRPC(StoreResponse, true, requestRPC.ID, nil, requestRPC.Destination, requestRPC.Source)
	
	// Send the response using the network
	handler.Network.SendResponse(responseRPC)
	
	// Return the response RPC
	return responseRPC
}

func (handler *MessageHandler) SendFindNodeRequest(source *Contact, destination *Contact, target *KademliaID) (*RPC, error) {
	RPC := NewRPC(FindNodeRequest, false, NewRandomKademliaID(), NewPayload(target, nil, nil), source, destination)
	response, err := handler.Node.Network.SendRequest(RPC)
	return response, err
}

func (handler *MessageHandler) SendFindNodeResponse(requestRPC *RPC) *RPC {
	// Get the k closest nodes to the target
	contacts := handler.Node.RoutingTable.FindClosestContacts(requestRPC.Payload.Key)
	rpc := NewRPC(FindNodeResponse, true, requestRPC.ID, NewPayload(nil, nil, contacts), requestRPC.Destination, requestRPC.Source)
	handler.Node.Network.SendResponse(rpc)
	return rpc
}

func (handler *MessageHandler) SendFindValueRequest(source *Contact, destination *Contact, key *KademliaID) (*RPC, error) {
	//TODO: implement
	RPC := NewRPC(FindValueRequest, false, NewRandomKademliaID(), nil, source, destination)
	response, err := handler.Node.Network.SendRequest(RPC)
	return response, err
}

func (handler *MessageHandler) SendFindValueResponse(requestRPC *RPC) *RPC {
	//TODO: implement
	// handler.Network.SendResponse(rpc)
	return NewRPC(FindValueResponse, true, requestRPC.ID, NewPayload(nil, requestRPC.Payload.Data, requestRPC.Payload.Contacts), requestRPC.Destination, requestRPC.Source)
}

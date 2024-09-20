package kademlia_node

import (
	"encoding/json"
	"fmt"
)

type MessageHandler struct {
	Node *Node
}

func NewMessageHandler(node *Node) *MessageHandler {
	handler := &MessageHandler{Node: node}
	return handler
}

func (handler *MessageHandler) ProcessRequest(rpc *RPC) (*RPC, error) {
	if !ValidateRPC(rpc) || rpc.IsResponse {
		fmt.Errorf("invalid RPC")
		return nil, fmt.Errorf("invalid RPC")
	}

	fmt.Println("RPC: ", rpc)
	// Add the source to the routing table or update it
	handler.Node.RoutingTable.AddContact(rpc.Source)
	fmt.Println("Added contact to routing table")

	switch rpc.Type {
	case PingRequest:
		fmt.Println("Received PingRequest")
		rpc := handler.SendPingResponse(rpc)
		return rpc, nil
	case StoreRequest:
		// TODO: Store the data
		rpc := handler.SendStoreResponse(rpc)
		return rpc, nil
	case FindNodeRequest:
		rpc := handler.SendFindNodeResponse(rpc)
		return rpc, nil
	case FindValueRequest:
		// TODO: Find the value
		rpc := handler.SendFindValueResponse(rpc)
		return rpc, nil
	}
	return nil, fmt.Errorf("invalid RPC")
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
	RPC := newRPC(PingRequest, false, NewRandomKademliaID(), nil, source, destination)
	response, err := handler.Node.Network.SendRequest(RPC)
	return response, err
}

func (handler *MessageHandler) SendPingResponse(requestRPC *RPC) *RPC {
	//TODO: implement
	rpc := newRPC(PingResponse, true, requestRPC.ID, nil, requestRPC.Destination, requestRPC.Source)
	handler.Node.Network.SendResponse(rpc)
	return rpc
}

func (handler *MessageHandler) SendStoreRequest(source *Contact, destination *Contact, data []byte) (*RPC, error) {
	//TODO: implement
	RPC := newRPC(StoreRequest, false, NewRandomKademliaID(), newPayload(NewRandomKademliaID(), data, nil), source, destination)
	response, err := handler.Node.Network.SendRequest(RPC)
	return response, err
}

func (handler *MessageHandler) SendStoreResponse(requestRPC *RPC) *RPC {
	//TODO: implement
	// handler.Network.SendResponse(rpc)
	return newRPC(StoreResponse, true, requestRPC.ID, nil, requestRPC.Destination, requestRPC.Source)
}

func (handler *MessageHandler) SendFindNodeRequest(source *Contact, destination *Contact, target *KademliaID) (*RPC, error) {
	RPC := newRPC(FindNodeRequest, false, NewRandomKademliaID(), newPayload(target, nil, nil), source, destination)
	response, err := handler.Node.Network.SendRequest(RPC)
	return response, err
}

func (handler *MessageHandler) SendFindNodeResponse(requestRPC *RPC) *RPC {
	// Get the k closest nodes to the target
	contacts := handler.Node.RoutingTable.FindClosestContacts(requestRPC.Payload.Key)
	rpc := newRPC(FindNodeResponse, true, requestRPC.ID, newPayload(nil, nil, contacts), requestRPC.Destination, requestRPC.Source)
	handler.Node.Network.SendResponse(rpc)
	return rpc
}

func (handler *MessageHandler) SendFindValueRequest(source *Contact, destination *Contact, key *KademliaID) (*RPC, error) {
	//TODO: implement
	RPC := newRPC(FindValueRequest, false, NewRandomKademliaID(), nil, source, destination)
	response, err := handler.Node.Network.SendRequest(RPC)
	return response, err
}

func (handler *MessageHandler) SendFindValueResponse(requestRPC *RPC) *RPC {
	//TODO: implement
	// handler.Network.SendResponse(rpc)
	return newRPC(FindValueResponse, true, requestRPC.ID, newPayload(nil, requestRPC.Payload.Data, requestRPC.Payload.Contacts), requestRPC.Destination, requestRPC.Source)
}

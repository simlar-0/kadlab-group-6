package kademlia_node

import (
	"encoding/json"
	"fmt"
)

type MessageHandlerInterface interface {
	ProcessRequest(rpc *RPC) (*RPC, error)
	SerializeMessage(rpc *RPC) ([]byte, error)
	DeserializeMessage(data []byte) (*RPC, error)
	SendPingRequest(source *Contact, destination *Contact) (*RPC, error)
	SendPingResponse(requestRPC *RPC) *RPC
	SendStoreRequest(source *Contact, destination *Contact, data []byte) (*RPC, error)
	SendStoreResponse(requestRPC *RPC) *RPC
	SendFindNodeRequest(source *Contact, destination *Contact, target *KademliaID) (*RPC, error)
	SendFindNodeResponse(requestRPC *RPC) *RPC
	SendFindValueRequest(source *Contact, destination *Contact, key *KademliaID) (*RPC, error)
	SendFindValueResponse(requestRPC *RPC) *RPC
}

type MessageHandler struct {
	Node *Node
}

func NewMessageHandler(node *Node) *MessageHandler {
	handler := &MessageHandler{Node: node}
	return handler
}

func (handler *MessageHandler) ProcessRequest(rpc *RPC) (*RPC, error) {
	if !ValidateRPC(rpc) || rpc.IsResponse {
		return nil, fmt.Errorf("invalid RPC")
	}

	fmt.Println("RPC: ", rpc)
	// Add the source to the routing table or update it
	handler.Node.RoutingTable.AddContact(rpc.Source)
	fmt.Println("Added contact to routing table")

	switch rpc.Type {
	case PingRequest:
		rpc := handler.SendPingResponse(rpc)
		return rpc, nil
	case StoreRequest:
		handler.Node.StoreData(rpc.Payload.Data)
		rpc := handler.SendStoreResponse(rpc)
		return rpc, nil
	case FindNodeRequest:
		rpc := handler.SendFindNodeResponse(rpc)
		return rpc, nil
	case FindValueRequest:
		rpc := handler.SendFindValueResponse(rpc)
		return rpc, nil
	default:
		return nil, fmt.Errorf("invalid RPC")
	}
}

func (handler *MessageHandler) SerializeMessage(rpc *RPC) (data []byte, err error) {
	data, err = json.Marshal(rpc)
	if err != nil {
		return nil, err
	}
	return data, nil
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
	rpc := NewRPC(PingRequest, false, NewRandomKademliaID(), nil, source, destination)
	response, err := handler.Node.Network.SendRequest(rpc)
	return response, err
}

func (handler *MessageHandler) SendPingResponse(requestRPC *RPC) *RPC {
	rpc := NewRPC(PingResponse, true, requestRPC.ID, nil, requestRPC.Destination, requestRPC.Source)
	handler.Node.Network.SendResponse(rpc)
	return rpc
}

func (handler *MessageHandler) SendStoreRequest(source *Contact, destination *Contact, data []byte) (*RPC, error) {
	rpc := NewRPC(StoreRequest, false, NewRandomKademliaID(), NewPayload(NewRandomKademliaID(), data, nil), source, destination)
	response, err := handler.Node.Network.SendRequest(rpc)
	return response, err
}

func (handler *MessageHandler) SendStoreResponse(requestRPC *RPC) *RPC {
	rpc := NewRPC(StoreResponse, true, requestRPC.ID, nil, requestRPC.Destination, requestRPC.Source)
	handler.Node.Network.SendResponse(rpc)
	return rpc
}

func (handler *MessageHandler) SendFindNodeRequest(source *Contact, destination *Contact, target *KademliaID) (*RPC, error) {
	rpc := NewRPC(FindNodeRequest, false, NewRandomKademliaID(), NewPayload(target, nil, nil), source, destination)
	response, err := handler.Node.Network.SendRequest(rpc)
	return response, err
}

func (handler *MessageHandler) SendFindNodeResponse(requestRPC *RPC) *RPC {
	contacts := handler.Node.RoutingTable.FindClosestContacts(requestRPC.Payload.Key)
	rpc := NewRPC(FindNodeResponse, true, requestRPC.ID, NewPayload(nil, nil, contacts), requestRPC.Destination, requestRPC.Source)
	handler.Node.Network.SendResponse(rpc)
	return rpc
}

func (handler *MessageHandler) SendFindValueRequest(source *Contact, destination *Contact, key *KademliaID) (*RPC, error) {
	fmt.Print(key)
	rpc := NewRPC(FindValueRequest, false, NewRandomKademliaID(), NewPayload(key, nil, nil), source, destination)
	response, err := handler.Node.Network.SendRequest(rpc)
	return response, err
}

func (handler *MessageHandler) SendFindValueResponse(requestRPC *RPC) *RPC {
	data, _ := handler.Node.GetData(requestRPC.Payload.Key)
	rpc := NewRPC(FindValueResponse, true, requestRPC.ID, NewPayload(nil, data, nil), requestRPC.Destination, requestRPC.Source)
	handler.Node.Network.SendResponse(rpc)
	return rpc
}

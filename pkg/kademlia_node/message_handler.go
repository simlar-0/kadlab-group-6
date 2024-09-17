package kademlia_node

import (
	"encoding/json"
	"fmt"
)

type MessageHandler struct {
	Network *Network
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{}
}

func (handler *MessageHandler) ProcessRequest(rpc *RPC) {
	if !ValidateRPC(rpc) || rpc.IsResponse {
		fmt.Errorf("invalid RPC")
		return
	}

	// TODO: Update the RoutingTable
	switch rpc.Type {
	case PingRequest:
		// TODO: Update the RoutingTable
		handler.SendPingResponse(rpc)
	case StoreRequest:
		// TODO: Store the data
		handler.SendStoreResponse(rpc)
	case FindNodeRequest:
		// TODO: Find the closest nodes
		handler.SendFindNodeResponse(rpc)
	case FindValueRequest:
		// TODO: Find the value
		handler.SendFindValueResponse(rpc)
	}
}

func (handler *MessageHandler) SerializeMessage(rpc *RPC) (data []byte, err error) {
	return json.Marshal(rpc)
}

func (handler *MessageHandler) DeserializeMessage(data []byte) (*RPC, error) {
	var rpc *RPC
	err := json.Unmarshal(data, rpc)
	return rpc, err
}

func (handler *MessageHandler) SendPingRequest(sender *Contact, destination *Contact) (*RPC, error) {
    RPC := newRPC(PingRequest, false, NewRandomKademliaID(), nil, sender, destination)
    response, err := handler.Network.SendRequest(RPC)
    return response, err
}

func (handler *MessageHandler) SendPingResponse(rpc *RPC) *RPC {
	//TODO: implement
	return newRPC(PingResponse, true, rpc.ID, nil, rpc.Destination, rpc.Sender)
}

func (handler *MessageHandler) SendStoreRequest(sender *Contact, destination *Contact, data []byte) (*RPC, error) {
    RPC := newRPC(StoreRequest, false, NewRandomKademliaID(), newPayload(NewRandomKademliaID(),data,nil,), sender, destination)
    response, err := handler.Network.SendRequest(RPC)
    return response, err
}


func (handler *MessageHandler) SendStoreResponse(rpc *RPC) *RPC {
	//TODO: implement
	return newRPC(StoreResponse, true, rpc.ID, nil, rpc.Destination, rpc.Sender)
}

func (handler *MessageHandler) SendFindNodeRequest(sender *Contact, destination *Contact, key *KademliaID) (*RPC, error) {
	RPC := newRPC(FindNodeRequest, false, NewRandomKademliaID(), nil, sender, destination)
	response, err := handler.Network.SendRequest(RPC)
	return response, err
}

func (handler *MessageHandler) SendFindNodeResponse(rpc *RPC) *RPC {
	//TODO: implement
	return newRPC(FindNodeResponse, true, rpc.ID, newPayload(nil, nil, []Contact{}), rpc.Destination, rpc.Sender)
}

func (handler *MessageHandler) SendFindValueRequest(sender *Contact, destination *Contact, key *KademliaID) (*RPC, error) {
    RPC := newRPC(FindValueRequest, false, NewRandomKademliaID(), nil, sender, destination)
    response, err := handler.Network.SendRequest(RPC)
    return response, err
}

func (handler *MessageHandler) SendFindValueResponse(rpc *RPC) *RPC {
	//TODO: implement
	return newRPC(FindValueResponse, true, rpc.ID, newPayload(nil, rpc.Payload.Data, rpc.Payload.Contacts), rpc.Destination, rpc.Sender)
}

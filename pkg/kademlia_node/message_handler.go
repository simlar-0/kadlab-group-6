package kademlia_node

import (
	"encoding/json"
	"fmt"
)

type MessageHandler struct {
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{}
}

func (handler *MessageHandler) ProcessRequest(rpc *RPC, network *Network) {
	if !ValidateRPC(rpc) || rpc.IsResponse {
		fmt.Errorf("invalid RPC")
		return
	}

	// TODO: Update the RoutingTable
	switch rpc.Type {
	case PingRequest:
		// TODO: Update the RoutingTable
		network.SendResponse(handler.NewPingResponse(rpc))
	case StoreRequest:
		// TODO: Store the data
		network.SendResponse(handler.NewStoreResponse(rpc))
	case FindNodeRequest:
		// TODO: Find the closest nodes
		network.SendResponse(handler.NewFindNodeResponse(rpc))
	case FindValueRequest:
		// TODO: Find the value
		network.SendResponse(handler.NewFindValueResponse(rpc))
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

func (handler *MessageHandler) NewPingRequest(sender *Contact, destination *Contact) *RPC {
	//TODO: implement
	return newRPC(PingRequest, false, NewRandomKademliaID(), nil, sender, destination)
}

func (handler *MessageHandler) NewPingResponse(rpc *RPC) *RPC {
	//TODO: implement
	return newRPC(PingResponse, true, rpc.ID, nil, rpc.Destination, rpc.Sender)
}

func (handler *MessageHandler) NewStoreRequest(sender *Contact, destination *Contact, key *KademliaID, data []byte) *RPC {
	//TODO: implement
	return newRPC(StoreRequest, false, NewRandomKademliaID(), newPayload(key, data, nil), sender, destination)
}

func (handler *MessageHandler) NewStoreResponse(rpc *RPC) *RPC {
	//TODO: implement
	return newRPC(StoreResponse, true, rpc.ID, nil, rpc.Destination, rpc.Sender)
}

func (handler *MessageHandler) NewFindNodeRequest(sender *Contact, destination *Contact, key *KademliaID) *RPC {
	//TODO: implement
	return newRPC(FindNodeRequest, false, NewRandomKademliaID(), newPayload(key, nil, nil), sender, destination)
}

func (handler *MessageHandler) NewFindNodeResponse(rpc *RPC) *RPC {
	//TODO: implement
	return newRPC(FindNodeResponse, true, rpc.ID, newPayload(nil, nil, []Contact{}), rpc.Destination, rpc.Sender)
}

func (handler *MessageHandler) NewFindValueRequest(sender *Contact, destination *Contact, key *KademliaID) *RPC {
	//TODO: implement
	return newRPC(FindValueRequest, false, NewRandomKademliaID(), newPayload(key, nil, nil), sender, destination)
}

func (handler *MessageHandler) NewFindValueResponse(rpc *RPC) *RPC {
	//TODO: implement
	return newRPC(FindValueResponse, true, rpc.ID, newPayload(nil, rpc.Payload.Data, rpc.Payload.Contacts), rpc.Destination, rpc.Sender)
}

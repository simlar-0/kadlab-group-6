package kademlia_node

import (
	"fmt"
)

type MessageHandler struct {
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{}
}

func RequestHandler(rpc RPC, network *Network) {
	if !ValidateRPC(rpc) {
		fmt.Errorf("invalid RPC")
		return
	}
	// TODO: Update the RoutingTable
	rpc.Sender = network.me
	// Switch on the type of the RPC
	switch rpc.Type {
	case PingRequest:
		// TODO: Update the RoutingTable
		rpc = NewPingResponse(rpc)
	case StoreRequest:
		// TODO: Store the data
		rpc = NewStoreResponse(rpc)
	case FindNodeRequest:
		// TODO: Find the closest nodes 
		rpc = NewFindNodeResponse(rpc)
	case FindValueRequest:
		// TODO: Find the value
		rpc = NewFindValueResponse(rpc)
	}

	SendResponse(rpc, network)
}

func NewPingRequest(sender *Contact, destination *Contact) RPC {
	//TODO: implement
	return newRPC(PingRequest, false, NewRandomKademliaID(), nil, sender, destination)
}

func NewPingResponse(rpc RPC) RPC {
	//TODO: implement
	return newRPC(PingResponse, true, rpc.ID, nil, rpc.Destination, rpc.Sender)
}

func NewStoreRequest(sender *Contact, destination *Contact, key *KademliaID, data []byte) RPC {
	//TODO: implement
	return newRPC(StoreRequest, false, NewRandomKademliaID(), newPayload(key, data, nil), sender, destination)
}

func NewStoreResponse(rpc RPC) RPC {
	//TODO: implement
	return newRPC(StoreResponse, true, rpc.ID, nil, rpc.Destination, rpc.Sender)
}

func NewFindNodeRequest(sender *Contact, destination *Contact, key *KademliaID) RPC {
	//TODO: implement
	return newRPC(FindNodeRequest, false, NewRandomKademliaID(), newPayload(key, nil, nil), sender, destination)
}

func NewFindNodeResponse(rpc RPC) RPC {
	//TODO: implement
	return newRPC(FindNodeResponse, true, rpc.ID, newPayload(nil, nil, []Contact{}), rpc.Destination, rpc.Sender)
}

func NewFindValueRequest(sender *Contact, destination *Contact, key *KademliaID) RPC {
	//TODO: implement
	return newRPC(FindValueRequest, false, NewRandomKademliaID(), newPayload(key, nil, nil), sender, destination)
}

func NewFindValueResponse(rpc RPC) RPC {
	//TODO: implement
	return newRPC(FindValueResponse, true, rpc.ID, newPayload(nil, rpc.Payload.Data, rpc.Payload.Contacts), rpc.Destination, rpc.Sender)
}

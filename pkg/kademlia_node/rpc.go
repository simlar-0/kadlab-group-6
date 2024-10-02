package kademlia_node

import (
	"fmt"
)

type RPC struct {
	ID          *KademliaID `json:"ID"`
	Type        RPCType     `json:"Type"`
	IsResponse  bool        `json:"IsResponse"`
	Destination *Contact    `json:"Destination"`
	Source      *Contact    `json:"Source"`
	Payload     *Payload    `json:"Payload"`
}

type Payload struct {
	Key      *KademliaID
	Data     []byte
	Contacts []*Contact
}

type RPCType string

// Define constants for valid RPC types
const (
	PingRequest       RPCType = "PING_REQUEST"
	PingResponse      RPCType = "PING_RESPONSE"
	StoreRequest      RPCType = "STORE_REQUEST"
	StoreResponse     RPCType = "STORE_RESPONSE"
	FindNodeRequest   RPCType = "FIND_NODE_REQUEST"
	FindNodeResponse  RPCType = "FIND_NODE_RESPONSE"
	FindValueRequest  RPCType = "FIND_VALUE_REQUEST"
	FindValueResponse RPCType = "FIND_VALUE_RESPONSE"
)

func newPayload(Key *KademliaID, Data []byte, Contacts []*Contact) *Payload {
	return &Payload{Key: Key, Data: Data, Contacts: Contacts}
}

func newRPC(Type RPCType, IsResponse bool, ID *KademliaID, Payload *Payload, Source *Contact, Destination *Contact) *RPC {
	return &RPC{Type: Type, IsResponse: IsResponse, ID: ID, Payload: Payload, Source: Source, Destination: Destination}
}

func ValidateRPC(rpc *RPC) bool {
	// Check if the RPC type is valid
	switch rpc.Type {
	case PingRequest, StoreRequest, FindNodeRequest, FindValueRequest, PingResponse, StoreResponse, FindNodeResponse, FindValueResponse:
		return true
	default:
		return false
	}
}

// String returns the string representation of the RPC
func (rpc *RPC) String() string {
	return fmt.Sprintf(`RPC(ID: "%s", Type: "%s", IsResponse: "%t", Destination: "%s", Source: "%s", Payload: "%s")`, rpc.ID, rpc.Type, rpc.IsResponse, rpc.Destination, rpc.Source, rpc.Payload)
}

// String returns the string representation of the Payload
func (payload *Payload) String() string {
	return fmt.Sprintf(`Payload(Key: "%s", Data: "%s", Contacts: "%s")`, payload.Key, payload.Data, payload.Contacts)
}

package kademlia_node

type RPC struct {
	// Type of the RPC
	Type RPCType
	// ID of the RPC
	ID *KademliaID
	// The data of the RPC message
	Payload *Payload
	// Destination of the RPC message
	Destination *Contact
	// Sender of the RPC message
	Sender *Contact
}

type Payload struct {
	// Key of the RPC message
	Key *KademliaID
	// Data of the RPC message, used in Store
	Data []byte
	// List of contacts, used in FindNode and FindData
	Contacts []Contact
}

type RPCType string

// Define constants for valid RPC types
const (
	PingRPC      RPCType = "PING"
	StoreRPC     RPCType = "STORE"
	FindNodeRPC  RPCType = "FIND_NODE"
	FindValueRPC RPCType = "FIND_VALUE"
)

func newPayload(Key *KademliaID, Data []byte, Contacts []Contact) *Payload {
	return &Payload{Key: Key, Data: Data, Contacts: Contacts}
}

func newRPC(Type RPCType, ID *KademliaID, Payload *Payload, Sender *Contact, Destination *Contact) *RPC {
	return &RPC{Type: Type, ID: ID, Payload: Payload, Sender: Sender, Destination: Destination}
}

func ValidateRPC(rpc RPC) bool {
	// Check if the RPC type is valid
	switch rpc.Type {
	case PingRPC, StoreRPC, FindNodeRPC, FindValueRPC:
		return true
	default:
		return false
	}
}

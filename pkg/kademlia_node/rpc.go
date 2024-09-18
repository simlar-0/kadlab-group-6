package kademlia_node

type RPC struct {
	Type        RPCType `json:"type"`
	IsResponse  bool   `json:"isResponse"`
	ID          *KademliaID `json:"id"`
	Payload     *Payload `json:"payload"`
	Destination *Contact `json:"destination"`
	Sender      *Contact `json:"sender"`
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

func newRPC(Type RPCType, IsResponse bool, ID *KademliaID, Payload *Payload, Sender *Contact, Destination *Contact) *RPC {
	return &RPC{Type: Type, IsResponse: IsResponse, ID: ID, Payload: Payload, Sender: Sender, Destination: Destination}
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

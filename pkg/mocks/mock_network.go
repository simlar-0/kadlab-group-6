package mocks

import (
	kademlia "kadlab-group-6/pkg/kademlia_node"
	"net"
)

type MockNetwork struct {
	sentMessages []*kademlia.RPC
	Node         *kademlia.Node
}

func NewMockNetwork(node *kademlia.Node) *MockNetwork {
	return &MockNetwork{
		sentMessages: []*kademlia.RPC{},
		Node:         node,
	}
}

func (m *MockNetwork) SendRequest(rpc *kademlia.RPC) (*kademlia.RPC, error) {
	m.sentMessages = append(m.sentMessages, rpc)

	return rpc, nil
}

func (m *MockNetwork) SendResponse(rpc *kademlia.RPC) {
	m.sentMessages = append(m.sentMessages, rpc)
}

func (m *MockNetwork) GetSentMessages() []*kademlia.RPC {
	return m.sentMessages
}

func (m *MockNetwork) Listen() {
	// Do nothing
}

func (m *MockNetwork) Write(*net.UDPConn, []byte, *net.UDPAddr) {
	// Do nothing
}

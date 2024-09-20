package tests

import (
	kademlia "kadlab-group-6/pkg/kademlia_node"
)

type MockNetwork struct {
	sentMessages []*kademlia.RPC
}

func NewMockNetwork() *MockNetwork {
	return &MockNetwork{
		sentMessages: []*kademlia.RPC{},
	}
}

func (m *MockNetwork) SendRequest(rpc *kademlia.RPC) (*kademlia.RPC, error) {
	m.sentMessages = append(m.sentMessages, rpc)

	return nil, nil
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

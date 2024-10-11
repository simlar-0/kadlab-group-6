package mocks

import (
	kademlia "kadlab-group-6/pkg/kademlia_node"
)

type MockNode struct{}

func (node *MockNode) Store(data []byte) (key *kademlia.KademliaID, err error) {
	return kademlia.NewRandomKademliaID(), nil
}

func (node *MockNode) LookupData(hash string) (content []byte, source *MockNode, err error) {
	return nil, nil, nil
}

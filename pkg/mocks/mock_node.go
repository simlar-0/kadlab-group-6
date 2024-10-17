package mocks

import (
	kademlia "kadlab-group-6/pkg/kademlia_node"
)

type MockNode struct{}

func (node *MockNode) Store(data []byte) (key *kademlia.KademliaID, err error) {
	return kademlia.NewKademliaID("0000000000000000000000000000000000000042"), nil
}

func (node *MockNode) LookupData(hash string) (content []byte, source *kademlia.Contact, err error) {
	if hash == "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b" {
		c := &kademlia.Contact{
			Id:   kademlia.NewKademliaID("0000000000000000000000000000000000001111"),
			Ip:   "111.111.111.111",
			Port: 1111,
		}

		return []byte("test"), c, nil
	}
	return nil, nil, nil
}

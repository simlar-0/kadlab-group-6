package tests

import (
	kademlia "kadlab-group-6/pkg/kademlia_node"
	mocks "kadlab-group-6/pkg/mocks"
	"os"
	"testing"
)

func initTestNode() *kademlia.Node {
	nodeID := kademlia.NewKademliaID("0000000000000000000000000000000000000000")
	me := kademlia.NewContact(nodeID, "127.0.0.1", 8000)
	node := &kademlia.Node{
		K:     20,
		Me:    me,
		Alpha: 3,
	}
	node.Network = mocks.NewMockNetwork(node)
	node.MessageHandler = mocks.NewMockMessageHandler(node)
	node.RoutingTable = kademlia.NewRoutingTable(node)
	return node
}

func TestNewNode(t *testing.T) {
	nodeID := kademlia.NewKademliaID("0000000000000000000000000000000000000000")

	os.Setenv("K", "20")
	os.Setenv("ALPHA", "3")
	node := kademlia.NewNode(nodeID)

	if node.Me.Id != nodeID {
		t.Errorf("Expected node ID %v, got %v", nodeID, node.Me.Id)
	}
	if node.K <= 0 {
		t.Errorf("Expected K to be greater than 0, got %v", node.K)
	}
	if node.Alpha <= 0 {
		t.Errorf("Expected Alpha to be greater than 0, got %v", node.Alpha)
	}
}

func TestJoin(t *testing.T) {
	node := initTestNode()
	targetID := kademlia.NewKademliaID("0000000000000000000000000000000000000001")
	targetNode := kademlia.NewContact(targetID, "127.0.0.1", 8001)

	err := node.Join(targetNode)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check if target node is in the routing table
	contacts := node.RoutingTable.FindClosestContacts(node.Me.Id)
	found := false
	for _, contact := range contacts {
		if contact.Id.Equals(targetNode.Id) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected target node %v to be in routing table, but it was not found", targetNode)
	}
}

func TestLookupContact(t *testing.T) {
	node := initTestNode()

	contact1 := kademlia.NewContact(kademlia.NewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1", 8000)
	contact2 := kademlia.NewContact(kademlia.NewKademliaID("0000000000000000000000000000000000000002"), "127.0.0.1", 8000)
	contact3 := kademlia.NewContact(kademlia.NewKademliaID("0000000000000000000000000000000000000003"), "127.0.0.1", 8000)

	node.RoutingTable.AddContact(contact1)
	node.RoutingTable.AddContact(contact2)
	node.RoutingTable.AddContact(contact3)

	targetID := kademlia.NewKademliaID("0000000000000000000000000000000000000004")
	targetContact := kademlia.NewContact(targetID, "127.0.0.1", 8004)

	closestContacts := node.LookupContact(targetContact)

	if len(closestContacts) == 0 {
		t.Errorf("Expected to find closest contacts, but got none")
	}

	expectedContacts := []*kademlia.Contact{contact1, contact2, contact3}
	for _, expectedContact := range expectedContacts {
		found := false
		for _, contact := range closestContacts {
			if contact.Id.Equals(expectedContact.Id) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected contact %v to be in the closest contacts, but it was not found", expectedContact)
		}
	}

}

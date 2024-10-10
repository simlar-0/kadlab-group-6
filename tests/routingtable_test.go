package tests

import (
	"fmt"
	kademlia "kadlab-group-6/pkg/kademlia_node"
	mocks "kadlab-group-6/pkg/mocks"
	"testing"
)

func initNodeRT() *kademlia.Node {
	nodeID := kademlia.NewKademliaID("0000000000000000000000000000000000000000")
	me := kademlia.NewContact(nodeID, "", 0)
	node := &kademlia.Node{
		K:  20,
		Me: me,
	}
	node.RoutingTable = kademlia.NewRoutingTable(node)
	node.MessageHandler = kademlia.NewMessageHandler(node)
	node.Network = mocks.NewMockNetwork(node)

	return node
}

func TestNewRoutingTable(t *testing.T) {
	node := &kademlia.Node{
		Me: &kademlia.Contact{
			Id: kademlia.NewKademliaID("0000000000000000000000000000000000000001"),
		},
		K: 20,
	}
	rt := kademlia.NewRoutingTable(node)

	if rt == nil {
		t.Errorf("Expected routing table to be initialized, got nil")
	}
	if rt.Node != node {
		t.Errorf("Expected node to be %v, got %v", node, node.RoutingTable.Node)
	}
	if len(rt.Buckets) != kademlia.IDLength*8 {
		t.Errorf("Expected %d buckets, got %d", kademlia.IDLength*8, len(node.RoutingTable.Buckets))
	}
}

func TestAddContactRT(t *testing.T) {
	node := initNodeRT()

	contact := &kademlia.Contact{
		Id: kademlia.NewKademliaID("0000000000000000000000000000000000000002"),
	}

	node.RoutingTable.AddContact(contact)

	bucketIndex := node.RoutingTable.GetBucketIndex(contact.Id)
	bucket := node.RoutingTable.Buckets[bucketIndex]

	contacts := bucket.GetContactsAndCalcDistance(node.Me.Id)
	found := false
	for _, c := range contacts {
		if c.Id.Equals(contact.Id) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected contact %v to be in bucket, but it was not found", contact)
	}
}

// TestAddContactToFullBucketPingOK tests adding a contact to a full bucket
// and when full the least recently seen contact is pinged. The ping is successful.
// The least recently seen contact should be moved to the front of the bucket.
// And the new contact should be ignored.
func TestAddContactToFullBucketPingOK(t *testing.T) {
	node := initNodeRT()
	node.MessageHandler = mocks.NewMockMessageHandler(node)

	counter := 0
	for i := 40; i < 61; i++ {
		// string formatting to get leading zeros
		s := fmt.Sprintf("%040d", i)
		fmt.Println(s)
		contact := &kademlia.Contact{
			Id: kademlia.NewKademliaID(s)}
		node.RoutingTable.AddContact(contact)
		counter++
	}

	t.Logf("Tried adding %d contacts to bucket 6", counter)
	node.RoutingTable.Buckets[6].PrintBucket()

	if node.RoutingTable.Buckets[6].Len() != 20 {
		t.Errorf("Expected bucket 6 to be full (20), got %d", node.RoutingTable.Buckets[6].Len())
	}

	// Expects the contacts with IDs 60 to not be in the bucket
	contact60 := kademlia.Contact{Id: kademlia.NewKademliaID("0000000000000000000000000000000000000060")}
	if node.RoutingTable.Buckets[6].Contains(contact60) {
		t.Errorf("Expected contact %v to not be in bucket", contact60)
	}

	// Expects the contacts with IDs 41 to the least recently seen contact
	contact41 := kademlia.Contact{Id: kademlia.NewKademliaID("0000000000000000000000000000000000000041")}
	if !node.RoutingTable.Buckets[6].GetLeastRecentlySeenContact().Id.Equals(contact41.Id) {
		t.Errorf("Expected contact %v to be least recently seen contact", contact41)
	}
}

// TestAddContactToFullBucketPingError tests adding a contact to a full bucket
// and when full the least recently seen contact is pinged. The ping is unsuccessful.
// The least recently seen contact should be removed from the bucket.
// And the new contact should be added to the bucket.
func TestAddContactToFullBucketPingError(t *testing.T) {
	node := initNodeRT()
	node.MessageHandler = mocks.NewMockMessageHandlerError(node)

	counter := 0
	for i := 40; i < 61; i++ {
		// string formatting to get leading zeros
		s := fmt.Sprintf("%040d", i)
		fmt.Println(s)
		contact := &kademlia.Contact{
			Id: kademlia.NewKademliaID(s)}
		node.RoutingTable.AddContact(contact)
		counter++
	}

	t.Logf("Tried adding %d contacts to bucket 6", counter)
	node.RoutingTable.Buckets[6].PrintBucket()

	if node.RoutingTable.Buckets[6].Len() != 20 {
		t.Errorf("Expected bucket 6 to be full (20), got %d", node.RoutingTable.Buckets[6].Len())
	}

	// Expects the contacts with IDs 60 to be in the bucket
	contact60 := kademlia.Contact{Id: kademlia.NewKademliaID("0000000000000000000000000000000000000060")}
	if !node.RoutingTable.Buckets[6].Contains(contact60) {
		t.Errorf("Expected contact %v to not be in bucket", contact60)
	}

	// Expects the contact with ID 40 to not be in the bucket
	contact40 := kademlia.Contact{Id: kademlia.NewKademliaID("0000000000000000000000000000000000000040")}
	if node.RoutingTable.Buckets[6].Contains(contact40) {
		t.Errorf("Expected contact %v to not be in bucket", contact40)
	}

	// Expects the contacts with IDs 41 to the least recently seen contact
	contact41 := kademlia.Contact{Id: kademlia.NewKademliaID("0000000000000000000000000000000000000041")}
	if !node.RoutingTable.Buckets[6].GetLeastRecentlySeenContact().Id.Equals(contact41.Id) {
		t.Errorf("Expected contact %v to be least recently seen contact", contact41)
	}
}

func TestFindClosestContacts(t *testing.T) {
	node := initNodeRT()

	contact1 := &kademlia.Contact{
		Id: kademlia.NewKademliaID("0000000000000000000000000000000000000002"),
	}
	contact2 := &kademlia.Contact{
		Id: kademlia.NewKademliaID("0000000000000000000000000000000000000003"),
	}

	node.RoutingTable.AddContact(contact1)
	node.RoutingTable.AddContact(contact2)

	target := kademlia.NewKademliaID("0000000000000000000000000000000000000004")
	closestContacts := node.RoutingTable.FindClosestContacts(target)

	if len(closestContacts) != 2 {
		t.Errorf("Expected 2 closest contacts, got %d", len(closestContacts))
	}
	if !closestContacts[0].Id.Equals(contact1.Id) {
		t.Errorf("Expected closest contact to be %v, got %v", contact1, closestContacts[0])
	}
	if !closestContacts[1].Id.Equals(contact2.Id) {
		t.Errorf("Expected second closest contact to be %v, got %v", contact2, closestContacts[1])
	}

}

func TestUpdateRoutingTable(t *testing.T) {
	node := initNodeRT()

	contact1 := &kademlia.Contact{
		Id: kademlia.NewKademliaID("0000000000000000000000000000000000000002"),
	}
	contact2 := &kademlia.Contact{
		Id: kademlia.NewKademliaID("0000000000000000000002000000000000000000"),
	}

	contacts := []*kademlia.Contact{contact1, contact2}
	node.RoutingTable.UpdateRoutingTable(contacts)

	bucketIndex1 := node.RoutingTable.GetBucketIndex(contact1.Id)
	bucket1 := node.RoutingTable.Buckets[bucketIndex1]
	bucketIndex2 := node.RoutingTable.GetBucketIndex(contact2.Id)
	bucket2 := node.RoutingTable.Buckets[bucketIndex2]

	contacts1 := bucket1.GetContactsAndCalcDistance(node.Me.Id)
	found1 := false
	for _, c := range contacts1 {
		if c.Id.Equals(contact1.Id) {
			found1 = true
			break
		}
	}
	if !found1 {
		t.Errorf("Expected contact1 %v to be in bucket, but it was not found", contact1)
	}

	contacts2 := bucket2.GetContactsAndCalcDistance(node.Me.Id)
	found2 := false
	for _, c := range contacts2 {
		if c.Id.Equals(contact2.Id) {
			found2 = true
			break
		}
	}
	if !found2 {
		t.Errorf("Expected contact2 %v to be in bucket, but it was not found", contact2)
	}
}

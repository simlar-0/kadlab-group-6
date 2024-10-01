package tests

import (
	kademlia "kadlab-group-6/pkg/kademlia_node"
	"testing"
)

func TestNewBucket(t *testing.T) {
	k := 20
	b := kademlia.NewBucket(k)
	if b == nil {
		t.Errorf("Expected new bucket to be created")
	}
	if b.K != k {
		t.Errorf("Expected bucket K value to be %d, got %d", k, b.K)
	}
	if b.List == nil {
		t.Errorf("Expected bucket List to be initialized")
	}
}

func TestAddContact(t *testing.T) {
	k := 20
	b := kademlia.NewBucket(k)
	contact := kademlia.NewContact(kademlia.NewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1", 8080)
	b.AddContact(*contact)
	if b.Len() != 1 {
		t.Errorf("Expected bucket length to be 1, got %d", b.Len())
	}
	b.AddContact(*contact)
	if b.Len() != 1 {
		t.Errorf("Expected bucket length to remain 1 after adding the same contact, got %d", b.Len())
	}
}

func TestRemoveContact(t *testing.T) {
	k := 20
	b := kademlia.NewBucket(k)
	contact := kademlia.NewContact(kademlia.NewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1", 8080)
	b.AddContact(*contact)
	b.RemoveContact(*contact)
	if b.Len() != 0 {
		t.Errorf("Expected bucket length to be 0 after removing the contact, got %d", b.Len())
	}
}

func TestGetContactAndCalcDistance(t *testing.T) {
	k := 20
	b := kademlia.NewBucket(k)
	target := kademlia.NewKademliaID("0000000000000000000000000000000000000000")
	contact := kademlia.NewContact(kademlia.NewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1", 8080)
	b.AddContact(*contact)
	contacts := b.GetContactAndCalcDistance(target)
	if len(contacts) != 1 {
		t.Errorf("Expected 1 contact, got %d", len(contacts))
	}
	if !contacts[0].Id.Equals(contact.Id) {
		t.Errorf("Expected contact ID to be %s, got %s", contact.Id.String(), contacts[0].Id.String())
	}
}

func TestLen(t *testing.T) {
	k := 20
	b := kademlia.NewBucket(k)
	if b.Len() != 0 {
		t.Errorf("Expected bucket length to be 0, got %d", b.Len())
	}
	contact := kademlia.NewContact(kademlia.NewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1", 8080)
	b.AddContact(*contact)
	if b.Len() != 1 {
		t.Errorf("Expected bucket length to be 1, got %d", b.Len())
	}
}

func TestGetLeastRecentlySeenContact(t *testing.T) {
	k := 20
	b := kademlia.NewBucket(k)
	contact1 := kademlia.NewContact(kademlia.NewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1", 8080)
	contact2 := kademlia.NewContact(kademlia.NewKademliaID("0000000000000000000000000000000000000002"), "127.0.0.1", 8081)
	b.AddContact(*contact1)
	b.AddContact(*contact2)
	leastRecentlySeen := b.GetLeastRecentlySeenContact()
	if !leastRecentlySeen.Id.Equals(contact1.Id) {
		t.Errorf("Expected least recently seen contact ID to be %s, got %s", contact1.Id.String(), leastRecentlySeen.Id.String())
	}
}

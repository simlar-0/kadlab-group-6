package tests

import (
	kademlia "kadlab-group-6/pkg/kademlia_node"
	"testing"
)

func TestNewContact(t *testing.T) {
	id := kademlia.NewKademliaID("0000000000000000000000000000000000000001")
	address := "127.0.0.1"
	port := 8080
	contact := kademlia.NewContact(id, address, port)

	if contact.Id.String() != id.String() {
		t.Errorf("Expected Id to be %s, got %s", id.String(), contact.Id.String())
	}
	if contact.Ip != address {
		t.Errorf("Expected Ip to be %s, got %s", address, contact.Ip)
	}
	if contact.Port != port {
		t.Errorf("Expected Port to be %d, got %d", port, contact.Port)
	}
	if contact.Distance != nil {
		t.Errorf("Expected Distance to be nil, got %s", contact.Distance.String())
	}
}

func TestCalcDistance(t *testing.T) {
	id1 := kademlia.NewKademliaID("0000000000000000000000000000000000000001")
	id2 := kademlia.NewKademliaID("0000000000000000000000000000000000000002")
	contact := kademlia.NewContact(id1, "127.0.0.1", 8080)
	contact.CalcDistance(id2)

	expectedDistance := kademlia.NewKademliaID("0000000000000000000000000000000000000003")
	if !contact.Distance.Equals(expectedDistance) {
		t.Errorf("Expected Distance to be %s, got %s", expectedDistance.String(), contact.Distance.String())
	}

	id1 = kademlia.NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	id2 = kademlia.NewKademliaID("F000000000000000000000000000000000000000")
	contact = kademlia.NewContact(id1, "127.0.0.1", 8080)
	contact.CalcDistance(id2)

	expectedDistance = kademlia.NewKademliaID("0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	if !contact.Distance.Equals(expectedDistance) {
		t.Errorf("Expected Distance to be %s, got %s", expectedDistance.String(), contact.Distance.String())
	}
}

func TestLess(t *testing.T) {
	id1 := kademlia.NewKademliaID("0000000000000000000000000000000000000001")
	id2 := kademlia.NewKademliaID("0000000000000000000000000000000000000002")
	id3 := kademlia.NewKademliaID("0000000000000000000000000000000000000003")

	contact1 := kademlia.NewContact(id1, "127.0.0.1", 8080)
	contact2 := kademlia.NewContact(id2, "127.0.0.1", 8081)
	contact1.CalcDistance(id3)
	contact2.CalcDistance(id3)

	if contact1.Less(contact2) != contact1.Distance.Less(contact2.Distance) {
		t.Errorf("Expected Less comparison to be %v, got %v", contact1.Distance.Less(contact2.Distance), contact1.Less(contact2))
	}
}

func TestString(t *testing.T) {
	id := kademlia.NewKademliaID("0000000000000000000000000000000000000001")
	contact := kademlia.NewContact(id, "127.0.0.1", 8080)
	expectedString := `contact(Id: "0000000000000000000000000000000000000001", Address: "127.0.0.1":"8080", Distance: "<nil>")`

	if contact.String() != expectedString {
		t.Errorf("Expected String to be %s, got %s", expectedString, contact.String())
	}
}

func TestContactCandidates_Append(t *testing.T) {
	id1 := kademlia.NewKademliaID("0000000000000000000000000000000000000001")
	id2 := kademlia.NewKademliaID("0000000000000000000000000000000000000002")
	contact1 := kademlia.NewContact(id1, "127.0.0.1", 8080)
	contact2 := kademlia.NewContact(id2, "127.0.0.1", 8081)

	candidates := &kademlia.ContactCandidates{}
	candidates.Append([]*kademlia.Contact{contact1, contact2})

	if len(candidates.Contacts) != 2 {
		t.Errorf("Expected length to be 2, got %d", len(candidates.Contacts))
	}
	if candidates.Contacts[0] != contact1 {
		t.Errorf("Expected first contact to be %v, got %v", &contact1, candidates.Contacts[0])
	}
	if candidates.Contacts[1] != contact2 {
		t.Errorf("Expected second contact to be %v, got %v", &contact2, candidates.Contacts[1])
	}
}

func TestContactCandidates_GetContacts(t *testing.T) {
	id1 := kademlia.NewKademliaID("0000000000000000000000000000000000000001")
	id2 := kademlia.NewKademliaID("0000000000000000000000000000000000000002")
	contact1 := kademlia.NewContact(id1, "127.0.0.1", 8080)
	contact2 := kademlia.NewContact(id2, "127.0.0.1", 8081)

	candidates := &kademlia.ContactCandidates{}
	candidates.Append([]*kademlia.Contact{contact1, contact2})

	contacts := candidates.GetContacts(1)
	if len(contacts) != 1 {
		t.Errorf("Expected length of contacts to be 1, got %d", len(contacts))
	}
	if contacts[0] != contact1 {
		t.Errorf("Expected first contact to be %v, got %v", &contact1, contacts[0])
	}

	contacts = candidates.GetContacts(2)
	if len(contacts) != 2 {
		t.Errorf("Expected length of contacts to be 2, got %d", len(contacts))
	}
	if contacts[0] != contact1 {
		t.Errorf("Expected first contact to be %v, got %v", &contact1, contacts[0])
	}

	contacts = candidates.GetContacts(3)
	if len(contacts) != 2 {
		t.Errorf("Expected length of contacts to be 2, got %d", len(contacts))
	}
}

func TestContactCandidates_Sort(t *testing.T) {
	id1 := kademlia.NewKademliaID("0000000000000000000000000000000000000001")
	id2 := kademlia.NewKademliaID("0000000000000000000000000000000000000002")
	id3 := kademlia.NewKademliaID("0000000000000000000000000000000000000003")
	contact1 := kademlia.NewContact(id1, "127.0.0.1", 8080)
	contact2 := kademlia.NewContact(id2, "127.0.0.1", 8081)
	contact3 := kademlia.NewContact(id3, "127.0.0.1", 8082)

	contact1.CalcDistance(id3)
	contact2.CalcDistance(id3)
	contact3.CalcDistance(id3)

	candidates := &kademlia.ContactCandidates{}
	candidates.Append([]*kademlia.Contact{contact3, contact1, contact2})
	candidates.Sort()

	if candidates.Contacts[0] != contact3 {
		t.Errorf("Expected first contact to be %v, got %v", contact3, candidates.Contacts[0])
	}
	if candidates.Contacts[1] != contact2 {
		t.Errorf("Expected second contact to be %v, got %v", contact2, candidates.Contacts[1])
	}
	if candidates.Contacts[2] != contact1 {
		t.Errorf("Expected third contact to be %v, got %v", contact1, candidates.Contacts[2])
	}
}

func TestContactCandidates_Len(t *testing.T) {
	id1 := kademlia.NewKademliaID("0000000000000000000000000000000000000001")
	id2 := kademlia.NewKademliaID("0000000000000000000000000000000000000002")
	contact1 := kademlia.NewContact(id1, "127.0.0.1", 8080)
	contact2 := kademlia.NewContact(id2, "127.0.0.1", 8081)

	candidates := &kademlia.ContactCandidates{}
	candidates.Append([]*kademlia.Contact{contact1, contact2})

	if candidates.Len() != 2 {
		t.Errorf("Expected length to be 2, got %d", candidates.Len())
	}
}

func TestContactCandidates_Swap(t *testing.T) {
	id1 := kademlia.NewKademliaID("0000000000000000000000000000000000000001")
	id2 := kademlia.NewKademliaID("0000000000000000000000000000000000000002")
	contact1 := kademlia.NewContact(id1, "127.0.0.1", 8080)
	contact2 := kademlia.NewContact(id2, "127.0.0.1", 8081)

	candidates := &kademlia.ContactCandidates{}
	candidates.Append([]*kademlia.Contact{contact1, contact2})
	candidates.Swap(0, 1)

	if candidates.Contacts[0] != contact2 {
		t.Errorf("Expected first contact to be %v, got %v", &contact2, candidates.Contacts[0])
	}
	if candidates.Contacts[1] != contact1 {
		t.Errorf("Expected second contact to be %v, got %v", &contact1, candidates.Contacts[1])
	}
}

func TestContactCandidates_Less(t *testing.T) {
	id1 := kademlia.NewKademliaID("0000000000000000000000000000000000000001")
	id2 := kademlia.NewKademliaID("0000000000000000000000000000000000000002")
	id3 := kademlia.NewKademliaID("0000000000000000000000000000000000000003")
	contact1 := kademlia.NewContact(id1, "127.0.0.1", 8080)
	contact2 := kademlia.NewContact(id2, "127.0.0.1", 8081)

	contact1.CalcDistance(id3)
	contact2.CalcDistance(id3)

	candidates := &kademlia.ContactCandidates{}
	candidates.Append([]*kademlia.Contact{contact1, contact2})
	if !candidates.Less(1, 0) {
		t.Errorf("Expected contact at index 0 to be less than contact at index 1")
	}
}

package tests

import (
	kademlia "kadlab-group-6/pkg/kademlia_node"
	"testing"
)

func TestNewShortlist(t *testing.T) {
	target := kademlia.NewKademliaID("0000000000000000000000000000000000000000")
	k := 3
	sl := kademlia.NewShortlist(target, k)

	if sl.Target != target {
		t.Errorf("Expected target %v, got %v", target, sl.Target)
	}
	if sl.K != k {
		t.Errorf("Expected K %d, got %d", k, sl.K)
	}
	if sl.Contacts.Len() != 0 {
		t.Errorf("Expected empty contacts list, got %d", sl.Contacts.Len())
	}
}

func TestAddContact(t *testing.T) {
	target := kademlia.NewKademliaID("0000000000000000000000000000000000000000")
	k := 3
	sl := kademlia.NewShortlist(target, k)

	contact := kademlia.NewContact(kademlia.NewKademliaID("1000000000000000000000000000000000000000"), "127.0.0.1", 8080)
	sl.AddContact(contact)

	if sl.Contacts.Len() != 1 {
		t.Errorf("Expected 1 contact, got %d", sl.Contacts.Len())
	}
	if !sl.Contains(contact) {
		t.Errorf("Expected contact to be in the shortlist")
	}
}

func TestAddContactListFull(t *testing.T) {
	target := kademlia.NewKademliaID("0000000000000000000000000000000000000000")
	k := 3
	sl := kademlia.NewShortlist(target, k)

	contact1 := kademlia.NewContact(kademlia.NewKademliaID("1000000000000000000000000000000000000000"), "127.0.0.1", 8080)
	contact2 := kademlia.NewContact(kademlia.NewKademliaID("2000000000000000000000000000000000000000"), "127.0.0.1", 8080)
	contact3 := kademlia.NewContact(kademlia.NewKademliaID("4000000000000000000000000000000000000000"), "127.0.0.1", 8080)
	contact4 := kademlia.NewContact(kademlia.NewKademliaID("3000000000000000000000000000000000000000"), "127.0.0.1", 8080)

	sl.AddContact(contact1)
	sl.AddContact(contact2)
	sl.AddContact(contact3)
	sl.AddContact(contact4)

	if sl.Contacts.Len() != k {
		t.Errorf("Expected %d contacts, got %d", k, sl.Contacts.Len())
	}
	if !sl.Contains(contact1) {
		t.Errorf("Expected contact1 to be in the shortlist")
	}
	if !sl.Contains(contact2) {
		t.Errorf("Expected contact2 to be in the shortlist")
	}
	if !sl.Contains(contact4) {
		t.Errorf("Expected contact3 to be in the shortlist")
	}
	if sl.Contains(contact3) {
		t.Errorf("Expected contact4 to not be in the shortlist")
	}
}

func TestRemoveContact(t *testing.T) {
	target := kademlia.NewKademliaID("0000000000000000000000000000000000000000")
	k := 3
	sl := kademlia.NewShortlist(target, k)

	contact1 := kademlia.NewContact(kademlia.NewKademliaID("1000000000000000000000000000000000000000"), "127.0.0.1", 8080)
	contact2 := kademlia.NewContact(kademlia.NewKademliaID("2000000000000000000000000000000000000000"), "127.0.0.1", 8080)
	contact3 := kademlia.NewContact(kademlia.NewKademliaID("3000000000000000000000000000000000000000"), "127.0.0.1", 8080)
	sl.AddContact(contact1)
	sl.AddContact(contact2)
	sl.AddContact(contact3)
	sl.RemoveContact(contact2)

	if sl.Contacts.Len() != k-1 {
		t.Errorf("Expected %d contacts, got %d", k, sl.Contacts.Len())
	}
	if !sl.Contains(contact1) {
		t.Errorf("Expected contact1 to be in the shortlist")
	}
	if sl.Contains(contact2) {
		t.Errorf("Expected contact2 to not be in the shortlist")
	}
	if !sl.Contains(contact3) {
		t.Errorf("Expected contact3 to be in the shortlist")
	}
}

func TestSort(t *testing.T) {
	target := kademlia.NewKademliaID("0000000000000000000000000000000000000000")
	k := 3
	sl := kademlia.NewShortlist(target, k)

	contact1 := kademlia.NewContact(kademlia.NewKademliaID("1000000000000000000000000000000000000000"), "127.0.0.1", 8080)
	contact2 := kademlia.NewContact(kademlia.NewKademliaID("2000000000000000000000000000000000000000"), "127.0.0.1", 8081)
	contact3 := kademlia.NewContact(kademlia.NewKademliaID("3000000000000000000000000000000000000000"), "127.0.0.1", 8082)

	sl.AddContact(contact3)
	sl.AddContact(contact1)
	sl.AddContact(contact2)

	sl.Sort()

	contacts := []*kademlia.Contact{}
	for elt := sl.Contacts.Front(); elt != nil; elt = elt.Next() {
		contacts = append(contacts, elt.Value.(*kademlia.Contact))
	}

	if contacts[0].Id.String() != contact1.Id.String() {
		t.Errorf("Expected contact1 to be first, got %v", contacts[0].Id)
	}
	if contacts[1].Id.String() != contact2.Id.String() {
		t.Errorf("Expected contact2 to be second, got %v", contacts[1].Id)
	}
	if contacts[2].Id.String() != contact3.Id.String() {
		t.Errorf("Expected contact3 to be third, got %v", contacts[2].Id)
	}
}

func TestGetClosestContacts(t *testing.T) {
	target := kademlia.NewKademliaID("0000000000000000000000000000000000000000")
	k := 3
	sl := kademlia.NewShortlist(target, k)

	contact1 := kademlia.NewContact(kademlia.NewKademliaID("1000000000000000000000000000000000000000"), "127.0.0.1", 8080)
	contact2 := kademlia.NewContact(kademlia.NewKademliaID("2000000000000000000000000000000000000000"), "127.0.0.1", 8081)
	contact3 := kademlia.NewContact(kademlia.NewKademliaID("3000000000000000000000000000000000000000"), "127.0.0.1", 8082)

	sl.AddContact(contact1)
	sl.AddContact(contact2)
	sl.AddContact(contact3)

	contacts := sl.GetClosestContacts(2)

	if len(contacts) != 2 {
		t.Errorf("Expected 2 contacts, got %d", len(contacts))
	}
	if contacts[0].Id.String() != contact1.Id.String() {
		t.Errorf("Expected contact1 to be first, got %v", contacts[0].Id)
	}
	if contacts[1].Id.String() != contact2.Id.String() {
		t.Errorf("Expected contact2 to be second, got %v", contacts[1].Id)
	}
}

func TestGetClosestContactsNotContacted(t *testing.T) {
	target := kademlia.NewKademliaID("0000000000000000000000000000000000000000")
	k := 3
	sl := kademlia.NewShortlist(target, k)

	contact1 := kademlia.NewContact(kademlia.NewKademliaID("1000000000000000000000000000000000000000"), "127.0.0.1", 8080)
	contact2 := kademlia.NewContact(kademlia.NewKademliaID("2000000000000000000000000000000000000000"), "127.0.0.1", 8081)
	contact3 := kademlia.NewContact(kademlia.NewKademliaID("3000000000000000000000000000000000000000"), "127.0.0.1", 8082)

	sl.AddContact(contact1)
	sl.AddContact(contact2)
	sl.AddContact(contact3)

	contacted := map[*kademlia.KademliaID]bool{
		contact1.Id: true,
		contact2.Id: true,
		contact3.Id: false,
	}

	contacts := sl.GetClosestContactsNotContacted(2, contacted)

	if len(contacts) != 1 {
		t.Errorf("Expected 1 contact, got %d", len(contacts))
	}
	if contacts[0].Id.String() != contact3.Id.String() {
		t.Errorf("Expected contact3 to be first, got %v", contacts[0].Id)
	}

	contacted[contact3.Id] = true

	contacts = sl.GetClosestContactsNotContacted(2, contacted)

	if len(contacts) != 0 {
		t.Errorf("Expected 0 contacts, got %d", len(contacts))
	}

	contacted[contact2.Id] = false
	contacted[contact3.Id] = false

	contacts = sl.GetClosestContactsNotContacted(2, contacted)

	if len(contacts) != 2 {
		t.Errorf("Expected 2 contacts, got %d", len(contacts))
	}
	if contacts[0].Id.String() != contact2.Id.String() {
		t.Errorf("Expected contact2 to be first, got %v", contacts[0].Id)
	}
	if contacts[1].Id.String() != contact3.Id.String() {
		t.Errorf("Expected contact3 to be second, got %v", contacts[1].Id)
	}

	contacted[contact1.Id] = false

	contacts = sl.GetClosestContactsNotContacted(2, contacted)

	if len(contacts) != 2 {
		t.Errorf("Expected 2 contacts, got %d", len(contacts))
	}
	if contacts[0].Id.String() != contact1.Id.String() {
		t.Errorf("Expected contact1 to be first, got %v", contacts[0].Id)
	}
	if contacts[1].Id.String() != contact2.Id.String() {
		t.Errorf("Expected contact2 to be second, got %v", contacts[1].Id)
	}

}

func TestGetClosestContact(t *testing.T) {
	target := kademlia.NewKademliaID("0000000000000000000000000000000000000000")
	k := 3
	sl := kademlia.NewShortlist(target, k)

	contact1 := kademlia.NewContact(kademlia.NewKademliaID("1000000000000000000000000000000000000000"), "127.0.0.1", 8080)
	contact2 := kademlia.NewContact(kademlia.NewKademliaID("2000000000000000000000000000000000000000"), "127.0.0.1", 8081)
	contact3 := kademlia.NewContact(kademlia.NewKademliaID("3000000000000000000000000000000000000000"), "127.0.0.1", 8082)

	sl.AddContact(contact1)
	sl.AddContact(contact2)
	sl.AddContact(contact3)

	contact := sl.GetClosestContact()

	if contact.Id.String() != contact1.Id.String() {
		t.Errorf("Expected contact1 to be closest, got %v", contact.Id)
	}

	sl.RemoveContact(contact1)

	contact = sl.GetClosestContact()

	if contact.Id.String() != contact2.Id.String() {
		t.Errorf("Expected contact2 to be closest, got %v", contact.Id)
	}

	sl.RemoveContact(contact2)
	sl.RemoveContact(contact3)

	contact = sl.GetClosestContact()

	if contact != nil {
		t.Errorf("Expected no contact to be closest, got %v", contact.Id)
	}
}

func TestAllContacted(t *testing.T) {
	target := kademlia.NewKademliaID("0000000000000000000000000000000000000000")
	k := 2
	sl := kademlia.NewShortlist(target, k)

	contact1 := kademlia.NewContact(kademlia.NewKademliaID("1000000000000000000000000000000000000000"), "127.0.0.1", 8080)
	contact2 := kademlia.NewContact(kademlia.NewKademliaID("2000000000000000000000000000000000000000"), "127.0.0.1", 8081)

	sl.AddContact(contact1)
	sl.AddContact(contact2)

	contacted := map[*kademlia.KademliaID]bool{
		contact1.Id: true,
		contact2.Id: true,
	}

	if !sl.AllContacted(contacted) {
		t.Errorf("Expected all contacts to be contacted")
	}

	contacted[contact2.Id] = false

	if sl.AllContacted(contacted) {
		t.Errorf("Expected not all contacts to be contacted")
	}
}

func TestContains(t *testing.T) {
	target := kademlia.NewKademliaID("0000000000000000000000000000000000000000")
	k := 3
	sl := kademlia.NewShortlist(target, k)

	contact := kademlia.NewContact(kademlia.NewKademliaID("1000000000000000000000000000000000000000"), "127.0.0.1", 8080)
	sl.AddContact(contact)

	if !sl.Contains(contact) {
		t.Errorf("Expected contact to be in the shortlist")
	}

	sl.RemoveContact(contact)

	if sl.Contains(contact) {
		t.Errorf("Expected contact to be removed from the shortlist")
	}
}

func TestString(t *testing.T) {
	target := kademlia.NewKademliaID("0000000000000000000000000000000000000000")
	k := 3
	sl := kademlia.NewShortlist(target, k)

	contact1 := kademlia.NewContact(kademlia.NewKademliaID("1000000000000000000000000000000000000000"), "127.0.0.1", 8080)
	contact2 := kademlia.NewContact(kademlia.NewKademliaID("2000000000000000000000000000000000000000"), "127.0.0.1", 8081)

	sl.AddContact(contact1)
	sl.AddContact(contact2)

	expected := "[" + contact1.String() + ", " + contact2.String() + "]"
	if sl.String() != expected {
		t.Errorf("Expected %v, got %v", expected, sl.String())
	}
}

func TestLen(t *testing.T) {
	target := kademlia.NewKademliaID("0000000000000000000000000000000000000000")
	k := 3
	sl := kademlia.NewShortlist(target, k)

	if sl.Len() != 0 {
		t.Errorf("Expected length 0, got %d", sl.Len())
	}

	contact := kademlia.NewContact(kademlia.NewKademliaID("1000000000000000000000000000000000000000"), "127.0.0.1", 8080)
	sl.AddContact(contact)

	length := sl.Len()

	if length != 1 {
		t.Errorf("Expected length 1, got %d", sl.Len())
	}
}

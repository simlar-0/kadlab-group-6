package kademlia_node

import (
	"container/list"
	"sort"
	"strings"
)

// shortlist represents a list of contacts sorted by distance to a target ID
type shortlist struct {
	Contacts *list.List
	Target   *KademliaID
	K        int
}

// NewShortlist creates a new shortlist
func NewShortlist(target *KademliaID, k int) *shortlist {
	return &shortlist{
		Contacts: list.New(),
		Target:   target,
		K:        k,
	}
}

// AddContact adds a contact to the shortlist
func (shortlist *shortlist) AddContact(contact *Contact) {
	// If the contact is already in the shortlist, skip
	if !shortlist.Contains(contact) {
		contact.CalcDistance(shortlist.Target)
		shortlist.Contacts.PushBack(contact)
		if shortlist.Contacts.Len() > shortlist.K {
			shortlist.Sort()
			shortlist.Contacts.Remove(shortlist.Contacts.Back())
		}
	}
}

// RemoveContact removes a contact from the shortlist
func (shortlist *shortlist) RemoveContact(contact *Contact) {
	for elt := shortlist.Contacts.Front(); elt != nil; elt = elt.Next() {
		if elt.Value.(*Contact).Id.Equals(contact.Id) {
			shortlist.Contacts.Remove(elt)
			break
		}
	}
}

// Sort sorts the contacts in the shortlist by distance to the target
func (shortlist *shortlist) Sort() {
	var contacts []*Contact
	for elt := shortlist.Contacts.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(*Contact)
		contacts = append(contacts, contact)
	}
	sort.Sort(byDistance(contacts))

	// Clear the list and reinsert sorted contacts
	shortlist.Contacts.Init()
	for _, contact := range contacts {
		shortlist.Contacts.PushBack(contact)
	}
}

// GetClosestContacts returns the closest contacts from the shortlist
func (shortlist *shortlist) GetClosestContactsNotContacted(k int, contacted map[*KademliaID]bool) []*Contact {
	var contacts []*Contact
	for elt := shortlist.Contacts.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(*Contact)
		if _, ok := contacted[contact.Id]; !ok {
			contacts = append(contacts, contact)
			if len(contacts) == k {
				break
			}
		}
	}
	return contacts
}

// GetClosestContacts returns the closest contacts from the shortlist
func (shortlist *shortlist) GetClosestContacts(k int) []*Contact {
	var contacts []*Contact
	for elt := shortlist.Contacts.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(*Contact)
		contacts = append(contacts, contact)
		if len(contacts) == k {
			break
		}
	}
	return contacts
}

// TrimToK trims the shortlist to length k
func (shortlist *shortlist) TrimToK() {
	for shortlist.Contacts.Len() > shortlist.K {
		shortlist.Contacts.Remove(shortlist.Contacts.Back())
	}
}

// CheckIfAllContacted checks each contact in the shortlist against a map of contacted contacts
// and is of length k
func (shortlist *shortlist) AllContacted(contacted map[*KademliaID]bool) bool {
	for elt := shortlist.Contacts.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(*Contact)
		if _, ok := contacted[contact.Id]; !ok {
			return false
		}
	}
	return true && shortlist.Contacts.Len() == shortlist.K
}

// CheckClosestContact checks if the closest contact in the shortlist is closer than the given contact
func (shortlist *shortlist) CheckClosestContact(contact *Contact) bool {
	if shortlist.Contacts.Len() == 0 {
		return true
	}
	return contact.Distance.Less(shortlist.Contacts.Front().Value.(*Contact).Distance)
}

// Contains checks if the shortlist contains a specific contact
func (shortlist *shortlist) Contains(contact *Contact) bool {
	for elt := shortlist.Contacts.Front(); elt != nil; elt = elt.Next() {
		if elt.Value.(*Contact).Id.Equals(contact.Id) {
			return true
		}
	}
	return false
}

// String returns a simple string representation of a shortlist
func (shortlist *shortlist) String() string {
	var contacts []string
	for elt := shortlist.Contacts.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(*Contact)
		contacts = append(contacts, contact.String())
	}
	return "[" + strings.Join(contacts, ", ") + "]"
}

// Len returns the length of the shortlist
func (shortlist *shortlist) Len() int {
	return shortlist.Contacts.Len()
}

// byDistance is a wrapper type for sorting contacts by distance
type byDistance []*Contact

func (bd byDistance) Len() int {
	return len(bd)
}

func (bd byDistance) Less(i, j int) bool {
	return bd[i].Distance.Less(bd[j].Distance)
}

func (bd byDistance) Swap(i, j int) {
	bd[i], bd[j] = bd[j], bd[i]
}

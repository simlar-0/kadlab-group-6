package kademlia_node

import (
	"fmt"
	"sort"
)

// Contact definition
// stores the KademliaID, the ip address and the distance
type Contact struct {
	Id       *KademliaID `json:"Id"`
	Ip       string      `json:"Ip"`
	Port     int         `json:"Port"`
	Distance *KademliaID `json:"Distance"`
}

func NewContact(id *KademliaID, address string, port int) *Contact {
	return &Contact{id, address, port, nil}
}

// CalcDistance calculates the distance to the target and
// fills the contacts distance field
func (contact *Contact) CalcDistance(target *KademliaID) {
	contact.Distance = contact.Id.CalcDistance(target)
}

// Less returns true if contact.distance < otherContact.distance
func (contact *Contact) Less(otherContact *Contact) bool {
	return contact.Distance.Less(otherContact.Distance)
}

func (contact *Contact) String() string {
	return fmt.Sprintf(`contact(Id: "%s", Address: "%s":"%d", Distance: "%s")`, contact.Id, contact.Ip, contact.Port, contact.Distance)
}

// ContactCandidates definition, stores an array of Contacts
type ContactCandidates struct {
	Contacts []*Contact
}

// Append an array of Contacts to the ContactCandidates
func (candidates *ContactCandidates) Append(contacts []*Contact) {
	candidates.Contacts = append(candidates.Contacts, contacts...)
}

// GetContacts returns the first count number of Contacts
func (candidates *ContactCandidates) GetContacts(count int) []*Contact {
	if count > len(candidates.Contacts) {
		count = len(candidates.Contacts)
	}
	return candidates.Contacts[:count]
}

// Sort the Contacts in ContactCandidates
func (candidates *ContactCandidates) Sort() {
	sort.Sort(candidates)
}

// Len returns the length of the ContactCandidates
func (candidates *ContactCandidates) Len() int {
	return len(candidates.Contacts)
}

// Swap the position of the Contacts at i and j
// WARNING does not check if either i or j is within range
func (candidates *ContactCandidates) Swap(i, j int) {
	candidates.Contacts[i], candidates.Contacts[j] = candidates.Contacts[j], candidates.Contacts[i]
}

// Less returns true if the Contact at index i is smaller than
// the Contact at index j
func (candidates *ContactCandidates) Less(i, j int) bool {
	return candidates.Contacts[i].Less(candidates.Contacts[j])
}

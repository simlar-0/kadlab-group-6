package kademlia_node

import (
	"container/list"
	"fmt"
)

// bucket definition
// contains a List
type bucket struct {
	List *list.List
	K    int
}

// NewBucket returns a new instance of a bucket
func NewBucket(k int) *bucket {
	bucket := &bucket{}
	bucket.List = list.New()
	bucket.K = k
	return bucket
}

// AddContact adds the Contact to the front of the bucket
// or moves it to the front of the bucket if it already existed
func (bucket *bucket) AddContact(contact Contact) {
	var element *list.Element
	for e := bucket.List.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).Id

		if (contact).Id.Equals(nodeID) {
			element = e
		}
	}

	if element == nil {
		if bucket.List.Len() < bucket.K {
			bucket.List.PushFront(contact)
		}
	} else {
		bucket.List.MoveToFront(element)
	}
}

// RemoveContact removes the Contact from the bucket
func (bucket *bucket) RemoveContact(contact Contact) {
	for elt := bucket.List.Front(); elt != nil; elt = elt.Next() {
		if elt.Value.(Contact).Id.Equals(contact.Id) {
			bucket.List.Remove(elt)
			break
		}
	}
}

// GetContactAndCalcDistance returns an array of Contacts where
// the distance has already been calculated
func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []*Contact {
	var contacts []*Contact

	for elt := bucket.List.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, &contact)
	}

	return contacts
}

// Len return the size of the bucket
func (bucket *bucket) Len() int {
	return bucket.List.Len()
}

// Returns the least recently seen contact in the bucket
func (bucket *bucket) GetLeastRecentlySeenContact() Contact {
	return bucket.List.Back().Value.(Contact)
}

// Contains returns true if the bucket contains the contact
func (bucket *bucket) Contains(contact Contact) bool {
	for elt := bucket.List.Front(); elt != nil; elt = elt.Next() {
		if elt.Value.(Contact).Id.Equals(contact.Id) {
			return true
		}
	}
	return false
}

// PrintBucket prints the bucket
func (bucket *bucket) PrintBucket() {
	for elt := bucket.List.Front(); elt != nil; elt = elt.Next() {
		fmt.Println("Bucket element: {Id: %s, Ip: %s, Port: %d}", elt.Value.(Contact).Id, elt.Value.(Contact).Ip, elt.Value.(Contact).Port)
	}
}

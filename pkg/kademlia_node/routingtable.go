package kademlia_node

import (
	"sync"
)

// RoutingTable definition
// keeps a refrence contact of me and an array of buckets
type RoutingTable struct {
	Me      *Contact
	Buckets [IDLength * 8]*bucket
	K       int
	Mu      sync.RWMutex
}

// NewRoutingTable returns a new instance of a RoutingTable
func NewRoutingTable(me *Contact, k int) *RoutingTable {
	routingTable := &RoutingTable{}
	for i := 0; i < IDLength*8; i++ {
		routingTable.Buckets[i] = newBucket(k)
	}
	routingTable.Me = me
	routingTable.K = k
	routingTable.Mu = sync.RWMutex{}
	return routingTable
}

// AddContact add a new contact to the correct Bucket
func (routingTable *RoutingTable) AddContact(contact *Contact) {
	routingTable.Mu.Lock()
	defer routingTable.Mu.Unlock()

	bucketIndex := routingTable.getBucketIndex(contact.Id)
	bucket := routingTable.Buckets[bucketIndex]
	bucket.AddContact(*contact)
}

// FindClosestContacts finds the count closest Contacts to the target in the RoutingTable
func (routingTable *RoutingTable) FindClosestContacts(target *KademliaID) []*Contact {
	routingTable.Mu.Lock()
	defer routingTable.Mu.Unlock()

	var candidates ContactCandidates
	bucketIndex := routingTable.getBucketIndex(target)
	bucket := routingTable.Buckets[bucketIndex]

	candidates.Append(bucket.GetContactAndCalcDistance(target))

	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < IDLength*8) && candidates.Len() < routingTable.K; i++ {
		if bucketIndex-i >= 0 {
			bucket = routingTable.Buckets[bucketIndex-i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
		if bucketIndex+i < IDLength*8 {
			bucket = routingTable.Buckets[bucketIndex+i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
	}

	candidates.Sort()

	return candidates.GetContacts(candidates.Len())
}

// getBucketIndex get the correct Bucket index for the KademliaID
func (routingTable *RoutingTable) getBucketIndex(id *KademliaID) int {
	distance := id.CalcDistance(routingTable.Me.Id)
	for i := 0; i < IDLength; i++ {
		for j := 0; j < 8; j++ {
			if (distance[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}

	return IDLength*8 - 1
}

func (routingTable *RoutingTable) UpdateRoutingTable(contacts []*Contact) {
	routingTable.Mu.Lock()
	defer routingTable.Mu.Unlock()

	// Get the bucket index
	// Get the bucket
	// Check if the bucket is full
	// If the bucket is full, ping the least recently seen contact
	// If the contact is still alive, move it to the front of the bucket
}

func (routingTable *RoutingTable) Refresh(KademliaID *KademliaID) {
	// Refresh the bucket containing the KademliaID
}

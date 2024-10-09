package kademlia_node

import (
	"fmt"
	"math/bits"
	"sync"
)

type RoutingTable struct {
	Node    *Node
	Buckets []*bucket
	Mutex   sync.RWMutex
}

// NewRoutingTable returns a new instance of a RoutingTable
func NewRoutingTable(node *Node) *RoutingTable {
	routingTable := &RoutingTable{
		Node: node}
	routingTable.Buckets = make([]*bucket, IDLength*8)
	for i := 0; i < IDLength*8; i++ {
		routingTable.Buckets[i] = newBucket(node.K)
	}
	return routingTable
}

// AddContact add a new contact to the correct Bucket
func (routingTable *RoutingTable) AddContact(contact *Contact) {
	routingTable.Mutex.Lock()
	defer routingTable.Mutex.Unlock()

	// Update the distance of the contact
	contact.CalcDistance(routingTable.Node.Me.Id)

	bucketIndex := routingTable.GetBucketIndex(contact.Id)

	bucket := routingTable.Buckets[bucketIndex]
	// Check if the bucket is full
	if bucket.Len() >= routingTable.Node.K {
		// Ping the least recently seen contact
		leastRecent := bucket.GetLeastRecentlySeenContact()
		_, err := routingTable.Node.MessageHandler.SendPingRequest(routingTable.Node.Me, &leastRecent)
		if err == nil {
			// If the contact is still alive, move it to the front of the bucket
			bucket.AddContact(leastRecent)
		} else {
			// If the contact is not alive, remove it from the bucket
			bucket.RemoveContact(leastRecent)
		}
	}
	bucket.AddContact(*contact)
}

// FindClosestContacts finds the count closest Contacts to the target in the RoutingTable
func (routingTable *RoutingTable) FindClosestContacts(target *KademliaID) []*Contact {
	routingTable.Mutex.Lock()
	defer routingTable.Mutex.Unlock()

	var candidates ContactCandidates
	bucketIndex := routingTable.GetBucketIndex(target)
	bucket := routingTable.Buckets[bucketIndex]

	candidates.Append(bucket.GetContactAndCalcDistance(target))

	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < IDLength*8) && candidates.Len() < routingTable.Node.K; i++ {
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

// GetBucketIndex get the correct Bucket index for the KademliaID
func (routingTable *RoutingTable) GetBucketIndex(id *KademliaID) int {
	distance := id.CalcDistance(routingTable.Node.Me.Id)
	leadingZeros := 0

	for i := 0; i < IDLength; i++ {
		if distance[i] == 0 {
			leadingZeros += 8
		} else {
			leadingZeros += bits.LeadingZeros8(distance[i])
			break
		}
	}

	return IDLength*8 - leadingZeros - 1
}

// UpdateRoutingTable updates the RoutingTable with a list of new contacts
func (routingTable *RoutingTable) UpdateRoutingTable(contacts []*Contact) {
	for _, contact := range contacts {
		routingTable.AddContact(contact)
	}
}

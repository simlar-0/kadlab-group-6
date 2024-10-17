package kademlia_node

import (
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
		routingTable.Buckets[i] = NewBucket(node.K)
	}
	return routingTable
}

// AddContact add a new contact to the correct Bucket
func (routingTable *RoutingTable) AddContact(contact *Contact) {
	routingTable.Mutex.Lock()
	defer routingTable.Mutex.Unlock()

	contact.CalcDistance(routingTable.Node.Me.Id)

	bucketIndex := routingTable.GetBucketIndex(contact.Id)

	bucket := routingTable.Buckets[bucketIndex]
	if bucket.Len() >= routingTable.Node.K {
		leastRecent := bucket.GetLeastRecentlySeenContact()
		_, err := routingTable.Node.MessageHandler.SendPingRequest(routingTable.Node.Me, &leastRecent)
		if err == nil {
			bucket.AddContact(leastRecent)
		} else {
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

	// If the target is the node itself, start from the first bucket
	if bucketIndex == -1 {
		bucketIndex = 0
	}
	bucket := routingTable.Buckets[bucketIndex]

	candidates.Append(bucket.GetContactsAndCalcDistance(target))

	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < IDLength*8) && candidates.Len() < routingTable.Node.K; i++ {
		if bucketIndex-i >= 0 {
			bucket = routingTable.Buckets[bucketIndex-i]
			candidates.Append(bucket.GetContactsAndCalcDistance(target))
		}
		if bucketIndex+i < IDLength*8 {
			bucket = routingTable.Buckets[bucketIndex+i]
			candidates.Append(bucket.GetContactsAndCalcDistance(target))
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

package kademlia_node

import (
	"sync"
)

type RoutingTable struct {
	Node    *Node
	Buckets []*bucket
	Mutex   sync.RWMutex
}

var (
	routingTableInstance *RoutingTable
	routingSingleton     sync.Once
)

// NewRoutingTable returns a new instance of a RoutingTable
func NewRoutingTable(node *Node) *RoutingTable {
	routingSingleton.Do(func() {
		routingTableInstance = &RoutingTable{
			Node: node}
		routingTableInstance.Buckets = make([]*bucket, IDLength*8)
		for i := 0; i < IDLength*8; i++ {
			routingTableInstance.Buckets[i] = NewBucket(node.K)
		}
	})
	return routingTableInstance
}

// AddContact add a new contact to the correct Bucket
func (routingTable *RoutingTable) AddContact(contact *Contact) {
	routingTable.Mutex.Lock()
	defer routingTable.Mutex.Unlock()

	// Update the distance of the contact
	contact.CalcDistance(routingTable.Node.Me.Id)

	bucketIndex := routingTable.getBucketIndex(contact.Id)
	bucket := routingTable.Buckets[bucketIndex]
	// Check if the bucket is full
	if bucket.Len() >= routingTable.Node.K {
		// Ping the least recently seen contact
		leastRecent := bucket.GetLeastRecentlySeenContact()
		_, err := routingTable.Node.MessageHandler.SendPingRequest(routingTable.Node.Me, &leastRecent)
		if err != nil {
			// If the contact is still alive, move it to the front of the bucket
			bucket.AddContact(bucket.GetLeastRecentlySeenContact())
		} else {
			// If the contact is not alive, remove it from the bucket
			bucket.RemoveContact(bucket.GetLeastRecentlySeenContact())
		}
		return
	}
	bucket.AddContact(*contact)
}

// FindClosestContacts finds the count closest Contacts to the target in the RoutingTable
func (routingTable *RoutingTable) FindClosestContacts(target *KademliaID) []*Contact {
	routingTable.Mutex.Lock()
	defer routingTable.Mutex.Unlock()

	var candidates ContactCandidates
	bucketIndex := routingTable.getBucketIndex(target)
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

// getBucketIndex get the correct Bucket index for the KademliaID
func (routingTable *RoutingTable) getBucketIndex(id *KademliaID) int {
	distance := id.CalcDistance(routingTable.Node.Me.Id)
	for i := 0; i < IDLength; i++ {
		for j := 0; j < 8; j++ {
			if (distance[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}
	return IDLength*8 - 1
}

// UpdateRoutingTable updates the RoutingTable with a list of new contacts
func (routingTable *RoutingTable) UpdateRoutingTable(contacts []*Contact) {
	for _, contact := range contacts {
		routingTable.AddContact(contact)
	}
}

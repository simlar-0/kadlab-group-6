package kademlia_node

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"sync"
)

type Node struct {
	Me             *Contact
	RoutingTable   *RoutingTable
	data           map[KademliaID][]byte
	Network        NetworkInterface
	MessageHandler MessageHandlerInterface
	K              int
	Alpha          int
}

type NodeInterface interface {
	Store(data []byte) (key *KademliaID, err error)
	LookupData(hash string) (content []byte, source *Contact, err error)
}

func (node *Node) PrintData() {
	fmt.Println(node.data)
}

// NewNode returns a new instance of a Node
func NewNode(id *KademliaID) *Node {
	k, _ := strconv.Atoi(os.Getenv("K"))
	alpha, _ := strconv.Atoi(os.Getenv("ALPHA"))
	data := make(map[KademliaID][]byte)
	ip := GetLocalIp("eth0")
	port := GetRandomPortOrDefault()
	me := NewContact(id, ip, port)

	node := &Node{
		data:  data,
		Me:    me,
		K:     k,
		Alpha: alpha,
	}

	node.RoutingTable = NewRoutingTable(node)
	node.MessageHandler = NewMessageHandler(node)
	node.Network = NewNetwork(node)
	fmt.Println("Node created with ID: ", id)
	return node
}

func (node *Node) LookupContact(target *Contact) []*Contact {
	// Uses strict parallelism to find the k closest contacts to the destination
	// i.e. Alpha concurrent FindNode requests
	shortlist := NewShortlist(target.Id, node.K)
	contacted := make(map[*KademliaID]bool)

	// Get the initial k closest contacts to the destination
	initialContacts := node.RoutingTable.FindClosestContacts(target.Id)
	for _, contact := range initialContacts {
		shortlist.AddContact(contact)
	}

	closestContact := shortlist.GetClosestContact()

	for {
		// Get the alpha closest contacts from the shortlist not contacted
		alphaClosest := shortlist.GetClosestContactsNotContacted(node.Alpha, contacted)
		responseChannel := make(chan []*Contact, len(alphaClosest))
		var wg sync.WaitGroup

		if len(alphaClosest) == 0 {
			return shortlist.GetClosestContacts(shortlist.Len())
		}

		// Send asynchronous FindNode requests to the alpha closest (not contacted) contacts in the shortlist
		for _, contact := range alphaClosest {
			contacted[contact.Id] = true

			wg.Add(1)
			go func(c *Contact) {
				defer wg.Done()
				contacts, err := node.MessageHandler.SendFindNodeRequest(node.Me, c, target.Id)
				if err != nil {
					// Dead contacts are removed from the shortlist
					shortlist.RemoveContact(c)
					return
				}
				// Add the k closest contacts from the response to the shortlist
				responseChannel <- contacts.Payload.Contacts
			}(contact)
		}
		// Wait for all goroutines to finish
		go func() {
			wg.Wait()
			close(responseChannel)
		}()

		// Process responses
		for contacts := range responseChannel {
			for _, contact := range contacts {
				// if the contact is me, skip
				if !contact.Id.Equals(node.Me.Id) {
					shortlist.AddContact(contact)
				}
			}
		}

		// Check if all the contacts in the shortlist have been contacted
		// or if the target is in the shortlist
		// or if the closest contact has not changed
		newClosestContact := shortlist.GetClosestContact()
		if newClosestContact == nil {
			return shortlist.GetClosestContacts(shortlist.Len())
		}

		if shortlist.AllContacted(contacted) || shortlist.Contains(target) || closestContact.Id.Equals(newClosestContact.Id) {
			return shortlist.GetClosestContacts(shortlist.Len())
		}
		closestContact = newClosestContact
	}
}

func (node *Node) LookupData(hash string) (content []byte, source *Contact, err error) {
	// first part similar to LookupContact
	targetID := NewKademliaID(hash)
	shortlist := NewShortlist(targetID, node.K)
	contacted := make(map[*KademliaID]bool)

	// get the initial k closest contacts to the destination
	initialContacts := node.RoutingTable.FindClosestContacts(targetID)
	for _, contact := range initialContacts {
		shortlist.AddContact(contact)
	}

	for {
		// Get the alpha closest contacts from the shortlist not contacted
		alphaClosest := shortlist.GetClosestContactsNotContacted(node.Alpha, contacted)
		responseChannel := make(chan *RPC, len(alphaClosest))
		var wg sync.WaitGroup

		if len(alphaClosest) == 0 {
			return nil, nil, fmt.Errorf("data not found")
		}

		// Send asynchronous FindNode requests to the alpha closest (not contacted) contacts in the shortlist
		for _, contact := range alphaClosest {
			contacted[contact.Id] = true
			wg.Add(1)
			go func(c *Contact) {
				defer wg.Done()
				rpc, err := node.MessageHandler.SendFindValueRequest(node.Me, c, targetID)
				if err != nil {
					// Dead contacts are removed from the shortlist
					shortlist.RemoveContact(c)
					return
				}
				// Add the k closest contacts from the response to the shortlist
				responseChannel <- rpc
			}(contact)
		}

		wg.Wait()
		close(responseChannel)

		// Process responses
		for rpc := range responseChannel {
			if rpc.Payload.Data != nil {
				return rpc.Payload.Data, rpc.Source, nil
			}
			for _, contact := range rpc.Payload.Contacts {
				if !contact.Id.Equals(node.Me.Id) {
					shortlist.AddContact(contact)
				}
			}
		}

		newClosestContacts := shortlist.GetClosestContacts(node.K)
		if len(newClosestContacts) == 0 || shortlist.AllContacted(contacted) {
			return nil, nil, fmt.Errorf("data not found")
		}
	}
}

func GenerateKey(data []byte) *KademliaID {
	hash := sha256.Sum256(data)
	return NewKademliaID(hex.EncodeToString(hash[:]))
}

func (node *Node) Store(data []byte) (key *KademliaID, err error) {

	fmt.Println("Storing the data")
	// Generate the key for the data (e.g., using a hash function)
	key = GenerateKey(data)
	target := &Contact{Id: key}
	// Use LookupContact to find the k closest nodes to the key
	closestContacts := node.LookupContact(target)
	// Send STORE_REQUEST to each of the closest contacts
	for _, contact := range closestContacts {
		r, err := node.MessageHandler.SendStoreRequest(node.Me, contact, data)
		if err != nil {
			// Handle error (e.g., log it, retry, etc.)
			fmt.Printf("Failed to store data on node %s: %v\n", contact.Id, err)
			return nil, err
		}

		if r != nil {
			return key, nil
		}
	}
	return nil, nil
}

func (node *Node) StoreData(data []byte) (err error) {
	fmt.Println("\033[0;31mData stored\033[0m")
	key := GenerateKey(data)
	node.data[*key] = data
	return nil
}

func (node *Node) GetData(key *KademliaID) (data []byte, err error) {
	fmt.Println("\033[0;31mData retrieved \033[0m")
	data = node.data[*key]
	return data, nil
}

func (node *Node) Join(contact *Contact) (err error) {
	fmt.Println("Joining the network")
	// Ping the contact to see if it is alive
	_, e := node.MessageHandler.SendPingRequest(node.Me, contact)
	if e != nil {
		return e
	}
	// Add the contact to the routing table
	node.RoutingTable.AddContact(contact)
	// Perform a lookupNode on myself
	contacts := node.LookupContact(node.Me)
	// Update the routing table with the results
	node.RoutingTable.UpdateRoutingTable(contacts)
	// Refresh all buckets further away than the closest neighbor
	node.RefreshBuckets()
	fmt.Println("Joined the network")
	return nil
}

// RefreshBuckets refreshes all buckets further away than the closest neighbor
func (node *Node) RefreshBuckets() {
	// Get the closest neighbor
	neighbor := node.RoutingTable.FindClosestContacts(node.Me.Id)[0]
	// Get the bucket index of the neighbor
	bucketIndex := node.RoutingTable.GetBucketIndex(neighbor.Id)
	// Refresh all buckets further away than the neighbor
	for i := bucketIndex + 1; i < IDLength*8; i++ {
		target := NewRandomKademliaIDInBucket(i, node.Me.Id)
		contacts := node.LookupContact(NewContact(target, "", 0))
		node.RoutingTable.UpdateRoutingTable(contacts)
	}
}

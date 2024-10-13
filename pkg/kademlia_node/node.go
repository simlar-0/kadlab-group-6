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
	data           map[KademliaID][]byte
	RoutingTable   *RoutingTable
	Network        *Network
	MessageHandler *MessageHandler
	K              int
	Alpha          int
}

type NodeCommunication interface {
	Store(data []byte) (key *KademliaID, err error)
	LookupData(hash string) (content []byte, source *Node, err error)
}

// NewNode returns a new instance of a Node
func NewNode(id *KademliaID) *Node {
	k, _ := strconv.Atoi(os.Getenv("K"))
	alpha, _ := strconv.Atoi(os.Getenv("ALPHA"))
	data := map[KademliaID][]byte{}
	ip := GetLocalIp("eth0")
	port := GetRandomPortOrDefault()
	me := NewContact(id, ip, port)
	routingTable := NewRoutingTable(me, k)
	messageHandler := NewMessageHandler(routingTable)
	network := NewNetwork(me)
	network.MessageHandler = messageHandler
	messageHandler.Network = network
	fmt.Println("Node created with ID: ", id)
	return &Node{
		Me:             me,
		data:           data,
		RoutingTable:   routingTable,
		Network:        network,
		MessageHandler: messageHandler,
		K:              k,
		Alpha:          alpha,
	}
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

		wg.Wait()
		close(responseChannel)

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
		if shortlist.AllContacted(contacted) || shortlist.Contains(target) {
			return shortlist.GetClosestContacts(node.K)
		}
	}
}
/*
func (node *Node) LookupData(hash string) (content []byte, source *Node, err error) {
	// TODO
	return nil, nil, nil
}
*/
func (node *Node) LookupData(hash string) (content []byte, source *Node, err error) {
	targetID := NewKademliaID(hash)
	shortlist := NewShortlist(targetID, node.K)
	contacted := make(map[*KademliaID]bool)
  
	initialContacts := node.RoutingTable.FindClosestContacts(targetID)
	for _, contact := range initialContacts {
	  shortlist.AddContact(contact)
	}
  
	for {
		alphaClosest := shortlist.GetClosestContactsNotContacted(node.Alpha, contacted)
		responseChannel := make(chan *RPC, len(alphaClosest))
		var wg sync.WaitGroup
	
		if len(alphaClosest) == 0 {
			return nil, nil, fmt.Errorf("Data not found")
		}
	
		for _, contact := range alphaClosest {
			contacted[contact.Id] = true
			wg.Add(1)
			go func(c *Contact) {
			defer wg.Done()
			rpc, err := node.MessageHandler.SendFindValueRequest(node.Me, c, targetID)
			if err != nil {
				shortlist.RemoveContact(c)
				return
			}
			responseChannel <- rpc
			}(contact)
		}

		wg.Wait()
		close(responseChannel)

		for rpc := range responseChannel {
			if rpc.Payload.Data != nil {
				// how to return the node?
				return rpc.Payload.Data, rpc.Node, nil // return rpc.Payload.Data, rpc.Source, nil
			}
			for _, contact := range rpc.Payload.Contacts {
				if !contact.Id.Equals(node.Me.Id) {
					shortlist.AddContact(contact)
				}
			}
		}
	
		newClosestContacts := shortlist.GetClosestContacts(node.K)
		if len(newClosestContacts) == 0 || shortlist.AllContacted(contacted) {
			return nil, nil, fmt.Errorf("Data not found")
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
		_, err := node.MessageHandler.SendStoreRequest(node.Me, contact, data)
		if err != nil {
			// Handle error (e.g., log it, retry, etc.)
			fmt.Printf("Failed to store data on node %s: %v\n", contact.Id, err)
		}
	}
	fmt.Println("Data stored")

	return key, nil
}

func (node *Node) Ping(target *Contact) (err error) {
	_, e := node.MessageHandler.SendPingRequest(node.Me, target)
	if e != nil {
		return e
	}
	node.RoutingTable.AddContact(target)
	return nil
}

func (node *Node) Join(contact *Contact) (err error) {
	fmt.Println("Joining the network")

	node.Ping(contact)
	// Perform a lookupNode on myself
	contacts := node.LookupContact(node.Me)
	fmt.Println("Lookup complete: ", contacts)
	// Update the routing table with the results
	node.RoutingTable.UpdateRoutingTable(contacts)
	// Refresh all buckets further away than the closest neighbor
	//node.routingTable.Refresh(contact.id)
	return nil
}

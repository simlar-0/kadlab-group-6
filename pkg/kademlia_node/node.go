package kademlia_node

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

type Node struct {
	Me             *Contact
	RoutingTable   *RoutingTable
	Network        *Network
	MessageHandler *MessageHandler
	K              int
	Alpha          int
	B              int
}

// NewNode returns a new instance of a Node
func NewNode(id *KademliaID) *Node {
	k, _ := strconv.Atoi(os.Getenv("K"))
	alpha, _ := strconv.Atoi(os.Getenv("ALPHA"))
	b, _ := strconv.Atoi(os.Getenv("B"))

	ip := GetLocalIp("eth0")
	port := GetRandomPortOrDefault()
	me := NewContact(id, ip, port)
	routingTable := NewRoutingTable(me, k)
	messageHandler := NewMessageHandler(routingTable)
	network := NewNetwork(me, messageHandler)
	messageHandler.Network = network

	return &Node{
		Me:             me,
		RoutingTable:   routingTable,
		Network:        network,
		MessageHandler: messageHandler,
		K:              k,
		Alpha:          alpha,
		B:              b}
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
				shortlist.AddContact(contact)
			}
		}

		// Check if all the contacts in the shortlist have been contacted
		// or if the target is in the shortlist
		if shortlist.AllContacted(contacted) || shortlist.Contains(target) {
			return shortlist.GetClosestContacts(node.K)
		}
	}
}

func (node *Node) LookupData(hash string) {
	// TODO

}

func (node *Node) Store(data []byte) {
	// TODO
}

func (node *Node) Join(contact *Contact) {
	fmt.Println("Joining the network")
	// Ping the contact to see if it is alive
	_, err := node.MessageHandler.SendPingRequest(node.Me, contact)
	if err != nil {
		// handle no response
	}
	// Add the contact to the routing table
	//node.routingTable.AddContact(contact)
	// Perform a lookupNode on myself
	//contacts := node.LookupContact(node.me)
	// Update the routing table with the results
	//node.routingTable.UpdateRoutingTable(contacts)
	// Refresh all buckets further away than the closest neighbor
	//node.routingTable.Refresh(contact.id)

}

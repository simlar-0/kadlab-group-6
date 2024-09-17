package kademlia_node

type Node struct {
	me             *Contact
	routingTable   *RoutingTable
	Network        *Network
	messageHandler *MessageHandler
}

// NewNode returns a new instance of a Node
func NewNode() *Node {
	ip := GetLocalIp("eth0")
	port := GetRandomPortOrDefault()
	me := NewContact(NewRandomKademliaID(), ip, port)
	messageHandler := NewMessageHandler()
	network := NewNetwork(me, messageHandler)
	messageHandler.Network = network
	return &Node{
		me:             me,
		routingTable:   NewRoutingTable(me),
		Network:        network,
		messageHandler: messageHandler}
}

func (node *Node) LookupContact(target *Contact) {
	// TODO
	RPC, err := node.messageHandler.SendFindNodeRequest(node.me, target, node.me.id)
	if err != nil {
		// handle no response
	}
	node.messageHandler.Network.SendRequest(RPC)
}

func (node *Node) LookupData(hash string) {
	// TODO

}

func (node *Node) Store(data []byte) {
	// TODO
}

func (node *Node) Join(contact *Contact) {
	// Add the contact to the routing table
	node.routingTable.AddContact(contact)
	// Perform a lookupNode on myself
	node.LookupContact(node.me)
	// Refresh all buckets further away than the closest neighbor
	// Update the routing table with the results
}

func (node *Node) Refresh(KademliaID *KademliaID) {
	// TODO: Refresh the bucket containing the KademliaID
	// Performs a lookupNode on the KademliaID
	// Update the routing table with the results
}

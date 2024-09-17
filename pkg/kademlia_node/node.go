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
	port := GetRandomPort()
	me := NewContact(NewRandomKademliaID(), ip, port)
	messageHandler := NewMessageHandler()
	network := NewNetwork(me, messageHandler)
	return &Node{
		me:             me,
		routingTable:   NewRoutingTable(me),
		Network:        network,
		messageHandler: messageHandler}
}

func (node *Node) LookupContact(target *Contact) {
	// TODO
}

func (node *Node) LookupData(hash string) {
	// TODO
}

func (node *Node) Store(data []byte) {
	// TODO
}


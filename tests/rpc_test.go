package tests

import (
	"fmt"
	node "kadlab-group-6/pkg/kademlia_node"
	"testing"
)

func TestNewPayload(t *testing.T) {
	key := node.NewKademliaID("1000000000000000000000000000000000000000")
	data := []byte("test data")
	contact := node.NewContact(node.NewKademliaID("2000000000000000000000000000000000000000"), "1.2.3.4", 1234)
	contacts := []*node.Contact{contact}

	payload := node.NewPayload(key, data, contacts)

	if payload.Key.String() != key.String() {
		t.Errorf("Expected Key %s, got %s", key.String(), payload.Key.String())
	}
	if string(payload.Data) != string(data) {
		t.Errorf("Expected Data %s, got %s", string(data), string(payload.Data))
	}
	if len(payload.Contacts) != len(contacts) {
		t.Errorf("Expected %d Contacts, got %d", len(contacts), len(payload.Contacts))
	}
	if payload.Contacts[0].String() != contact.String() {
		t.Errorf("Expected Contact %s, got %s", contact.String(), payload.Contacts[0].String())
	}
}

func TestNewRPC(t *testing.T) {
	id := node.NewKademliaID("0000000000000000000000000000000000000000")
	payload := node.NewPayload(id, []byte("test data"), nil)
	source := node.NewContact(node.NewKademliaID("1000000000000000000000000000000000000000"), "1.2.3.4", 1234)
	destination := node.NewContact(node.NewKademliaID("2000000000000000000000000000000000000000"), "5.6.7.8", 5678)

	rpc := node.NewRPC(node.PingRequest, false, id, payload, source, destination)

	if rpc.ID.String() != id.String() {
		t.Errorf("Expected ID %s, got %s", id.String(), rpc.ID.String())
	}
	if rpc.Type != node.PingRequest {
		t.Errorf("Expected Type %s, got %s", node.PingRequest, rpc.Type)
	}
	if rpc.IsResponse != false {
		t.Errorf("Expected IsResponse false, got %t", rpc.IsResponse)
	}
	if rpc.Source.String() != source.String() {
		t.Errorf("Expected Source %s, got %s", source.String(), rpc.Source.String())
	}
	if rpc.Destination.String() != destination.String() {
		t.Errorf("Expected Destination %s, got %s", destination.String(), rpc.Destination.String())
	}
	if rpc.Payload.String() != payload.String() {
		t.Errorf("Expected Payload %s, got %s", payload.String(), rpc.Payload.String())
	}
}

func TestValidateRPC(t *testing.T) {
	id := node.NewKademliaID("0000000000000000000000000000000000000000")
	payload := node.NewPayload(id, []byte("test data"), nil)
	source := node.NewContact(node.NewKademliaID("1000000000000000000000000000000000000000"), "1.2.3.4", 1234)
	destination := node.NewContact(node.NewKademliaID("2000000000000000000000000000000000000000"), "5.6.7.8", 5678)

	validRPC := node.NewRPC(node.PingRequest, false, id, payload, source, destination)
	invalidRPC := node.NewRPC("INVALID_TYPE", false, id, payload, source, destination)

	if !node.ValidateRPC(validRPC) {
		t.Errorf("Expected valid RPC to be valid")
	}
	if node.ValidateRPC(invalidRPC) {
		t.Errorf("Expected invalid RPC to be invalid")
	}
}

func TestRPCString(t *testing.T) {
	id := node.NewKademliaID("0000000000000000000000000000000000000000")
	payload := node.NewPayload(id, []byte("test data"), nil)
	source := node.NewContact(node.NewKademliaID("1000000000000000000000000000000000000000"), "1.2.3.4", 1234)
	destination := node.NewContact(node.NewKademliaID("2000000000000000000000000000000000000000"), "5.6.7.8", 5678)

	rpc := node.NewRPC(node.PingRequest, false, id, payload, source, destination)
	expectedString := fmt.Sprintf(`RPC(ID: "%s", Type: "%s", IsResponse: "%t", Destination: "%s", Source: "%s", Payload: "%s")`, rpc.ID, rpc.Type, rpc.IsResponse, rpc.Destination, rpc.Source, rpc.Payload)

	if rpc.String() != expectedString {
		t.Errorf("Expected %s, got %s", expectedString, rpc.String())
	}
}

func TestPayloadString(t *testing.T) {
	key := node.NewKademliaID("1000000000000000000000000000000000000000")
	data := []byte("test data")
	contact := node.NewContact(node.NewKademliaID("2000000000000000000000000000000000000000"), "1.2.3.4", 1234)
	contacts := []*node.Contact{contact}

	payload := node.NewPayload(key, data, contacts)
	expectedString := fmt.Sprintf(`Payload(Key: "%s", Data: "%s", Contacts: "%s")`, payload.Key, payload.Data, payload.Contacts)

	if payload.String() != expectedString {
		t.Errorf("Expected %s, got %s", expectedString, payload.String())
	}
}

package tests

import (
	"bytes"
	kademlia "kadlab-group-6/pkg/kademlia_node"
	mock "kadlab-group-6/pkg/mocks"
	"testing"
)

func TestInit(t *testing.T) {
	n := &mock.MockNode{}

	cli := kademlia.CliInit(n)

	if cli != nil {
		t.Errorf("cli initialisation failed")
	}
}

func TestErrorPut(t *testing.T) {
	n := &mock.MockNode{}

	cli := kademlia.CliInit(n)

	r := cli.Put(nil)

	if r != nil {
		t.Errorf("Key should be nil")
	}
}

func TestPut(t *testing.T) {
	n := &mock.MockNode{}

	cli := kademlia.CliInit(n)

	args := []string{"content"}
	r := cli.Put(args)

	if !r.Equals(kademlia.NewKademliaID("0000000000000000000000000000000000000042")) {
		t.Errorf("Key should be key")
	}
}

func TestErrorGet(t *testing.T) {
	n := &mock.MockNode{}

	cli := kademlia.CliInit(n)

	r := cli.Put(nil)

	if r != nil {
		t.Errorf("Key should be nil")
	}
}

func TestGet(t *testing.T) {
	n := &mock.MockNode{}

	cli := kademlia.CliInit(n)

	args := []string{"9f86d081884c7d659a2feaa0c55ad015a3bf4f1b"}
	r, c := cli.Get(args)

	if c == nil || r == nil {
		t.Errorf("Get error")
	}

	if !bytes.Equal(r, []byte("test")) {
		t.Errorf("Wrong data content")
	}

	if !c.Id.Equals(kademlia.NewKademliaID("0000000000000000000000000000000000001111")) {
		t.Errorf("Wrong contact id")
	}

	if c.Ip != "111.111.111.111" {
		t.Errorf("Wrong contact ip")
	}

	if c.Port != 1111 {
		t.Errorf("Wrong contact port")
	}
}

package tests

import (
	"encoding/hex"
	kademlia "kadlab-group-6/pkg/kademlia_node"
	"testing"
)

func TestNewKademliaID(t *testing.T) {
	data := "0000000000000000000000000000000000000001"
	id := kademlia.NewKademliaID(data)

	expected, _ := hex.DecodeString(data)
	for i := 0; i < kademlia.IDLength; i++ {
		if id[i] != expected[i] {
			t.Errorf("Expected byte %d to be %x, got %x", i, expected[i], id[i])
		}
	}
}

func TestNewRandomKademliaID(t *testing.T) {
	id1 := kademlia.NewRandomKademliaID()
	id2 := kademlia.NewRandomKademliaID()

	if id1.Equals(id2) {
		t.Errorf("Expected two random KademliaIDs to be different, but they are the same")
	}
}

func TestNewRandomKademliaIDInBucket(t *testing.T) {
	referenceID := kademlia.NewKademliaID("0000000000000000000000000000000000000001")
	bucketIndex := 100
	id := kademlia.NewRandomKademliaIDInBucket(bucketIndex, referenceID)

	lowerBound := kademlia.KademliaID{}
	upperBound := kademlia.KademliaID{}
	for i := 0; i < kademlia.IDLength; i++ {
		if i < bucketIndex/8 {
			lowerBound[i] = referenceID[i]
			upperBound[i] = referenceID[i]
		} else if i == bucketIndex/8 {
			lowerBound[i] = referenceID[i] & (0xFF << (8 - bucketIndex%8))
			upperBound[i] = referenceID[i] | (0xFF >> (bucketIndex % 8))
		} else {
			lowerBound[i] = 0x00
			upperBound[i] = 0xFF
		}
	}

	for i := 0; i < kademlia.IDLength; i++ {
		if id[i] < lowerBound[i] || id[i] > upperBound[i] {
			t.Errorf("Expected byte %d to be within bounds [%x, %x], got %x", i, lowerBound[i], upperBound[i], id[i])
		}
	}
}

func TestKademliaID_Less(t *testing.T) {
	id1 := kademlia.NewKademliaID("0000000000000000000000000000000000000001")
	id2 := kademlia.NewKademliaID("0000000000000000000000000000000000000002")

	if !id1.Less(id2) {
		t.Errorf("Expected %s to be less than %s", id1.String(), id2.String())
	}
	if id2.Less(id1) {
		t.Errorf("Expected %s to be not less than %s", id2.String(), id1.String())
	}
	if id1.Less(id1) {
		t.Errorf("Expected %s to be not less than %s", id1.String(), id1.String())
	}
}

func TestKademliaID_Equals(t *testing.T) {
	id1 := kademlia.NewKademliaID("0000000000000000000000000000000000000001")
	id2 := kademlia.NewKademliaID("0000000000000000000000000000000000000001")
	id3 := kademlia.NewKademliaID("0000000000000000000000000000000000000002")

	if !id1.Equals(id2) {
		t.Errorf("Expected %s to be equal to %s", id1.String(), id2.String())
	}
	if id1.Equals(id3) {
		t.Errorf("Expected %s to be not equal to %s", id1.String(), id3.String())
	}
}

func TestKademliaID_CalcDistance(t *testing.T) {
	id1 := kademlia.NewKademliaID("0000000000000000000000000000000000000001")
	id2 := kademlia.NewKademliaID("0000000000000000000000000000000000000002")
	expectedDistance := kademlia.NewKademliaID("0000000000000000000000000000000000000003")

	distance := id1.CalcDistance(id2)
	if !distance.Equals(expectedDistance) {
		t.Errorf("Expected distance to be %s, got %s", expectedDistance.String(), distance.String())
	}
}

func TestKademliaID_String(t *testing.T) {
	data := "0000000000000000000000000000000000000001"
	id := kademlia.NewKademliaID(data)

	if id.String() != data {
		t.Errorf("Expected string representation to be %s, got %s", data, id.String())
	}
}

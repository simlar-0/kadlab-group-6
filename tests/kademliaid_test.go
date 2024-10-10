package tests

import (
	"encoding/hex"
	kademlia "kadlab-group-6/pkg/kademlia_node"
	"testing"
)

func TestNewKademliaID(t *testing.T) {
	data := "0000000000000000000000000000000000000001"
	t.Logf("Creating KademliaID from data: %s", data)
	id := kademlia.NewKademliaID(data)

	expected, _ := hex.DecodeString(data)
	t.Logf("Expected KademliaID: %s", hex.EncodeToString(expected))
	for i := 0; i < kademlia.IDLength; i++ {
		if id[i] != expected[i] {
			t.Errorf("Expected byte %d to be %x, got %x", i, expected[i], id[i])
		}
	}
	t.Logf("Generated KademliaID: %s", id.String())
}

func TestNewRandomKademliaID(t *testing.T) {
	id1 := kademlia.NewRandomKademliaID()
	id2 := kademlia.NewRandomKademliaID()

	t.Logf("Generated random KademliaID1: %s", id1.String())
	t.Logf("Generated random KademliaID2: %s", id2.String())

	if id1.Equals(id2) {
		t.Errorf("Expected two random KademliaIDs to be different, but they are the same")
	}
}

func TestNewRandomKademliaIDInBucket(t *testing.T) {
	referenceID := kademlia.NewKademliaID("0000000000000000000000000000000000000001")
	bucketIndex := 100
	t.Logf("Generating random KademliaID in bucket %d with reference ID %s", bucketIndex, referenceID.String())
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

	t.Logf("Lower bound: %s", lowerBound.String())
	t.Logf("Upper bound: %s", upperBound.String())
	t.Logf("Generated ID: %s", id.String())

	for i := 0; i < kademlia.IDLength; i++ {
		if id[i] < lowerBound[i] || id[i] > upperBound[i] {
			t.Errorf("Expected byte %d to be within bounds [%x, %x], got %x", i, lowerBound[i], upperBound[i], id[i])
		}
	}
}

func TestKademliaID_Less(t *testing.T) {
	id1 := kademlia.NewKademliaID("0000000000000000000000000000000000000001")
	id2 := kademlia.NewKademliaID("0000000000000000000000000000000000000002")

	t.Logf("Comparing IDs: %s < %s", id1.String(), id2.String())
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

	t.Logf("Comparing IDs for equality: %s == %s", id1.String(), id2.String())
	if !id1.Equals(id2) {
		t.Errorf("Expected %s to be equal to %s", id1.String(), id2.String())
	}
	t.Logf("Comparing IDs for inequality: %s != %s", id1.String(), id3.String())
	if id1.Equals(id3) {
		t.Errorf("Expected %s to be not equal to %s", id1.String(), id3.String())
	}
}

func TestKademliaID_CalcDistance(t *testing.T) {
	id1 := kademlia.NewKademliaID("0000000000000000000000000000000000000001")
	id2 := kademlia.NewKademliaID("0000000000000000000000000000000000000002")
	t.Logf("ID1: %s", id1.String())
	t.Logf("ID2: %s", id2.String())

	expectedDistance := kademlia.NewKademliaID("0000000000000000000000000000000000000003")
	t.Logf("Expected distance: %s", expectedDistance.String())

	distance := id1.CalcDistance(id2)
	t.Logf("Actual distance: %s", distance.String())
	if !distance.Equals(expectedDistance) {
		t.Errorf("Expected distance to be %s, got %s", expectedDistance.String(), distance.String())
	}
}

func TestKademliaID_String(t *testing.T) {
	data := "0000000000000000000000000000000000000001"
	t.Logf("Creating KademliaID from data: %s", data)
	id := kademlia.NewKademliaID(data)

	t.Logf("Checking string representation: expected %s, got %s", data, id.String())
	if id.String() != data {
		t.Errorf("Expected string representation to be %s, got %s", data, id.String())
	}
}

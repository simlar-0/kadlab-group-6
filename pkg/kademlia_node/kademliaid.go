package kademlia_node

import (
	"encoding/hex"
	"math/rand"
	"time"
)

// the static number of bytes in a KademliaID
const IDLength = 20

// type definition of a KademliaID
type KademliaID [IDLength]byte

// NewKademliaID returns a new instance of a KademliaID based on the string input
func NewKademliaID(data string) *KademliaID {
	decoded, _ := hex.DecodeString(data)

	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = decoded[i]
	}

	return &newKademliaID
}

// NewRandomKademliaID returns a new instance of a random KademliaID,
// change this to a better version if you like
func NewRandomKademliaID() *KademliaID {
	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = uint8(rand.Intn(256))
	}
	return &newKademliaID
}

// NewRandomKademliaIDInBucket returns a new instance of a random KademliaID
// that is within the bounds of the bucket index
func NewRandomKademliaIDInBucket(bucketIndex int, referenceID *KademliaID) *KademliaID {
	lowerBound := KademliaID{}
	upperBound := KademliaID{}

	// Calculate the lower and upper bounds for the bucket
	for i := 0; i < IDLength; i++ {
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

	// Generate a random ID within the bounds
	source := rand.NewSource(time.Now().UnixNano())
	randomgen := rand.New(source)
	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		if lowerBound[i] == upperBound[i] {
			newKademliaID[i] = lowerBound[i]
		} else {
			diff := int(upperBound[i] - lowerBound[i] + 1)
			if diff > 0 {
				newKademliaID[i] = uint8(randomgen.Intn(diff)) + lowerBound[i]
			} else {
				newKademliaID[i] = lowerBound[i]
			}
		}
	}

	return &newKademliaID
}

// Less returns true if kademliaID < otherKademliaID (bitwise)
func (kademliaID KademliaID) Less(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return kademliaID[i] < otherKademliaID[i]
		}
	}
	return false
}

// Equals returns true if kademliaID == otherKademliaID (bitwise)
func (kademliaID KademliaID) Equals(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return false
		}
	}
	return true
}

// CalcDistance returns a new instance of a KademliaID that is built 
// through a bitwise XOR operation betweeen kademliaID and target
func (kademliaID KademliaID) CalcDistance(target *KademliaID) *KademliaID {
	result := KademliaID{}
	for i := 0; i < IDLength; i++ {
		result[i] = kademliaID[i] ^ target[i]
	}
	return &result
}

// String returns a simple string representation of a KademliaID
func (kademliaID *KademliaID) String() string {
	return hex.EncodeToString(kademliaID[0:IDLength])
}

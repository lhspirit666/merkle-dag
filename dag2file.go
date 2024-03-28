package merkledag

import (
	"errors"
	"strings"
)

// Hash to file
func Hash2File(store KVStore, hash []byte, path string, hp HashPool) ([]byte, error) {
	node, err := findNode(store, hash, hp)
	if err != nil {
		return nil, err
	}


	return node.LookupFile(path)
}

func findNode(store KVStore, hash []byte, hp HashPool) (*Node, error) {
	hashFunc := hp.Get()
	defer hashFunc.Reset()
	_, err := hashFunc.Write(hash)
	if err != nil {
		return nil, err
	}
	nodeID := hashFunc.Sum(nil)

	return kademliaLookup(store, nodeID)
}

func kademliaLookup(store KVStore, targetID []byte) (*Node, error) {
	var closestNode *Node
	var closestDistance int

	it := store.Iterator()
	defer it.Close()
	for it.SeekToFirst(); it.Valid(); it.Next() {
		key := it.Key()

		distance := calculateDistance(key, targetID)

		if closestNode == nil || distance < closestDistance {
			closestNode = decodeNode(it.Value())
			closestDistance = distance
		}
	}

	if closestNode == nil {
		return nil, errors.New("no closest node found")
	}

	return closestNode, nil
}

func calculateDistance(nodeID1, nodeID2 []byte) int {
	distance := 0
	for i := 0; i < len(nodeID1); i++ {
		xor := nodeID1[i] ^ nodeID2[i]
		for xor != 0 {
			distance++
			xor &= xor - 1
		}
	}
	return distance
}
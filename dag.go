package merkledag

import "hash"

func Add(store KVStore, node Node, h hash.Hash) []byte {
	// TODO 将分片写入到KVStore中，并返回Merkle Root
	switch node.Type() {
	case FILE:
		content := node.(File).Bytes()
		h.Write(content)
		return store.Put(h.Sum(nil))
	case DIR:
		it := node.(Dir).It()
		childHashes := make([][]byte, 0)
		for it.Next() {
			childNode := it.Node()
			childHash := Add(store, childNode, h)
			childHashes = append(childHashes, childHash)
		}
		// Sort childHashes for consistent order if needed
		// ...
		for _, hash := range childHashes {
			h.Write(hash)
		}
		return store.Put(h.Sum(nil))
	}
	return nil
	
}

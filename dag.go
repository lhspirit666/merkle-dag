package merkledag

import "hash"

func Add(store KVStore, node Node, h hash.Hash) []byte {
	switch n := node.(type) {
	case File:
		// 对文件节点进行哈希计算
		hashBytes := n.Bytes()
		h.Reset()
		h.Write(hashBytes)
		hashValue := h.Sum(nil)

		// 将哈希值存入 KVStore
		if err := store.Put(hashValue, hashBytes); err != nil {
			panic(err) // 错误处理可以更加健壮
		}

		return hashValue

	case Dir:
		// 对目录节点进行哈希计算
		it := n.It()
		buf := make([]byte, 0)
		for it.Next() {
			childNode := it.Node()
			childHash := Add(store, childNode, h)
			// 将子节点的哈希值写入缓冲区
			buf = append(buf, childHash...)
		}
		// 计算目录节点的哈希值
		h.Reset()
		h.Write(buf)
		hashValue := h.Sum(nil)

		// 将哈希值存入 KVStore
		if err := store.Put(hashValue, buf); err != nil {
			panic(err) // 错误处理可以更加健壮
		}

		return hashValue

	default:
		panic("Unknown node type")
	}
	return nil
}

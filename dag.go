package merkledag

import "hash"

func Add(store KVStore, node Node, h hash.Hash) []byte {
	// 递归地处理节点并计算其哈希值
	switch n := node.(type) {
	case File:
		// 对文件节点进行哈希计算
		hashBytes := n.Bytes()
		h.Reset()
		h.Write(hashBytes)
		hashValue := h.Sum(nil)

		// 将哈希值存入KVStore
		if err := store.Put(hashValue, hashBytes); err != nil {
			panic(err) // 错误处理可以更加健壮
		}

		return hashValue

	case Dir:
		// 对目录节点进行哈希计算
		it := n.It()
		buf := bytes.Buffer{}
		for it.Next() {
			childNode := it.Node()
			childHash := Add(store, childNode, h)
			// 将子节点的哈希值写入缓冲区
			buf.Write(childHash)
		}
		// 计算目录节点的哈希值
		h.Reset()
		h.Write(buf.Bytes())
		hashValue := h.Sum(nil)

		// 将哈希值存入KVStore
		if err := store.Put(hashValue, buf.Bytes()); err != nil {
			panic(err) // 错误处理可以更加健壮
		}

		return hashValue

	default:
		panic("Unknown node type")
	}

	// 处理完节点后，返回nil表示没有错误
	return nil
}

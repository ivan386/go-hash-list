package hash_list

import "hash"

type List struct {
	hasher     hash.Hash
	block_size int
	hash_list  [][]byte
	block_len  int
	size int
}

func New(hasher hash.Hash, block_size int) *List {
	list := new(List)
	list.hasher = hasher
	list.block_size = block_size
	return list
}

func (list *List) BlockSize() int {
	return list.block_size
}

func (list *List) Size() int {
	return list.size
}

func (list *List) Reset() {
	list.hasher.Reset()
	list.hash_list = nil
	list.block_len = 0
	list.size = 0
}

func (list *List) Write(data []byte) (int, error) {
	writen := len(data)

	if list.block_len > 0 {
		part := list.block_size - list.block_len
		if part <= len(data) {
			list.hasher.Write(data[:part])
			list.AppendHash(list.hasher.Sum(nil))
			list.block_len = 0
			data = data[part:]
		}
	}

	for len(data) >= list.block_size {
		list.hasher.Reset()
		list.hasher.Write(data[:list.block_size])
		list.AppendHash(list.hasher.Sum(nil))
		data = data[list.block_size:]
	}

	if len(data) > 0 {
		if list.block_len == 0 {
			list.hasher.Reset()
		}
		list.hasher.Write(data)
		list.block_len += len(data)
	}

	return writen, nil
}

func (list *List) Sum(in []byte) []byte {
	for _, hash_value := range list.GetList() {
		in = append(in, hash_value...)
	}
	return in
}

func (list *List) GetList() [][]byte {
	if list.block_len > 0 {
		list.AppendHash(list.hasher.Sum(nil))
		list.block_len = 0
	}
	return list.hash_list
}

func (list *List) HashSize() int {
	return list.hasher.Size()
}

func (list *List) AppendHash(hashes ...[]byte) {
	for _, hash_value := range hashes {
		list.hash_list = append(list.hash_list, hash_value)
		list.size += len(hash_value)
	}
}

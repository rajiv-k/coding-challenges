package main

const (
	fnvOffsetBasis uint32 = 0x811c9dc5
	fnvPrime       uint32 = 0x01000193
)

// Ref: https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function#FNV-1a_hash
func fnv32a(data []byte) uint32 {
	hash := fnvOffsetBasis
	for _, b := range data {
		hash = hash ^ uint32(b)
		hash = hash * fnvPrime
	}
	return hash
}

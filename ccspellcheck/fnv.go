package main

const (
	fnvOffsetBasis uint32 = 0x811c9dc5
	fnvPrime       uint32 = 0x01000193
)

func fnv32(data []byte) uint32 {
	hash := fnvOffsetBasis
	for _, b := range data {
		hash = hash ^ uint32(b)
		hash = hash * fnvPrime
	}
	return hash
}

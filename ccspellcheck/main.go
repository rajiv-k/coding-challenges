package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Workiva/go-datastructures/bitarray"
)

const (
	programName    = "ccspellcheck"
	programVersion = "0.0.1"
)

var (
	version bool
)

type BloomFilter struct {
	bitArray bitarray.BitArray
	numHash  uint16
	hashFunc []func([]byte) uint32
	numBits  uint64
}

func NewBloomFilter(numHash uint16, numBits uint64) *BloomFilter {

	hashFuncs := make([]func([]byte) uint32, numHash)

	for i := uint16(0); i < numHash; i++ {
		seed := uint32(i)
		hashFuncs[i] = func(data []byte) uint32 {
			h := fnv32(data)
			return h ^ seed
		}
	}

	bf := &BloomFilter{
		numHash:  numHash,
		hashFunc: hashFuncs,
		bitArray: bitarray.NewBitArray(numBits),
		numBits:  numBits,
	}
	return bf
}

func (bf *BloomFilter) buildDictionary(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("ERROR: could not open words file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	i := 0
	for scanner.Scan() {
		for _, hf := range bf.hashFunc {
			index := uint64(hf(scanner.Bytes())) % bf.bitArray.Capacity()
			bf.bitArray.SetBit(index)
		}
		i += 1
	}
	log.Printf("buildDictionary: added %v words\n", i)
}

func (bf *BloomFilter) PossiblyContains(word string) bool {
	for _, hf := range bf.hashFunc {
		index := uint64(hf([]byte(word))) % bf.bitArray.Capacity()
		isSet, err := bf.bitArray.GetBit(index)
		if err != nil {
			log.Fatalf("ERROR: error while getting bit value from bitarray: %v\n", err)
		}
		if !isSet {
			return false
		}
	}
	return true
}

func main() {
	flag.BoolVar(&version, "version", false, "print version")
	flag.Parse()

	if version {
		fmt.Printf("%v %v\n", programName, programVersion)
		os.Exit(0)
	}

	bf := NewBloomFilter(uint16(4), uint64(0x003d0900))

	// FIXME(rajiv): Reuse previously created bloom filter, if it already exists.
	bf.buildDictionary("./words.txt")

	if flag.NArg() != 1 {
		log.Fatalf("usage: %v <word>\n", programName)
	}

	word := flag.Arg(0)

	if bf.PossiblyContains(word) {
		log.Printf("'%v' found!\n", word)
	} else {
		log.Printf("'%v' is not contained in the dictionary\n", word)
	}
}

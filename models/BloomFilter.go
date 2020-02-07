package models

import (
	"encoding/binary"
	"encoding/hex"
	"transactions/utils"

	"github.com/spaolacci/murmur3"
)

const (
	//BIP37CONSTANT is the constant to generate the seed when adding items to the bloom fitler
	BIP37CONSTANT uint32 = 0xfba4c795
)

// BloomFilter datastructure
type BloomFilter struct {
	Size          uint32
	BitField      []byte
	FunctionCount uint32
	Tweak         uint32
}

// CreateBloomFilter will create a new bloom filter
func CreateBloomFilter(Size uint32, FunctionCount uint32, Tweak uint32) *BloomFilter {
	return &BloomFilter{
		Size:          Size,
		BitField:      make([]byte, Size*8),
		FunctionCount: FunctionCount,
		Tweak:         Tweak,
	}
}

// Add an item to the bloom filter
func (bloomFilter *BloomFilter) Add(item string) {
	for i := uint32(0); i < bloomFilter.FunctionCount; i++ {
		seed := i*BIP37CONSTANT + bloomFilter.Tweak
		rawBytes := []byte(item)
		h := murmur3.Sum32WithSeed(rawBytes, seed)
		bit := h % (bloomFilter.Size * 8)
		bloomFilter.BitField[bit] = 1
	}
}

// FilterBytes will return the bits to byte representation
func (bloomFilter *BloomFilter) FilterBytes() []byte {
	return *utils.BitsToBytes(bloomFilter.BitField)
}

// FilterLoad creates a filter message to send to nodes
func (bloomFilter *BloomFilter) FilterLoad(flag uint8) string {
	payload := utils.EncodeVarInt(uint64(bloomFilter.Size))
	payloadBytes, _ := hex.DecodeString(payload)
	payloadBytes = append(payloadBytes, bloomFilter.FilterBytes()...)

	// Create empty buffers
	functionCountBuf := make([]byte, 4)
	tweakBuf := make([]byte, 4)
	flagBuf := make([]byte, 1)

	binary.LittleEndian.PutUint32(functionCountBuf, bloomFilter.FunctionCount)
	binary.LittleEndian.PutUint32(tweakBuf, bloomFilter.Tweak)
	copy(flagBuf, []byte{flag})
	// Append process
	payloadBytes = append(payloadBytes, functionCountBuf...)
	payloadBytes = append(payloadBytes, tweakBuf...)
	payloadBytes = append(payloadBytes, flagBuf...)
	return hex.EncodeToString(payloadBytes)
}

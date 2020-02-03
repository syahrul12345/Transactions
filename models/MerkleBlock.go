package models

import (
	"encoding/hex"
	"fmt"
	"transactions/utils"
)

// MerkleBlock represetns a merkle block
type MerkleBlock struct {
	Version      [4]byte
	PrevBlock    [32]byte
	MerkleRoot   [32]byte
	TimeStamp    [4]byte
	Bits         [4]byte
	Nonce        [4]byte
	Total        [4]byte
	NumberHashes uint64
	Hashes       [][32]byte
	Flag         []byte
}

//ParseMerkleBlock will parse the merkleblock dump string
func ParseMerkleBlock(merkleblockdump string) *MerkleBlock {
	buf, err := hex.DecodeString(merkleblockdump)
	if err != nil {
		fmt.Println("Failed to parse decode merkleblock string into bytes")
	}
	// Empty variables to hold the values
	var Version [4]byte
	var PrevBlock [32]byte
	var MerkleRoot [32]byte
	var TimeStamp [4]byte
	var Bits [4]byte
	var Nonce [4]byte
	var Total [4]byte
	// Copy from the txDump bytearray
	verBuf := buf[0:4]
	prevBlkBuf := buf[4:36]
	merkleRootBuf := buf[36:68]
	timeStampBuf := buf[68:72]
	bitBuf := buf[72:76]
	nonceBuf := buf[76:80]
	totalBuf := buf[80:84]

	utils.Reverse(&verBuf)
	utils.Reverse(&prevBlkBuf)
	utils.Reverse(&merkleRootBuf)
	utils.Reverse(&timeStampBuf)
	utils.Reverse(&totalBuf)

	copy(Version[:], verBuf)
	copy(PrevBlock[:], prevBlkBuf)
	copy(MerkleRoot[:], merkleRootBuf)
	copy(TimeStamp[:], timeStampBuf)
	copy(Bits[:], bitBuf)
	copy(Nonce[:], nonceBuf)
	copy(Total[:], totalBuf)

	// Now handle the hashes
	numHashes, newString := utils.ReadVarInt(hex.EncodeToString(buf[84:]))
	// Empty array to hold the hashes
	hashesBuf := [][32]byte{}
	for i := uint64(0); i < numHashes; i++ {
		rawBytes, _ := hex.DecodeString(newString)
		// 32byte bytearray to hold the raw bytes
		var tempHashBuf [32]byte
		tempHash := rawBytes[:32]
		copy(tempHashBuf[:], tempHash)
		hashesBuf = append(hashesBuf, tempHashBuf)
		// need to change new string
		newString = hex.EncodeToString(rawBytes[32:])
	}
	flagLength, newString := utils.ReadVarInt(newString)
	newStringBytes, _ := hex.DecodeString(newString)
	flag := newStringBytes[:flagLength]
	return &MerkleBlock{
		Version,
		PrevBlock,
		MerkleRoot,
		TimeStamp,
		Bits,
		Nonce,
		Total,
		numHashes,
		hashesBuf,
		flag,
	}
}

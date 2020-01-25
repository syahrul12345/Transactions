package models

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Block represents a block in the blockchaibn
type Block struct {
	Version    [4]byte
	PrevBlock  [32]byte
	MerkleRoot [32]byte
	TimeStamp  [4]byte
	Bits       [4]byte
	Nonce      [4]byte
}

//ParseBlock will parse a blockHeader and return the corresponding block
func ParseBlock(txDump string) *Block {
	buf, err := hex.DecodeString(txDump)
	if err != nil {
		fmt.Println("Failed to parse decode blockheader into bytes")
	}
	// Empty variables to hold the values
	var Version [4]byte
	var PrevBlock [32]byte
	var MerkleRoot [32]byte
	var TimeStamp [4]byte
	var Bits [4]byte
	var Nonce [4]byte
	// Copy from the txDump bytearray
	verBuf := buf[0:4]
	prevBlkBuf := buf[4:36]
	merkleRootBuf := buf[36:68]
	timeStampBuf := buf[68:72]
	bitBuf := buf[72:76]
	nonceBuf := buf[76:]

	reverse(&verBuf)
	reverse(&prevBlkBuf)
	reverse(&merkleRootBuf)
	reverse(&timeStampBuf)

	copy(Version[:], verBuf)
	copy(PrevBlock[:], prevBlkBuf)
	copy(MerkleRoot[:], merkleRootBuf)
	copy(TimeStamp[:], timeStampBuf)
	copy(Bits[:], bitBuf)
	copy(Nonce[:], nonceBuf)

	return &Block{
		Version,
		PrevBlock,
		MerkleRoot,
		TimeStamp,
		Bits,
		Nonce,
	}
}

// Serialize a block
func (block *Block) Serialize() string {
	// Need to reverse alll litle endians.
	//Make sure to deep copy
	tempVer := block.Version
	tempPrevBlk := block.PrevBlock
	tempMerkleRoot := block.MerkleRoot
	tempTimeStamp := block.TimeStamp
	blockVerBuf := tempVer[:]
	prevBlkBuf := tempPrevBlk[:]
	merkleRootBuf := tempMerkleRoot[:]
	timeStampBuf := tempTimeStamp[:]
	reverse(&blockVerBuf)
	reverse(&prevBlkBuf)
	reverse(&merkleRootBuf)
	reverse(&timeStampBuf)
	return hex.EncodeToString(blockVerBuf) + hex.EncodeToString(prevBlkBuf) + hex.EncodeToString(merkleRootBuf) + hex.EncodeToString(timeStampBuf) + hex.EncodeToString(block.Bits[:]) + hex.EncodeToString(block.Nonce[:])
}

//ID willGets the ID of the current block
func (block *Block) ID() string {
	return hex.EncodeToString(block.hash())
}

func (block *Block) hash() []byte {
	byteData, _ := hex.DecodeString(block.Serialize())
	firstRound := sha256.Sum256(byteData)
	secondRound := sha256.Sum256(firstRound[:])
	for i := len(secondRound)/2 - 1; i >= 0; i-- {
		opp := len(secondRound) - 1 - i
		secondRound[i], secondRound[opp] = secondRound[opp], secondRound[i]
	}
	return secondRound[:]
}

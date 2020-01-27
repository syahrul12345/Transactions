package models

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"transactions/utils"
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

	utils.Reverse(&verBuf)
	utils.Reverse(&prevBlkBuf)
	utils.Reverse(&merkleRootBuf)
	utils.Reverse(&timeStampBuf)

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
	utils.Reverse(&blockVerBuf)
	utils.Reverse(&prevBlkBuf)
	utils.Reverse(&merkleRootBuf)
	utils.Reverse(&timeStampBuf)
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

//Bip9 Checks if the block implements Bip9
func (block *Block) Bip9() bool {
	// Version number is already in little endian
	version := binary.BigEndian.Uint32(block.Version[:])
	if version>>29 == 0x001 {
		return true
	}
	return false
}

//Bip91 Checks if the block implements Bip91
func (block *Block) Bip91() bool {
	version := binary.BigEndian.Uint32(block.Version[:])
	if version>>4&1 == 1 {
		return true
	}
	return false
}

//Bip141 checks if the block implements Bip141
func (block *Block) Bip141() bool {
	version := binary.BigEndian.Uint32(block.Version[:])
	if version>>1&1 == 1 {
		return true
	}
	return false
}

// Target will find the target for the block
func (block *Block) Target() string {
	return utils.BitsToTarget(block.Bits)
}

// Difficulty calculates the current block difficulty in base 10
func (block *Block) Difficulty() string {
	target := block.Target()
	targetBig, _ := big.NewInt(0).SetString(target, 16)
	// Do for (256**(0x1d-3))/target
	// a ** (0x1d-3) / target
	a := big.NewInt(256)
	// b = (ex1d-3)
	b := big.NewInt(0x1d - 3)
	// find c = a ** (b-3) or c = 256**(0x1d-3)
	c := big.NewInt(0).Exp(a, b, big.NewInt(0))
	// numerator = 0xffff * C
	numrator := big.NewInt(0).Mul(big.NewInt(0xffff), c)
	return big.NewInt(0).Div(numrator, targetBig).Text(10)

}

// CheckPow will check if the block header is smaller than the required target
func (block *Block) CheckPow() bool {
	// hash256 operation on the block id
	byteData, _ := hex.DecodeString(block.Serialize())
	firstRound := sha256.Sum256(byteData)
	secondRound := sha256.Sum256(firstRound[:])

	// Reverse the sha & calculate the proof
	sha := make([]byte, 32)
	copy(sha, secondRound[:])
	utils.Reverse(&sha)
	proof := big.NewInt(0).SetBytes(sha)
	// Calculate the target
	target, _ := big.NewInt(0).SetString(block.Target(), 16)
	// Compare
	return proof.Cmp(target) == -1
	// proof := utils.ToBigHex(hex.EncodeToString(sha))

}

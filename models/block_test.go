package models

import (
	"encoding/hex"
	"testing"
)

func TestBlockParse(t *testing.T) {
	blockHeaderRaw := "020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d"
	block := ParseBlock(blockHeaderRaw)
	getVer := hex.EncodeToString(block.Version[:])
	wantVer := "20000002"
	if getVer != wantVer {
		t.Errorf("Failed to parse the block version, expected %s got %s", wantVer, getVer)
	}
	getPrevBlock := hex.EncodeToString(block.PrevBlock[:])
	wantPrevBlock := "000000000000000000fd0c220a0a8c3bc5a7b487e8c8de0dfa2373b12894c38e"
	if getPrevBlock != wantPrevBlock {
		t.Errorf("Failed to parse the block prevBlockID, expected %s got %s", wantPrevBlock, getPrevBlock)
	}
	getMerkleRoot := hex.EncodeToString(block.MerkleRoot[:])
	wantMerkleRoot := "be258bfd38db61f957315c3f9e9c5e15216857398d50402d5089a8e0fc50075b"
	if getMerkleRoot != wantMerkleRoot {
		t.Errorf("Failed to parse the block merkleRoot, expected %s got %s", wantMerkleRoot, getMerkleRoot)
	}
	getTimeStamp := hex.EncodeToString(block.TimeStamp[:])
	wantTimeStamp := "59a7771e"
	if getTimeStamp != wantTimeStamp {
		t.Errorf("Failed to parse the block timestamp, expected %s got %s", wantTimeStamp, getTimeStamp)
	}
	getBits := hex.EncodeToString(block.Bits[:])
	wantBits := "e93c0118"
	if getBits != wantBits {
		t.Errorf("Failed to parse the block bits, expected %s got %s", wantBits, getBits)
	}
	getNonce := hex.EncodeToString(block.Nonce[:])
	wantNonce := "a4ffd71d"
	if getNonce != wantNonce {
		t.Errorf("Failed to parse block nocne,expected %s got %s", wantNonce, getNonce)
	}
}

func TestSerialize(t *testing.T) {
	blockHeaderRaw := "020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d"
	block := ParseBlock(blockHeaderRaw)
	get := block.Serialize()
	if get != blockHeaderRaw {
		t.Errorf("Parsed block header %s\n into block struct, but failed to serialize back into %s. Got %s instead.", blockHeaderRaw, blockHeaderRaw, get)
	}
}
func TestID(t *testing.T) {
	blockHeaderRaw := "020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d"
	block := ParseBlock(blockHeaderRaw)
	get := block.ID()
	want := "0000000000000000007e9e4c586439b0cdbe13b1370bdd9435d76a644d047523"
	if get != want {
		t.Errorf("Expected the block to have a hash of %s but got %s", want, get)
	}
}

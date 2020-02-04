package models

import (
	"encoding/binary"
	"encoding/hex"
	"testing"
	"transactions/utils"
)

func TestMerkleParse(t *testing.T) {
	MerkleBlockDump := "00000020df3b053dc46f162a9b00c7f0d5124e2676d47bbe7c5d0793a500000000000000ef445fef2ed495c275892206ca533e7411907971013ab83e3b47bd0d692d14d4dc7c835b67d8001ac157e670bf0d00000aba412a0d1480e370173072c9562becffe87aa661c1e4a6dbc305d38ec5dc088a7cf92e6458aca7b32edae818f9c2c98c37e06bf72ae0ce80649a38655ee1e27d34d9421d940b16732f24b94023e9d572a7f9ab8023434a4feb532d2adfc8c2c2158785d1bd04eb99df2e86c54bc13e139862897217400def5d72c280222c4cbaee7261831e1550dbb8fa82853e9fe506fc5fda3f7b919d8fe74b6282f92763cef8e625f977af7c8619c32a369b832bc2d051ecd9c73c51e76370ceabd4f25097c256597fa898d404ed53425de608ac6bfe426f6e2bb457f1c554866eb69dcb8d6bf6f880e9a59b3cd053e6c7060eeacaacf4dac6697dac20e4bd3f38a2ea2543d1ab7953e3430790a9f81e1c67f5b58c825acf46bd02848384eebe9af917274cdfbb1a28a5d58a23a17977def0de10d644258d9c54f886d47d293a411cb6226103b55635"
	merkleBlock := ParseMerkleBlock(MerkleBlockDump)
	versionWant := "20000000"
	versionGet := hex.EncodeToString(merkleBlock.Version[:])
	if versionGet != versionWant {
		t.Errorf("Expected the version to be %s but got %s", versionGet, versionWant)
	}
	merkleRootWant := "ef445fef2ed495c275892206ca533e7411907971013ab83e3b47bd0d692d14d4"
	merkleRoot := merkleBlock.MerkleRoot[:]
	utils.Reverse(&merkleRoot)
	merkleRootGet := hex.EncodeToString(merkleRoot)
	if merkleRootGet != merkleRootWant {
		t.Errorf("Expected the merkle root to be %s but got %s", merkleRootWant, merkleRootGet)

	}
	prevBlockWant := "df3b053dc46f162a9b00c7f0d5124e2676d47bbe7c5d0793a500000000000000"
	prevBlock := merkleBlock.PrevBlock[:]
	utils.Reverse(&prevBlock)
	prevBlockGet := hex.EncodeToString(prevBlock)
	if prevBlockGet != prevBlockWant {
		t.Errorf("Expected the previous block to be %s but got %s", prevBlockWant, prevBlockGet)
	}
	rawBytes, _ := hex.DecodeString("dc7c835b")
	utils.Reverse(&rawBytes)
	getBytes := merkleBlock.TimeStamp[:]
	getTimeStamp := hex.EncodeToString(getBytes)
	wantTimeStamp := hex.EncodeToString(rawBytes)
	if getTimeStamp != wantTimeStamp {
		t.Errorf("Expected the timestamp to be %s,but got %s", wantTimeStamp, getTimeStamp)
	}
	wantBits := hex.EncodeToString([]byte{
		0x67,
		0xd8,
		0x00,
		0x1a,
	})
	getBits := hex.EncodeToString(merkleBlock.Bits[:])
	if wantBits != getBits {
		t.Errorf("Expected the bits to be %s but got %s", wantBits, getBits)
	}
	wantNonce := hex.EncodeToString([]byte{
		0xc1,
		0x57,
		0xe6,
		0x70,
	})
	getNonce := hex.EncodeToString(merkleBlock.Nonce[:])
	if wantBits != getBits {
		t.Errorf("Expected the bits to be %s but got %s", wantNonce, getNonce)
	}
	rawBytes, _ = hex.DecodeString("bf0d0000")
	wantTotal := binary.LittleEndian.Uint32(rawBytes)
	// Merkle Block is already in little endian
	getTotal := binary.BigEndian.Uint32(merkleBlock.Total[:])
	if wantTotal != getTotal {
		t.Errorf("Expected the total number of hashes to be %d but got %d", wantTotal, getTotal)
	}
	hexHashes := []string{
		"ba412a0d1480e370173072c9562becffe87aa661c1e4a6dbc305d38ec5dc088a",
		"7cf92e6458aca7b32edae818f9c2c98c37e06bf72ae0ce80649a38655ee1e27d",
		"34d9421d940b16732f24b94023e9d572a7f9ab8023434a4feb532d2adfc8c2c2",
		"158785d1bd04eb99df2e86c54bc13e139862897217400def5d72c280222c4cba",
		"ee7261831e1550dbb8fa82853e9fe506fc5fda3f7b919d8fe74b6282f92763ce",
		"f8e625f977af7c8619c32a369b832bc2d051ecd9c73c51e76370ceabd4f25097",
		"c256597fa898d404ed53425de608ac6bfe426f6e2bb457f1c554866eb69dcb8d",
		"6bf6f880e9a59b3cd053e6c7060eeacaacf4dac6697dac20e4bd3f38a2ea2543",
		"d1ab7953e3430790a9f81e1c67f5b58c825acf46bd02848384eebe9af917274c",
		"dfbb1a28a5d58a23a17977def0de10d644258d9c54f886d47d293a411cb62261",
	}
	for i, getHashes := range merkleBlock.Hashes {
		getHash := hex.EncodeToString(getHashes[:])
		wantBytes, _ := hex.DecodeString(hexHashes[i])
		utils.Reverse(&wantBytes)
		wantHash := hex.EncodeToString(wantBytes)
		if wantHash != getHash {
			t.Errorf("Expected the transaction hash in index %d to be %s, but got %s", i, wantHash, getHash)
		}
	}
	wantFlag := "b55635"
	getFlag := hex.EncodeToString(merkleBlock.Flag)
	if getFlag != wantFlag {
		t.Errorf("Expected teh flag to be %s but got %s", wantFlag, getFlag)
	}

}
func TestIsValid(t *testing.T) {
	MerkleBlockDump := "00000020df3b053dc46f162a9b00c7f0d5124e2676d47bbe7c5d0793a500000000000000ef445fef2ed495c275892206ca533e7411907971013ab83e3b47bd0d692d14d4dc7c835b67d8001ac157e670bf0d00000aba412a0d1480e370173072c9562becffe87aa661c1e4a6dbc305d38ec5dc088a7cf92e6458aca7b32edae818f9c2c98c37e06bf72ae0ce80649a38655ee1e27d34d9421d940b16732f24b94023e9d572a7f9ab8023434a4feb532d2adfc8c2c2158785d1bd04eb99df2e86c54bc13e139862897217400def5d72c280222c4cbaee7261831e1550dbb8fa82853e9fe506fc5fda3f7b919d8fe74b6282f92763cef8e625f977af7c8619c32a369b832bc2d051ecd9c73c51e76370ceabd4f25097c256597fa898d404ed53425de608ac6bfe426f6e2bb457f1c554866eb69dcb8d6bf6f880e9a59b3cd053e6c7060eeacaacf4dac6697dac20e4bd3f38a2ea2543d1ab7953e3430790a9f81e1c67f5b58c825acf46bd02848384eebe9af917274cdfbb1a28a5d58a23a17977def0de10d644258d9c54f886d47d293a411cb6226103b55635"
	merkleBlock := ParseMerkleBlock(MerkleBlockDump)
	valid := merkleBlock.IsValid()
	if !valid {
		t.Errorf("Failed to valdiate the merkle block created from the merkle block dump")
	}

}

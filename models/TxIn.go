package models

import (
	"encoding/binary"
	"encoding/hex"
	"transactions/utils"
)

type TxIn struct {
	PrevTx    string
	PrevIndex uint32
	ScriptSig Script
	Sequence  uint32
}

//ParseTxIn Creates a TxIn object given a hash where the version and input count has been removed
func ParseTxIn(cleanedHash string) (TxIn, []byte) {
	//Get the 32 bytes which represent the prevTX hash
	prevHash := cleanedHash[0:64]
	txHash := utils.ToBigHex(prevHash)

	//Get the next 4 bytes which represent the prev index
	prevIndex := cleanedHash[64:72]
	cleanedHash = cleanedHash[72:]
	prevIndexInt := utils.FromLittleHex(prevIndex)

	//Dummy script
	// script := ParseScript(cleanedHash)
	script, byteHash := ParseScript(cleanedHash)

	//Sequence is from the next 4 bytes
	dummySequenceBytes := byteHash[0:4]
	byteHash = byteHash[4:]
	sequenceInt := utils.FromLittleHex(hex.EncodeToString(dummySequenceBytes))
	return TxIn{
		txHash,
		prevIndexInt,
		*script,
		sequenceInt,
	}, byteHash
}

//Serializes a TxIn object
func (txIn *TxIn) Serialize() string {
	// Convert all of this to the correct string in little endian
	prevTx := txIn.PrevTx
	prevIndex := txIn.PrevIndex
	scriptSig := txIn.ScriptSig
	seq := txIn.Sequence

	//toBigHex reverses a string from little to big hex, vice versa
	prevTxString := utils.ToBigHex(prevTx)
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, prevIndex)
	prevIndexString := hex.EncodeToString(buf)
	scriptSigString := scriptSig.Serialize()
	buf2 := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf2, seq)
	seqString := hex.EncodeToString(buf2)
	return prevTxString + prevIndexString + scriptSigString + seqString

}

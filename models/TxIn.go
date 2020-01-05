package models

import (
	"transactions/utils"
)

type TxIn struct {
	PrevTx    string
	PrevIndex uint32
	ScriptSig Script
	Sequence  uint32
}

//ParseTxIn Creates a TxIn object given a hash where the version and input count has been removed
func ParseTxIn(cleanedHash string) TxIn {
	//Get the 32 bytes which represent the prevTX hash
	prevHash := cleanedHash[0:64]
	txHash := utils.ToBigHex(prevHash)

	//Get the next 4 bytes which represent the prev index
	prevIndex := cleanedHash[64:72]
	cleanedHash = cleanedHash[72:]
	prevIndexInt := utils.FromLittleHex(prevIndex)

	//Dummy script
	// script := ParseScript(cleanedHash)
	script := ParseScript(cleanedHash)
	//Sequence is from the next 4 bytes
	dummySequence := "0f000000"
	seequenceInt := utils.FromLittleHex(dummySequence)

	return TxIn{
		txHash,
		prevIndexInt,
		*script,
		seequenceInt,
	}
}

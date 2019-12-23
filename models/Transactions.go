package models

import (
	"encoding/binary"
	"strconv"
	"transactions/utils"
)

//Transaction represents a go-struct of the transaction dump
type Transaction struct {
	Version    uint32
	InputCount uint64
	TxIns      []TxIn
}

//TxIn represents a Transaction Input Object
type TxIn struct {
	PrevTx    string
	PrevIndex string
	ScriptSig string
	Sequence  string
}

//Parse a transaction from the given transactionhash
func Parse(txHash string) *Transaction {
	versionHash := txHash[0:8]
	version := GetVersion(versionHash)
	//Get the hasRemoved, which is the hash without the version and input count
	inputcount, hashRemoved := utils.GetInputs(txHash)
	txInList := make([]TxIn, 0)
	for i := 1; i <= int(inputcount); i++ {
		txIn := ParseTxIn(hashRemoved)
		txInList = append(txInList, txIn)
	}

	// res := utils.EncodeToLittleEndian(inputcount)
	return &Transaction{
		version,
		inputcount,
		txInList,
	}
}

//GetVersion will ge the version of the transaction hash. It takes the version hash as input
func GetVersion(versionHash string) uint32 {
	//Parse as a big endian integer
	versionInt, _ := strconv.ParseUint(versionHash, 16, 32)
	//Conver the integer to a byte array
	intBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(intBytes, uint32(versionInt))
	data := binary.LittleEndian.Uint32(intBytes)
	return data

}

//ParseTxIn Creates a TxIn object given a hash where the version and input count has been removed
func ParseTxIn(cleanedHash string) TxIn {
	//Get the 32 bytes which represent the prevTX hash
	prevHash := cleanedHash[0:64]
	txHash := utils.ToLittleHex(prevHash)

	//Get the next 4 bytes which represent the prev index
	prevIndex := cleanedHash[64:72]
	txIndex := utils.ToLittleHex(prevIndex)

	//Dummy script
	dummyScript := "deadbeef"
	dummyHash := utils.ToLittleHex(dummyScript)
	//Sequence
	dummySequence := "beefdead"
	sequenceHash := utils.ToLittleHex(dummySequence)
	return TxIn{
		txHash,
		txIndex,
		dummyHash,
		sequenceHash,
	}
}

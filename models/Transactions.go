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

// For every transaction, we want to remove the byte array once it';s parsed

//Parse a transaction from the given transactionhash
func Parse(txHash string) *Transaction {
	versionHash := txHash[0:8]
	version := GetVersion(versionHash)
	//Get the hasRemoved, which is the hash without the version and input count
	inputcount, hashRemoved := utils.ReadVarInt(txHash[8:])
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

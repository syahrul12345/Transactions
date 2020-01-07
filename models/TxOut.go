package models

import (
	"encoding/binary"
	"encoding/hex"
)

//TxOut is an output from the transaciton dump
type TxOut struct {
	Amount       uint64
	ScriptPubKey Script
}

//ParseTxOut will parse the transaction object and input a TxOut object
func ParseTxOut(cleanedHash string) (TxOut, []byte) {
	amountHash := cleanedHash[0:16]
	cleanedHash = cleanedHash[16:]
	amountHashByteArray, _ := hex.DecodeString(amountHash)
	amount := binary.LittleEndian.Uint64(amountHashByteArray)
	script, byteHash := ParseScript(cleanedHash)
	return TxOut{
		amount,
		*script,
	}, byteHash
}

//Serialize a txOut object
func (txOut TxOut) Serialize() string {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, txOut.Amount)
	amountString := hex.EncodeToString(buf)
	scriptPubKeyString := txOut.ScriptPubKey.Serialize()
	return amountString + scriptPubKeyString
}

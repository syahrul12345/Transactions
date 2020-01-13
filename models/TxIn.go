package models

import (
	"encoding/binary"
	"encoding/hex"
	"transactions/utils"
)

//TxIn represents one incoming transaction object
type TxIn struct {
	PrevTx    string
	PrevIndex uint32
	ScriptSig Script
	Sequence  uint32
}

// @dev As seen, the TxIn object has no amount in the struct. How then can we get the amount to be spent?
// @dev It is in the Output of the PREVIOUS transaction.
// @dev Hence, a call has to be made to geth the unspent UTXO set from the previous transaction.

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

//Serialize : Serializes a TxIn object
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

//fetchTx : Fetch the previous Output for a given TxInput. Must pass a TxFetcher object and testnet flag.
func (txIn *TxIn) fetchTx(testnet bool, txFetcher *TxFetcher) []TxOut {
	// Each TxIn has a prevTx, which is a hash of a block which will return a raw transaction
	return txFetcher.Fetch(txIn.PrevTx, true, true)
}

//Value : Get the value of the corresponding UTXO to be spent for this input
func (txIn *TxIn) Value(testnet bool, txFetcher *TxFetcher) uint64 {
	txOuts := txIn.fetchTx(testnet, txFetcher)
	txIndex := txIn.PrevIndex
	return txOuts[txIndex].Amount
}

//GetScriptPubKey will get the corresponding Script object which is the PubKey to the TxIn's Script Sig
func (txIn *TxIn) GetScriptPubKey(testnet bool, txFetcher *TxFetcher) *Script {
	txOuts := txIn.fetchTx(testnet, txFetcher)
	txIndex := txIn.PrevIndex
	return &txOuts[txIndex].ScriptPubKey
}

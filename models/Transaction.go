package models

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"transactions/utils"

	"github.com/syahrul12345/secp256k1"
)

//Transaction represents a go-struct of the transaction dump
type Transaction struct {
	Testnet  bool
	Version  uint32
	TxIns    []TxIn
	TxOuts   []TxOut
	Locktime uint32
}

// For every transaction, we want to remove the byte array once it';s parsed

//Parse a transaction from the given transactionhash
func Parse(txHash string, testnet bool) *Transaction {
	versionHash := txHash[0:8]
	version := GetVersion(versionHash)
	//Get the hasRemoved, which is the hash without the version and input count
	inputcount, hashRemoved := utils.ReadVarInt(txHash[8:])
	txInList := make([]TxIn, 0)
	for i := 1; i <= int(inputcount); i++ {
		txIn, byteHash := ParseTxIn(hashRemoved)
		hashRemoved = hex.EncodeToString(byteHash)
		txInList = append(txInList, txIn)
	}
	outputcount, hashRemoved := utils.ReadVarInt(hashRemoved)
	txOutList := make([]TxOut, 0)
	for k := 1; k <= int(outputcount); k++ {
		txOut, byteHash := ParseTxOut(hashRemoved)
		hashRemoved = hex.EncodeToString(byteHash)
		txOutList = append(txOutList, txOut)
	}

	lockTime := utils.FromLittleHex(hashRemoved)
	return &Transaction{
		testnet,
		version,
		txInList,
		txOutList,
		lockTime,
	}
}

//Serialize the TX
func (tx *Transaction) Serialize() string {
	// Lets first parse the version
	var res string
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, tx.Version)
	//Encode the version
	res = hex.EncodeToString(buf)
	//Encode the amount of txin
	res = res + utils.EncodeToLittleEndian(uint64(len(tx.TxIns)))
	//Encode the txIns
	for _, txIn := range tx.TxIns {
		res = res + txIn.Serialize()
	}
	//Encode the amount of txout
	res = res + utils.EncodeToLittleEndian(uint64(len(tx.TxOuts)))
	for _, txOut := range tx.TxOuts {
		res = res + txOut.Serialize()
	}
	//Encode the loctime
	buf2 := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf2, tx.Locktime)
	res = res + hex.EncodeToString(buf2)
	return res
}

//ID returns the transaction ID of the transaciton data
func (tx *Transaction) ID() string {
	return hex.EncodeToString(tx.hash())
}

func (tx *Transaction) hash() []byte {
	// Executes two rounds of sha256 which is hash 256
	byteData, _ := hex.DecodeString(tx.Serialize())
	firstRound := sha256.Sum256(byteData)
	secondRound := sha256.Sum256(firstRound[:])
	// reverse it
	for i := len(secondRound)/2 - 1; i >= 0; i-- {
		opp := len(secondRound) - 1 - i
		secondRound[i], secondRound[opp] = secondRound[opp], secondRound[i]
	}
	return secondRound[:]
}

//SigHash will calculate the signature hash for a particular TxIn, given the TxIn's index in the TxIns list.
func (tx *Transaction) SigHash(index uint64, redeemScript *Script) string {
	// Let's rebuild a custom serialization procedure.
	buf := make([]byte, 4)
	//Serialzie the byte
	binary.LittleEndian.PutUint32(buf, tx.Version)

	//Append var int for txin
	varintString := utils.EncodeToLittleEndian(uint64(len(tx.TxIns)))
	varIntByteArray, _ := hex.DecodeString(varintString)
	buf = append(buf, varIntByteArray...)
	txFetcher := CreateTxFetcher("https://blockchain.info/rawtx/", tx.Testnet)
	for i, txIn := range tx.TxIns {
		if uint64(i) == index {
			if redeemScript != nil {
				// We replace the scriptsig with the script pub key
				temp := &TxIn{
					txIn.PrevTx,
					txIn.PrevIndex,
					*redeemScript,
					txIn.Sequence,
				}
				tempSerialize, _ := hex.DecodeString(temp.Serialize())
				// Append the new script object
				buf = append(buf, tempSerialize...)
			} else {
				// We replace the scriptsig with the script pub key
				temp := &TxIn{
					txIn.PrevTx,
					txIn.PrevIndex,
					*txIn.GetScriptPubKey(tx.Testnet, txFetcher),
					txIn.Sequence,
				}
				tempSerialize, _ := hex.DecodeString(temp.Serialize())
				// Append the new script object
				buf = append(buf, tempSerialize...)
			}
		} else {
			temp := &TxIn{
				txIn.PrevTx,
				txIn.PrevIndex,
				Script{
					[][]byte{
						[]byte{0},
					},
				},
				txIn.Sequence,
			}
			tempSerialize, _ := hex.DecodeString(temp.Serialize())
			// Append the new script object
			buf = append(buf, tempSerialize...)
		}

	}
	//Append var int for txout
	varintString = utils.EncodeToLittleEndian(uint64(len(tx.TxOuts)))
	varIntByteArray, _ = hex.DecodeString(varintString)
	buf = append(buf, varIntByteArray...)
	for _, txOut := range tx.TxOuts {
		temp := &TxOut{
			txOut.Amount,
			txOut.ScriptPubKey,
		}
		tempSerialize, _ := hex.DecodeString(temp.Serialize())
		// Append the new script object
		buf = append(buf, tempSerialize...)
	}
	//Decode the locktime into a buffer and add it to the main buffer
	locktimeBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(locktimeBuf, tx.Locktime)
	buf = append(buf, locktimeBuf...)
	//Do the same thing for the SIG_ALL_HASH
	opcodeBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(opcodeBuf, 1)
	buf = append(buf, opcodeBuf...)
	//Do hash256
	firstRound := sha256.Sum256(buf)
	secondRound := sha256.Sum256(firstRound[:])

	return hex.EncodeToString(secondRound[:])
}

//VerifyInput : Verifies if an input, contained within the Tx object can be unlocked. Provide the index of the tx object
func (tx *Transaction) VerifyInput(index uint64) bool {
	txIn := tx.TxIns[index]
	txFetcher := CreateTxFetcher("https://blockchain.info/rawtx/", tx.Testnet)
	scriptPubKey := txIn.GetScriptPubKey(tx.Testnet, txFetcher)
	// Get the signature z of the scriptSig for that input
	var redeemScript *Script
	if scriptPubKey.ISP2SH() {
		fmt.Println("is p2sh")
		redeemScriptCommand := txIn.ScriptSig.Commands[len(txIn.ScriptSig.Commands)-1]
		redeemScriptString := utils.EncodeToLittleEndian(uint64(len(redeemScriptCommand))) + hex.EncodeToString(redeemScriptCommand)
		redeemScript, _ = ParseScript(redeemScriptString)
	} else {
		redeemScript = nil
	}
	z := tx.SigHash(index, redeemScript)
	combinedScript := txIn.ScriptSig.Add(scriptPubKey)
	return combinedScript.Evaluate("0x" + z)
}

// SignInput : Sign a input at inputIndex, with a provided privateKey.
func (tx *Transaction) SignInput(inputIndex uint64, privateKey string) bool {
	z := tx.SigHash(inputIndex, nil)
	signature, _ := secp256k1.Sign(privateKey, z)
	der := signature.DER()
	sigBytes, _ := hex.DecodeString(der)
	sigBytes = append(sigBytes, byte(1))
	// Set default compression as true
	sec := secp256k1.GetSec(privateKey, true)
	secBytes, _ := hex.DecodeString(sec)
	//
	tx.TxIns[inputIndex].ScriptSig = Script{
		[][]byte{
			sigBytes,
			secBytes,
		},
	}
	return tx.VerifyInput(inputIndex)

}

// Verify : Verifies if the transaction is correct
func (tx *Transaction) Verify() bool {
	if tx.Fee() < 0 {
		return false
	}
	for i := range tx.TxIns {
		fmt.Printf("Verifying for input %d\n", i)
		if !tx.VerifyInput(uint64(i)) {
			return false
		}
	}
	return true
}

//Fee : Get the Fee to be earned by the miner
func (tx *Transaction) Fee() uint64 {
	txFetcher := CreateTxFetcher("https://blockchain.info/rawtx/", tx.Testnet)
	// Fee is simply the sum of input - outpun
	var inputSum uint64
	var outputSum uint64
	//Asynchronous code part
	inputAmounts := make(chan uint64, len(tx.TxIns))
	for i, txIn := range tx.TxIns {
		//Do the call asynchronously
		go func(i int, txIn TxIn) {
			inputAmounts <- txIn.Value(tx.Testnet, txFetcher)
		}(i, txIn)
	}
	//Count will check when to close the channel
	var count int
	for val := range inputAmounts {
		count = count + 1
		if count == len(tx.TxIns) {
			close(inputAmounts)
		}
		inputSum = inputSum + val
	}
	for _, txOut := range tx.TxOuts {
		outputSum = outputSum + txOut.Amount
	}
	return inputSum - outputSum
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

// CreateTxIn : Create an TxIn object with the provided prevtx,previd and empty scriptsig
func CreateTxIn(prevTxID string, prevTxIndex uint32) *TxIn {
	return &TxIn{
		PrevTx:    prevTxID,
		PrevIndex: prevTxIndex,
		ScriptSig: Script{},
		Sequence:  0xffffffff,
	}
}

//CreateTxOut : Create a TxOut object with the provided Amount and script Object
func CreateTxOut(amount uint64, scriptpubkey *Script) TxOut {
	return TxOut{
		Amount:       amount,
		ScriptPubKey: scriptpubkey,
	}
}

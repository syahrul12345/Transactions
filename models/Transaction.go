package models

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"strconv"
	"transactions/utils"
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
func (tx *Transaction) SigHash(index uint64, testnet bool) string {
	// Let's rebuild a custom serialization procedure.
	buf := make([]byte, 4)
	//Serialzie the byte
	binary.LittleEndian.PutUint32(buf, tx.Version)

	//Append var int for txin
	varintString := utils.EncodeToLittleEndian(uint64(len(tx.TxIns)))
	varIntByteArray, _ := hex.DecodeString(varintString)
	buf = append(buf, varIntByteArray...)

	txFetcher := CreateTxFetcher("https://blockchain.info/rawtx/", true)
	for i, txIn := range tx.TxIns {
		if uint64(i) == index {
			// We replace the scriptsig with the script pub key
			temp := &TxIn{
				txIn.PrevTx,
				txIn.PrevIndex,
				*txIn.GetScriptPubKey(testnet, txFetcher),
				txIn.Sequence,
			}
			tempSerialize, _ := hex.DecodeString(temp.Serialize())
			// Append the new script object
			buf = append(buf, tempSerialize...)
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
	txFetcher := CreateTxFetcher("https://blockchain.info/rawtx/", true)
	scriptPubKey := txIn.GetScriptPubKey(tx.Testnet, txFetcher)
	// Get the signature z of the scriptSig for that input
	z := tx.SigHash(index, tx.Testnet)
	combinedScript := txIn.ScriptSig.Add(scriptPubKey)
	return combinedScript.Evaluate(z)
}

//Fee : Get the Fee to be earned by the miner
func (tx *Transaction) Fee() uint64 {
	txFetcher := CreateTxFetcher("https://blockchain.info/rawtx/", true)
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

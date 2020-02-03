package main

import (
	"transactions/models"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	// Broadcast a transaction
	// secret := "syahrulnizam"
	// secretNum := binary.BigEndian.Uint64([]byte(secret))
	// secretNumString := strconv.FormatUint(secretNum, 10)
	// fromAddress := secp256k1.GetTestnetAddressFromPrivKey(secretNumString)
	// targetAmount := 0.01
	// changeAmount := 0.09
	// toAddress := "mwJn1YPMq7y5F8J3LkC5Hxg9PHyZ5K4cFv"
	// // Get the output that we want to spend :
	// prevTxID := "8bdd30e76dd8ad6bda64a5bf2e234f4a1e552e74af22de9a8ba45ab7ba43e5bb"
	// prevTxIndex := uint32(1)
	// // Create the txin object. This tx in object will want to spend the output above
	// TxIns := []models.TxIn{}
	// // Create One txIn Object which is empty

	// // We'll put the correct sequence later
	// TxIns = append(TxIns, *models.CreateTxIn(prevTxID, prevTxIndex))
	// // Lets build the outputs
	// TxOuts := []models.TxOut{}

	// // First output
	// h160 := utils.Decode58(toAddress)
	// scriptPubKey := models.CreateScriptPubKeyForP2PKH(h160)
	// targetSatoshi := uint64(targetAmount * 100000000)
	// TxOuts = append(TxOuts, models.CreateTxOut(targetSatoshi, scriptPubKey))
	// // Lets send back the change to ourselves
	// h160 = utils.Decode58(fromAddress)
	// scriptPubKey = models.CreateScriptPubKeyForP2PKH(h160)
	// changeSatoshi := uint64(changeAmount * 100000000)
	// TxOuts = append(TxOuts, models.CreateTxOut(changeSatoshi, scriptPubKey))

	// // Let's create the transaction object
	// Tx := models.Transaction{
	// 	Testnet:  true,
	// 	Version:  1,
	// 	TxIns:    TxIns,
	// 	TxOuts:   TxOuts,
	// 	Locktime: 0,
	// }
	// res := Tx.IsCoinbase()
	// We will now need to sign the TxObject. The Tx we want to spend now has an empty scriptsig.
	// res := Tx.SignInput(0, secretNumString)
	// fmt.Println(res)
	// blockHeader := "020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d"
	// block := models.ParseBlock(blockHeader)
	// fmt.Println(block)
	tree := models.CreateMerkleTree(16)
	spew.Dump(tree)
}

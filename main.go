package main

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
	// // We will now need to sign the TxObject. The Tx we want to spend now has an empty scriptsig.
	// res := Tx.SignInput(0, secretNumString)
	// fmt.Println(res)
}

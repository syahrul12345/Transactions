package models

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"
)

func TestParse(t *testing.T) {
	requirements := []map[string]string{
		{
			"input": "0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600",
			"want":  "1",
		},
	}
	for _, requirement := range requirements {
		tx := Parse(requirement["input"], true)
		version := tx.Version
		want, _ := strconv.ParseInt(requirement["want"], 10, 32)
		if version != uint32(want) {
			t.Errorf("Version for the provided blockchash is not correct. Got %d got %d", version, want)
		}
	}
}
func TestParseInputs(t *testing.T) {
	requirements := []map[string]string{
		{
			"input":                   "0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600",
			"txInCount":               "1",
			"txInZeroPrevTx":          "d1c789a9c60383bf715f3f6ad9d14b91fe55f3deb369fe5d9280cb1a01793f81",
			"txInZeroPrevIndex":       "0",
			"txInZeroScriptSerialize": "6b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278a",
			"txInZeroSequence":        "0xfffffffe",
		},
	}
	for _, requirement := range requirements {
		tx := Parse(requirement["input"], true)

		txInCount := len(tx.TxIns)
		txInCountWant, _ := strconv.ParseInt(requirement["txInCount"], 10, 64)
		if txInCount != int(txInCountWant) {
			t.Errorf("Blockchain hash provided should have %d inputs, but received %d inputs", int(txInCountWant), txInCount)
		}

		txInZeroPrevTx := tx.TxIns[0].PrevTx
		if txInZeroPrevTx != requirement["txInZeroPrevTx"] {
			t.Errorf("Expected a previous Tx hash of %s but got %s for TxIn[%d]", requirement["txInZeroPrevTx"], txInZeroPrevTx, 0)
		}

		txInZeroPrevIndex := tx.TxIns[0].PrevIndex
		txInZeroPrevIndexWant, _ := strconv.ParseInt(requirement["txInZeroPrevIndex"], 10, 64)
		if txInZeroPrevIndex != uint32(txInZeroPrevIndexWant) {
			t.Errorf("Expected the selected TxIn[%d] to have a previous index of %d, but got %d", 0, txInZeroPrevIndexWant, txInZeroPrevIndex)
		}

		txInZeroScriptSerialize := tx.TxIns[0].ScriptSig.Serialize()
		txInZeroScriptSerializeWant := requirement["txInZeroScriptSerialize"]
		if txInZeroScriptSerialize != txInZeroScriptSerializeWant {
			t.Errorf("Expected the selected TxIn[%d] to have a script object of %s but got %s", 0, txInZeroScriptSerializeWant, txInZeroScriptSerialize)
		}

		txInZeroSequence := tx.TxIns[0].Sequence
		// Conver to bigendian endian. We know its 4 bytes already in bigendian
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, txInZeroSequence)
		sequenceGet := "0x" + hex.EncodeToString(buf)
		sequenceWant := requirement["txInZeroSequence"]
		if sequenceGet != sequenceWant {
			t.Errorf("Expected the selectedf TxIn[%d] to have a locktime of %s but got %s", 0, sequenceWant, sequenceGet)
		}
	}
}
func TestParseOutput(t *testing.T) {
	requirements := []map[string]interface{}{
		{
			"input":      "0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600",
			"TxOutCount": "2",
			"TxVars": []map[string]string{
				{
					"amount":       "32454049",
					"scriptpubkey": "1976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac",
				},
				{
					"amount":       "10011545",
					"scriptpubkey": "1976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac",
				},
			},
		},
	}
	for _, requirement := range requirements {
		txInputString := fmt.Sprintf("%v", requirement["input"])
		tx := Parse(txInputString, true)
		txOutCountWant, _ := strconv.ParseInt(fmt.Sprintf("%v", requirement["TxOutCount"]), 10, 64)
		if int64(len(tx.TxOuts)) != txOutCountWant {
			t.Errorf("Expected there to be %d outputs, but instead there were %d outputs", len(tx.TxOuts), txOutCountWant)
		}
		TxVars := requirement["TxVars"].([]map[string]string)
		for i, TxVar := range TxVars {
			txOut := tx.TxOuts[i]
			amount := txOut.Amount
			pubkey := txOut.ScriptPubKey.Serialize()
			amountWant, _ := strconv.ParseUint(TxVar["amount"], 10, 64)
			if amount != amountWant {
				t.Errorf("For transaction output at index %d,expected to have amount of %d but got %d", i, amountWant, amount)
			}
			if pubkey != TxVar["scriptpubkey"] {
				t.Errorf("For transaction output at index %d,expected to have script pub key %s,but got %s", i, TxVar["scriptpubkey"], pubkey)
			}
		}

	}
}

func TestLocktime(t *testing.T) {
	inputString := "0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600"
	tx := Parse(inputString, true)
	if tx.Locktime != 410393 {
		t.Errorf("Expected to get a locktime of %d but got %d", 410393, tx.Locktime)
	}
}

func TestFee(t *testing.T) {
	requirements := []map[string]string{
		{
			"input": "0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600",
			"fee":   "40000",
		},
		{
			"input": "010000000456919960ac691763688d3d3bcea9ad6ecaf875df5339e148a1fc61c6ed7a069e010000006a47304402204585bcdef85e6b1c6af5c2669d4830ff86e42dd205c0e089bc2a821657e951c002201024a10366077f87d6bce1f7100ad8cfa8a064b39d4e8fe4ea13a7b71aa8180f012102f0da57e85eec2934a82a585ea337ce2f4998b50ae699dd79f5880e253dafafb7feffffffeb8f51f4038dc17e6313cf831d4f02281c2a468bde0fafd37f1bf882729e7fd3000000006a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937feffffff567bf40595119d1bb8a3037c356efd56170b64cbcc160fb028fa10704b45d775000000006a47304402204c7c7818424c7f7911da6cddc59655a70af1cb5eaf17c69dadbfc74ffa0b662f02207599e08bc8023693ad4e9527dc42c34210f7a7d1d1ddfc8492b654a11e7620a0012102158b46fbdff65d0172b7989aec8850aa0dae49abfb84c81ae6e5b251a58ace5cfeffffffd63a5e6c16e620f86f375925b21cabaf736c779f88fd04dcad51d26690f7f345010000006a47304402200633ea0d3314bea0d95b3cd8dadb2ef79ea8331ffe1e61f762c0f6daea0fabde022029f23b3e9c30f080446150b23852028751635dcee2be669c2a1686a4b5edf304012103ffd6f4a67e94aba353a00882e563ff2722eb4cff0ad6006e86ee20dfe7520d55feffffff0251430f00000000001976a914ab0c0b2e98b1ab6dbf67d4750b0a56244948a87988ac005a6202000000001976a9143c82d7df364eb6c75be8c80df2b3eda8db57397088ac46430600",
			"fee":   "140500",
		},
	}
	for _, requirement := range requirements {
		tx := Parse(requirement["input"], false)
		fee := tx.Fee()
		feewant, _ := strconv.ParseUint(requirement["fee"], 10, 64)
		if fee != feewant {
			t.Errorf("Expected to get a fee of %d but got %d", feewant, fee)
		}
	}
}

func TestSigHash(t *testing.T) {
	txFetcher := CreateTxFetcher("https://blockchain.info/rawtx/", false)
	tx := txFetcher.FetchTx("452c629d67e41baec3ac6f04fe744b4b9617f8f859c63b3002f8684e7a4fee03")
	get := tx.SigHash(0, nil)
	want := "27e0c5994dec7824e56dec6b2fcb342eb7cdb0d0957c2fce9882f715e85d81a6"
	if get != want {
		t.Errorf("Expected the signature of the transaction to be %s, but got %s", want, get)
	}
}

func TestVerifyP2Pkh(t *testing.T) {
	txFetcher := CreateTxFetcher("https://blockchain.info/rawtx/", false)
	tx := txFetcher.FetchTx("452c629d67e41baec3ac6f04fe744b4b9617f8f859c63b3002f8684e7a4fee03")

	verified := tx.Verify()
	if !verified {
		t.Errorf("Unable to verify with tx hash %s", "452c629d67e41baec3ac6f04fe744b4b9617f8f859c63b3002f8684e7a4fee03")
	}

	txFetcher = CreateTxFetcher("https://blockchain.info/rawtx/", true)
	tx = txFetcher.FetchTx("5418099cc755cb9dd3ebc6cf1a7888ad53a1a3beb5a025bce89eb1bf7f1650a2")
	verified = tx.Verify()
	if !verified {
		t.Errorf("Unable to verify with tx hash %s", "5418099cc755cb9dd3ebc6cf1a7888ad53a1a3beb5a025bce89eb1bf7f1650a2")
	}
}
func TestSignInput(t *testing.T) {
	// Doesn't work looks like its an invalid transaction.
	txDump := "010000000199a24308080ab26e6fb65c4eccfadf76749bb5bfa8cb08f291320b3c21e56f0d0d00000000ffffffff02408af701000000001976a914d52ad7ca9b3d096a38e752c2018e6fbc40cdf26f88ac80969800000000001976a914507b27411ccf7f16f10297de6cef3f291623eddf88ac00000000"
	tx := Parse(txDump, true)
	res := tx.SignInput(0, "8675309")
	if !res {
		t.Error("Failed to sign the input succesfully")
	}
}

func TestVerifyP2SH(t *testing.T) {
	txFetcher := CreateTxFetcher("https://blockchain.info/rawtx/", false)
	tx := txFetcher.FetchTx("46df1a9484d0a81d03ce0ee543ab6e1a23ed06175c104a178268fad381216c2b")
	verified := tx.Verify()
	if !verified {
		t.Errorf("Unable to verify with tx hash %s", "46df1a9484d0a81d03ce0ee543ab6e1a23ed06175c104a178268fad381216c2b")
	}
}

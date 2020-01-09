package models

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//TxFetcher will get the details of the latest Tx From the Blockchain. By default, it should ping itself, but we can get it from a third party too.
//Crate only one txfetcher instance.
type TxFetcher struct {
	url     string
	testnet bool
	TxIDs   []string
}

type txOutTemp struct {
	Value  uint64 `json:"value"`
	Script string `json:"script"`
}

//CreateTxFetcher will crate a new TxFetcher Object
func CreateTxFetcher(url string, testnet bool) *TxFetcher {
	return &TxFetcher{
		url,
		testnet,
		[]string{},
	}
}

//GetURL will return the url of the TxFetcher object
func (txFetcher *TxFetcher) GetURL() string {
	return txFetcher.url
}

//Fetch will fetch a transaction dump if a txid is given. Have to provide testnet and fresh flags. TxID has to be in hexadecimal
func (txFetcher *TxFetcher) Fetch(txID string, testnet bool, fresh bool) []TxOut {
	url := fmt.Sprintf("%s%s?format=hex", txFetcher.GetURL(), txID)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	//Close the body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	txInfo := string(body)
	tx := Parse(txInfo, testnet)
	//Hold a a cache of 20 TxIds
	if len(txFetcher.TxIDs) < 20 {
		txFetcher.TxIDs = append(txFetcher.TxIDs, tx.ID())
	} else {
		txFetcher.TxIDs = txFetcher.TxIDs[1:19]
	}

	// Create an array of standard TxOuts
	return tx.TxOuts
}

//Write to the cache

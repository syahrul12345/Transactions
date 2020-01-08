package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//TxFetcher will get the details of the latest Tx From the Blockchain. By default, it should ping itself, but we can get it from a third party too
type TxFetcher struct {
	url     string
	testnet bool
}

//BitcoinResponse is the block information
type BitcoinResponse struct {
	Ver  uint64
	Time uint64
	Bits uint64
	NTX  uint64 `json:"n_tx"`
}

//CreateTxFetcher will crate a new TxFetcher Object
func CreateTxFetcher(url string, testnet bool) *TxFetcher {
	return &TxFetcher{
		url,
		testnet,
	}
}

//GetUrl will return the url of the TxFetcher object
func (txFetcher *TxFetcher) GetUrl() string {
	return txFetcher.url
}

//Fetch will fetch a transaction dump if a txid is given. Have to provide testnet and fresh flags. TxID has to be in hexadecimal
func (txFetcher *TxFetcher) Fetch(txID string, testnet bool, fresh bool) string {
	url := fmt.Sprintf("%s%s", txFetcher.GetUrl(), txID)
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
	var bitcoinResponse BitcoinResponse
	json.Unmarshal(body, &bitcoinResponse)
	fmt.Println(&bitcoinResponse)
	return string(body)
}

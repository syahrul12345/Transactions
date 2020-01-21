module transactions

go 1.13

require (
	github.com/btcsuite/btcutil v1.0.1 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/ethereum/go-ethereum v1.9.9 // indirect
	github.com/kisielk/godepgraph v0.0.0-20190626013829-57a7e4a651a9 // indirect
	github.com/syahrul12345/secp256k1 v0.0.0-20200121092528-2176b9d91c67
	golang.org/x/crypto v0.0.0-20200117160349-530e935923ad
)

// Only when want to use local packages
// replace github.com/syahrul12345/secp256k1 => ./secp256k1

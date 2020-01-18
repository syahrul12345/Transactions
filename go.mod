module transactions

go 1.13

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/ethereum/go-ethereum v1.9.9 // indirect
	github.com/syahrul12345/secp256k1 v0.0.0-20200118083920-72131c241529
	golang.org/x/crypto v0.0.0-20200115085410-6d4e4cb37c7d
)

// Only when want to use local packages
// replace github.com/syahrul12345/secp256k1 => ./secp256k1

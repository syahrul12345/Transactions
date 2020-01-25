package models

import (
	"encoding/hex"
	"testing"
)

func TestBlockParse(t *testing.T) {
	blockHeaderRaw := "020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d"
	block := ParseBlock(blockHeaderRaw)
	getVer := hex.EncodeToString(block.Version[:])
	wantVer := "20000002"
	if getVer != wantVer {
		t.Errorf("Failed to parse the block version, expected %s got %s", wantVer, getVer)
	}
	getPrevBlock := hex.EncodeToString(block.PrevBlock[:])
	wantPrevBlock := "000000000000000000fd0c220a0a8c3bc5a7b487e8c8de0dfa2373b12894c38e"
	if getPrevBlock != wantPrevBlock {
		t.Errorf("Failed to parse the block prevBlockID, expected %s got %s", wantPrevBlock, getPrevBlock)
	}
	getMerkleRoot := hex.EncodeToString(block.MerkleRoot[:])
	wantMerkleRoot := "be258bfd38db61f957315c3f9e9c5e15216857398d50402d5089a8e0fc50075b"
	if getMerkleRoot != wantMerkleRoot {
		t.Errorf("Failed to parse the block merkleRoot, expected %s got %s", wantMerkleRoot, getMerkleRoot)
	}
	getTimeStamp := hex.EncodeToString(block.TimeStamp[:])
	wantTimeStamp := "59a7771e"
	if getTimeStamp != wantTimeStamp {
		t.Errorf("Failed to parse the block timestamp, expected %s got %s", wantTimeStamp, getTimeStamp)
	}
	getBits := hex.EncodeToString(block.Bits[:])
	wantBits := "e93c0118"
	if getBits != wantBits {
		t.Errorf("Failed to parse the block bits, expected %s got %s", wantBits, getBits)
	}
	getNonce := hex.EncodeToString(block.Nonce[:])
	wantNonce := "a4ffd71d"
	if getNonce != wantNonce {
		t.Errorf("Failed to parse block nocne,expected %s got %s", wantNonce, getNonce)
	}
}

func TestSerialize(t *testing.T) {
	blockHeaderRaw := "020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d"
	block := ParseBlock(blockHeaderRaw)
	get := block.Serialize()
	if get != blockHeaderRaw {
		t.Errorf("Parsed block header %s\n into block struct, but failed to serialize back into %s. Got %s instead.", blockHeaderRaw, blockHeaderRaw, get)
	}
}
func TestID(t *testing.T) {
	blockHeaderRaw := "020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d"
	block := ParseBlock(blockHeaderRaw)
	get := block.ID()
	want := "0000000000000000007e9e4c586439b0cdbe13b1370bdd9435d76a644d047523"
	if get != want {
		t.Errorf("Expected the block to have a hash of %s but got %s", want, get)
	}
}

func TestBip9(t *testing.T) {
	blockHeaderRaw := "020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d"
	block := ParseBlock(blockHeaderRaw)
	if !block.Bip9() {
		t.Errorf("Expected the block to implement bip9 but it didnt..")
	}
	blockHeaderRaw = "0400000039fa821848781f027a2e6dfabbf6bda920d9ae61b63400030000000000000000ecae536a304042e3154be0e3e9a8220e5568c3433a9ab49ac4cbb74f8df8e8b0cc2acf569fb9061806652c27"
	block = ParseBlock(blockHeaderRaw)
	if block.Bip9() {
		t.Errorf("Expected teh block to NOT implement bip9 ")
	}
}

func TestBip91(t *testing.T) {
	blockHeaderRaw := "1200002028856ec5bca29cf76980d368b0a163a0bb81fc192951270100000000000000003288f32a2831833c31a25401c52093eb545d28157e200a64b21b3ae8f21c507401877b5935470118144dbfd1"
	block := ParseBlock(blockHeaderRaw)
	if !block.Bip91() {
		t.Errorf("Expected the block to implement bip91 but it didnt")
	}

	blockHeaderRaw = "020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d"
	block = ParseBlock(blockHeaderRaw)
	if block.Bip91() {
		t.Errorf("Expected teh block to NOT implement bip91 ")
	}
}
func TestBip141(t *testing.T) {
	blockHeaderRaw := "020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d"
	block := ParseBlock(blockHeaderRaw)
	if !block.Bip141() {
		t.Errorf("Expected the block to implement bip141 but it didnt")
	}

	blockHeaderRaw = "0000002066f09203c1cf5ef1531f24ed21b1915ae9abeb691f0d2e0100000000000000003de0976428ce56125351bae62c5b8b8c79d8297c702ea05d60feabb4ed188b59c36fa759e93c0118b74b2618"
	block = ParseBlock(blockHeaderRaw)
	if block.Bip141() {
		t.Errorf("Expected teh block to NOT implement bip141 ")
	}
}

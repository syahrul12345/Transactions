package models

import (
	"encoding/hex"
	"testing"
)

func TestMerkleTreeInit(t *testing.T) {
	MerkleTree := CreateMerkleTree(9)
	wantlevel0 := 1
	getlevel0 := len(MerkleTree.Nodes[0])
	if wantlevel0 != getlevel0 {
		t.Errorf("Expected level 0 to have %d nodes, but got %d", wantlevel0, getlevel0)
	}
	wantLevel1 := 2
	getLevel1 := len(MerkleTree.Nodes[1])
	if wantLevel1 != getLevel1 {
		t.Errorf("Expected level 1 to have %d nodes, but got %d", wantLevel1, getLevel1)
	}
	wantLevel2 := 3
	getLevel2 := len(MerkleTree.Nodes[2])
	if wantLevel2 != wantLevel2 {
		t.Errorf("Expected level 2 to have %d nodes, but got %d", wantLevel2, getLevel2)
	}
	wantLevel3 := 5
	getLevel3 := len(MerkleTree.Nodes[3])
	if wantLevel3 != getLevel3 {
		t.Errorf("Expected level 3 to have %d nodes, but got %d", wantLevel3, getLevel3)
	}
	wantLevel4 := 9
	getLevel4 := len(MerkleTree.Nodes[4])
	if wantLevel4 != getLevel4 {
		t.Errorf("Expected level 4 to have %d nodes, but got %d", wantLevel4, getLevel4)
	}
}

func TestMerkleTreePopulateFirst(t *testing.T) {
	hexHashes := []string{
		"9745f7173ef14ee4155722d1cbf13304339fd00d900b759c6f9d58579b5765fb",
		"5573c8ede34936c29cdfdfe743f7f5fdfbd4f54ba0705259e62f39917065cb9b",
		"82a02ecbb6623b4274dfcab82b336dc017a27136e08521091e443e62582e8f05",
		"507ccae5ed9b340363a0e6d765af148be9cb1c8766ccc922f83e4ae681658308",
		"a7a4aec28e7162e1e9ef33dfa30f0bc0526e6cf4b11a576f6c5de58593898330",
		"bb6267664bd833fd9fc82582853ab144fece26b7a8a5bf328f8a059445b59add",
		"ea6d7ac1ee77fbacee58fc717b990c4fcccf1b19af43103c090f601677fd8836",
		"457743861de496c429912558a106b810b0507975a49773228aa788df40730d41",
		"7688029288efc9e9a0011c960a6ed9e5466581abf3e3a6c26ee317461add619a",
		"b1ae7f15836cb2286cdd4e2c37bf9bb7da0a2846d06867a429f654b2e7f383c9",
		"9b74f89fa3f93e71ff2c241f32945d877281a6a50a6bf94adac002980aafe5ab",
		"b3a92b5b255019bdaf754875633c2de9fec2ab03e6b8ce669d07cb5b18804638",
		"b5c0b915312b9bdaedd2b86aa2d0f8feffc73a2d37668fd9010179261e25e263",
		"c9d52c5cb1e557b92c84c52e7c4bfbce859408bedffc8a5560fd6e35e10b8800",
		"c555bc5fc3bc096df0a0c9532f07640bfb76bfe4fc1ace214b8b228a1297a4c2",
		"f9dbfafc3af3400954975da24eb325e326960a25b87fffe23eef3e7ed2fb610e",
	}
	// tree := CreateMerkleTree(uint64(len(hexHashes)))
	var flagBuf []byte
	for i := 0; i < 31; i++ {
		flagBuf = append(flagBuf, 1)
	}
	hashes32 := [][32]byte{}
	for _, hash := range hexHashes {
		var hash32buf [32]byte
		hashBytes, _ := hex.DecodeString(hash)
		copy(hash32buf[:], hashBytes)
		hashes32 = append(hashes32, hash32buf)
	}
	tree := CreateMerkleTree(uint64(len(hexHashes)))
	tree.Populate(flagBuf, hashes32)
	getRoot := tree.Root()
	wantRoot := "597c4bafe3832b17cbbabe56f878f4fc2ad0f6a402cee7fa851a9cb205f87ed1"
	if getRoot != wantRoot {
		t.Errorf("Expected the root to be %s but got %s", wantRoot, getRoot)
	}
}

func TestMerkleTreePopulate2(t *testing.T) {
	hexHashes := []string{
		"42f6f52f17620653dcc909e58bb352e0bd4bd1381e2955d19c00959a22122b2e",
		"94c3af34b9667bf787e1c6a0a009201589755d01d02fe2877cc69b929d2418d4",
		"959428d7c48113cb9149d0566bde3d46e98cf028053c522b8fa8f735241aa953",
		"a9f27b99d5d108dede755710d4a1ffa2c74af70b4ca71726fa57d68454e609a2",
		"62af110031e29de1efcad103b3ad4bec7bdcf6cb9c9f4afdd586981795516577",
	}
	var flagBuf []byte
	for i := 0; i < 11; i++ {
		flagBuf = append(flagBuf, 1)
	}
	hashes32 := [][32]byte{}
	for _, hash := range hexHashes {
		var hash32buf [32]byte
		hashBytes, _ := hex.DecodeString(hash)
		copy(hash32buf[:], hashBytes)
		hashes32 = append(hashes32, hash32buf)
	}
	tree := CreateMerkleTree(uint64(len(hexHashes)))
	tree.Populate(flagBuf, hashes32)
	getRoot := tree.Root()
	wantRoot := "a8e8bd023169b81bc56854137a135b97ef47a6a7237f4c6e037baed16285a5ab"
	if getRoot != wantRoot {
		t.Errorf("Expected the root to be %s but got %s", wantRoot, getRoot)
	}
}

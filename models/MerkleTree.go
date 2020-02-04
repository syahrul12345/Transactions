package models

import (
	"encoding/hex"
	"fmt"
	"math"
	"transactions/utils"
)

// MerkleTree represents the MerkleTree
type MerkleTree struct {
	// Total is the total number of nodes
	Total uint64
	// MaxDepth is the maximum level
	MaxDepth uint64
	//Nodes is an array of an array of strings. Each array of string, represetns ONE level in the merkle tree
	Nodes        [][]string
	CurrentDepth uint64
	CurrentIndex uint64
}

// CreateMerkleTree creates a merkle tree with the number of leaf nodes
func CreateMerkleTree(total uint64) *MerkleTree {
	// First we create an empty binary tree
	// We have log2total levels
	floatTotal := float64(total)
	MaxDepth := math.Ceil(math.Log2(floatTotal))
	Nodes := [][]string{}
	for depth := float64(0); depth < MaxDepth+1; depth++ {
		// Calculate the number of nodes in one level
		NumberOfNodes := math.Ceil(floatTotal / math.Pow(2, (MaxDepth-depth)))
		// Empty array to store the current level
		currentLevel := []string{}
		for count := float64(0); count < NumberOfNodes; count++ {
			currentLevel = append(currentLevel, "")
		}
		// Add it to the total level
		Nodes = append(Nodes, currentLevel)
	}
	return &MerkleTree{
		Total:        total,
		MaxDepth:     uint64(MaxDepth),
		Nodes:        Nodes,
		CurrentDepth: 0,
		CurrentIndex: 0,
	}
}

//Up Moves the pointer up
func (merkleTree *MerkleTree) Up() {
	merkleTree.CurrentDepth--
	merkleTree.CurrentIndex /= 2
}

//Left moves the pointer left
func (merkleTree *MerkleTree) Left() {
	merkleTree.CurrentDepth++
	merkleTree.CurrentIndex *= 2
}

//Right moves the pointer to the right
func (merkleTree *MerkleTree) Right() {
	merkleTree.CurrentDepth++
	merkleTree.CurrentIndex = merkleTree.CurrentIndex*2 + 1
}

//Root returns the root of hte merkle tree
func (merkleTree *MerkleTree) Root() string {
	return merkleTree.Nodes[0][0]
}

//SetCurrentNode sets a value to the current node
func (merkleTree *MerkleTree) SetCurrentNode(val string) {
	merkleTree.Nodes[merkleTree.CurrentDepth][merkleTree.CurrentIndex] = val
}

//GetCurrentNode gets the value of the current node
func (merkleTree *MerkleTree) GetCurrentNode() string {
	return merkleTree.Nodes[merkleTree.CurrentDepth][merkleTree.CurrentIndex]
}

//GetLeftNode returns the left node of the current node
func (merkleTree *MerkleTree) GetLeftNode() string {
	return merkleTree.Nodes[merkleTree.CurrentDepth+1][merkleTree.CurrentIndex*2]
}

//GetRightNode returns the left node of the current node
func (merkleTree *MerkleTree) GetRightNode() string {
	return merkleTree.Nodes[merkleTree.CurrentDepth+1][merkleTree.CurrentIndex*2+1]
}

//IsLeaf checks if hte current node is a leaf
func (merkleTree *MerkleTree) IsLeaf() bool {
	return merkleTree.CurrentDepth == merkleTree.MaxDepth
}

//RightExists check if there is a Right node to the current node
func (merkleTree *MerkleTree) RightExists() bool {
	return len(merkleTree.Nodes[merkleTree.CurrentDepth+1]) > int(merkleTree.CurrentIndex*2+1)
}

// Populate the merkle tree
func (merkleTree *MerkleTree) Populate(flagBits []byte, hashes [][32]byte) {
	for merkleTree.Root() == "" {
		if merkleTree.IsLeaf() {
			flagBits = flagBits[1:]
			hash := hashes[0]
			hashes = hashes[1:]
			merkleTree.SetCurrentNode(hex.EncodeToString(hash[:]))
			merkleTree.Up()
		} else {
			leftHash := merkleTree.GetLeftNode()
			// There is no left hash, so the hash can be in the hashes input or it might need to be calculated
			if leftHash == "" {
				// If 0, means it's in the hashes input
				flag := flagBits[0]
				flagBits = flagBits[1:]
				if flag == 0 {
					hash := hashes[0]
					hashes = hashes[1:]
					merkleTree.SetCurrentNode(hex.EncodeToString(hash[:]))
					// Pop from the flagbits
					merkleTree.Up()
				} else {
					// Move the pointer to the left child node, we will fill it up in the next iteration
					merkleTree.Left()
				}
			} else if merkleTree.RightExists() == true {
				// There is a left hash, we check if the right node exists
				rightHash := merkleTree.GetRightNode()
				// Check if the rightNode has a hash
				if rightHash == "" {
					// It doesnt have a hash, so we move the pointer here
					merkleTree.Right()
				} else {
					// We have the left hash (it wont enter this part of the code if there is no left hash) and the right hash, thus
					// we can calculate the hash of the current node
					leftHashBytes, _ := hex.DecodeString(leftHash)
					rightHashBytes, _ := hex.DecodeString(rightHash)
					merkleTree.SetCurrentNode(utils.MerkleParent(leftHashBytes, rightHashBytes))
					// We move up
					merkleTree.Up()
				}
			} else {
				// This part handles case where we have the left child node, but there is no right child node!
				leftHashBytes, _ := hex.DecodeString(leftHash)
				merkleTree.SetCurrentNode(utils.MerkleParent(leftHashBytes, leftHashBytes))
				merkleTree.Up()
			}
		}
	}
	if len(hashes) != 0 {
		fmt.Println("Not all hashes are consumed!")
	}
	for _, bit := range flagBits {
		if bit != 0 {
			fmt.Println("Not all bits are consumed")
		}
	}
}

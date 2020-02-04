package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"runtime"
	"strconv"
)

const (
	twoweeks uint64 = 60 * 60 * 24 * 14
)

//EncodeToLittleEndian : a large number represented from as an uint64 into hexadecimal litte-endian format in stirng represetnation.
//Input string mus be in base 10
func EncodeToLittleEndian(input uint64) string {
	// Some of the nubmers huge, use bigint
	someNumber := strconv.FormatUint(input, 10)
	bigInt, _ := big.NewInt(0).SetString(someNumber, 10)
	check1, _ := big.NewInt(0).SetString("0xfd"[2:], 16)
	check2, _ := big.NewInt(0).SetString("0x10000"[2:], 16)
	check3, _ := big.NewInt(0).SetString("0x100000000"[2:], 16)
	check4, _ := big.NewInt(0).SetString("0x10000000000000000"[2:], 16)
	if bigInt.Cmp(check1) < 0 {
		// If it's less than 256, it can be fit in one byte but
		return hex.EncodeToString(bigInt.Bytes())
	} else if bigInt.Cmp(check2) < 0 {
		numberInt, _ := strconv.ParseUint(someNumber, 10, 16)
		numberByte := make([]byte, 2)
		binary.LittleEndian.PutUint16(numberByte, uint16(numberInt))
		str := hex.EncodeToString(numberByte)
		return "fd" + str
	} else if bigInt.Cmp(check3) < 0 {
		numberInt, _ := strconv.ParseUint(someNumber, 10, 32)
		numberByte := make([]byte, 4)
		binary.LittleEndian.PutUint32(numberByte, uint32(numberInt))
		str := hex.EncodeToString(numberByte)
		return "fe" + str
	} else if bigInt.Cmp(check4) < 0 {
		numberInt, _ := strconv.ParseUint(someNumber, 10, 64)
		numberByte := make([]byte, 8)
		binary.LittleEndian.PutUint64(numberByte, numberInt)
		str := hex.EncodeToString(numberByte)
		return "ff" + str

	}
	return "FAILED"
}

//ReadVarInt will read the variable integer will decode the number of inputs from the incoming Transaction dump.
//It also returns the transaction dump without the version and input count
func ReadVarInt(txHash string) (uint64, string) {
	//Check the byte at position 9,10
	i := txHash[0:2]
	marker, _ := strconv.ParseUint(i, 16, 8)
	// Check the value of the marker
	if marker == 0xfd {
		//253 to 2^16-1
		//Take the next 2 bytes or 4 chars
		numberString := txHash[2:6]
		number, _ := strconv.ParseUint(numberString, 16, 16)
		//Create a empty byte array of the approriate size
		numberBytes := make([]byte, hex.DecodedLen(len(numberString)))
		//Put the parsed number in the byte array
		binary.BigEndian.PutUint16(numberBytes, uint16(number))
		//Create a new reader that reads the byte array filled with the bytes
		reader := bytes.NewReader(numberBytes)
		//Conver the byte array in little endian byte array
		var res uint16
		binary.Read(reader, binary.LittleEndian, &res)
		return uint64(res), txHash[6:]
	}
	if marker == 0xfe {
		//2^16 to 2^32 -1
		//take the next 4 bytes or 8 chars
		numberString := txHash[2:10]
		number, _ := strconv.ParseUint(numberString, 16, 32)
		//Create a byte array of the approriate size
		numberBytes := make([]byte, hex.DecodedLen(len(numberString)))
		//Put the parsed number in the byte array
		binary.BigEndian.PutUint32(numberBytes, uint32(number))
		//Create a reader
		reader := bytes.NewReader(numberBytes)
		//Conver the byte array to little endian
		var res uint32
		binary.Read(reader, binary.LittleEndian, &res)
		return uint64(res), txHash[10:]
	}
	if marker == 0xff {
		//For number between 2^32 and 2^54-1
		//Take the next 8 bytes or 16 chars
		numberString := txHash[2:18]
		number, _ := strconv.ParseUint(numberString, 16, 64)
		//Create an empoty byte array
		numberBytes := make([]byte, hex.DecodedLen(len(numberString)))
		//Put the number in the byte array as big endian
		binary.BigEndian.PutUint64(numberBytes, number)
		//Create a reader
		reader := bytes.NewReader(numberBytes)
		//Conver it to small endian
		var res uint64
		binary.Read(reader, binary.LittleEndian, &res)
		return uint64(res), txHash[18:]
	}
	//Marker itself is the number
	return uint64(marker), txHash[2:]
}

//ToBigHex will convert the string into the big-endian hexadecimal of the number, in string
//representation
func ToBigHex(input string) string {
	decodedHash, _ := hex.DecodeString(input)
	//reverse the decodedHash as it's actually in little endian
	for i := len(decodedHash)/2 - 1; i >= 0; i-- {
		opp := len(decodedHash) - 1 - i
		decodedHash[i], decodedHash[opp] = decodedHash[opp], decodedHash[i]
	}
	//Encode it to a string representation of the bytes in hexadecimal
	res := hex.EncodeToString(decodedHash)
	return res
}

//FromLittleHex will convert a hexadecimal stirng representing a number as a little endian to the correct number
func FromLittleHex(input string) uint32 {
	prevIndex, _ := hex.DecodeString(input)
	data := binary.LittleEndian.Uint32(prevIndex)
	return data
}

//GetFunctionName will get the function name
func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// Reverse will reverse a byte array
func Reverse(heightByte *[]byte) {
	a := *heightByte
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
	*heightByte = a
}

// BitsToTarget will return a human readable string representing the targert
func BitsToTarget(bits [4]byte) string {
	// Make a copy
	buf := make([]byte, 4)
	copy(buf, bits[:])

	exponentByte := buf[len(buf)-1] - 3
	coefficientBits := buf[:len(buf)-1]
	// Reverse it as it's little endian
	Reverse(&coefficientBits)
	coefficientString := hex.EncodeToString(coefficientBits)
	coefficient, _ := strconv.ParseUint(coefficientString, 16, 24)
	// First term
	firstTem := big.NewInt(0).SetUint64(coefficient)
	// Second term
	// 256 ** (exponent-3)
	secondTerm := big.NewInt(0).Exp(big.NewInt(256), big.NewInt(0).SetBytes([]byte{exponentByte}), big.NewInt(0))
	// Answer
	target := big.NewInt(0).Mul(firstTem, secondTerm)
	return target.Text(16)
}

// TargetToBits accepts a hexadecimal string and convert it to a byte array of size4
func TargetToBits(target string) [4]byte {
	rawBytes, _ := hex.DecodeString(target)
	var exponent byte
	var coefficient [3]byte
	if rawBytes[0] > 0x7f {
		// Exponent is the length of the rawbytes
		exponent = byte(len(rawBytes) + 1)
		coefficient = [3]byte{
			0x00,
			rawBytes[0],
			rawBytes[1],
		}
	} else {
		exponent = byte(len(rawBytes))
		coefficient = [3]byte{
			rawBytes[0],
			rawBytes[1],
			rawBytes[2],
		}
	}
	// Reverse it, it's in little endian
	buf := [4]byte{
		coefficient[2],
		coefficient[1],
		coefficient[0],
		exponent,
	}
	return buf
}

// CalculateNewBits will calcualte the new bits using the previousBits and the timedifferential
func CalculateNewBits(previousBits [4]byte, timeDifferential uint64) [4]byte {
	if timeDifferential > twoweeks*4 {
		timeDifferential = twoweeks * 4
	}
	if timeDifferential < twoweeks/4 {
		timeDifferential = timeDifferential / 4
	}
	previousTarget, _ := big.NewInt(0).SetString(BitsToTarget(previousBits), 16)
	timeDifferentialBig := big.NewInt(0).SetUint64(timeDifferential)
	twoWeeksBig := big.NewInt(0).SetUint64(twoweeks)
	// Time differential divided by 2 weeks
	first := big.NewInt(0).Mul(previousTarget, timeDifferentialBig)
	//  new target = previoustarget * timedifferentialdivided
	newTarget := big.NewInt(0).Div(first, twoWeeksBig)
	// previousTarget, _ := strconv.ParseUint(BitsToTarget(previousBits), 16, 32)
	return TargetToBits(newTarget.Text(16))
}

// MerkleParent calculates the parent of two hashes
func MerkleParent(hash0 []byte, hash1 []byte) string {
	buf := append(hash0, hash1...)
	first := sha256.Sum256(buf)
	second := sha256.Sum256(first[:])
	return hex.EncodeToString(second[:])
}

// MerkleParentLevel calculates the MerkleParents of every two hashes in a list of hahses, and stores the parents.
// Input should be a list of hashes in stirng
func MerkleParentLevel(hashes []string) []string {
	hashesBytes := [][]byte{}
	for _, hash := range hashes {
		hashByte, _ := hex.DecodeString(hash)
		hashesBytes = append(hashesBytes, hashByte)
	}
	// Check if it is an odd length
	if len(hashesBytes)%2 == 1 {
		lastObject := hashesBytes[len(hashesBytes)-1]
		hashesBytes = append(hashesBytes, lastObject)
	}
	parents := []string{}
	for i := 0; i < len(hashesBytes); i = i + 2 {
		parent := MerkleParent(hashesBytes[i], hashesBytes[i+1])
		parents = append(parents, parent)
	}
	return parents
}

// MerkleRoot returns the merkle root, a single hash from a set of hashes
func MerkleRoot(hashes []string) string {
	currentLevel := hashes
	for len(currentLevel) > 1 {
		currentLevel = MerkleParentLevel(currentLevel)
	}
	return currentLevel[0]
}

// BytesToBits converts a bytearray to a byte array of only 1 and 0
func BytesToBits(bytes []byte) []byte {
	var bitsBuf []byte
	for _, someByte := range bytes {
		for i := 0; i < 8; i++ {
			bitsBuf = append(bitsBuf, someByte&1)
			someByte >>= 1
		}
	}
	return bitsBuf
}

// BitsToBytes converts a byte array of only 1 to 0 to the a proper byte array
func BitsToBytes(bits []byte) *[]byte {
	if len(bits)%8 != 0 {
		fmt.Println("bit_field does not have a length that is divisible by 8")
		return nil
	}
	buf := make([]byte, len(bits)/8)
	for i, bit := range bits {
		byteIndex, bitIndex := i/8, i%8
		if bit == 1 {
			buf[byteIndex] |= 1 << bitIndex
		}
	}
	return &buf
}

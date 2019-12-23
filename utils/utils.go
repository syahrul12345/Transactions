package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"strconv"
)

//Decodes a string representing a hexadecimal encoded number in little-endian
//into a human-readable format. Accepts the target to be decoded
func FromLittle(hexadecimal string) {

}

//Encodes a large number represented from as an uint64 into hexadecimal litte-endian format in stirng represetnation.
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
		return bigInt.Text(16)
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

//GetInputs will decode the number of inputs from the incoming Transaction dump.
//It also returns the transaction dump without the version and input count
func GetInputs(txHash string) (uint64, string) {
	//Check the byte at position 9,10
	i := txHash[8:10]
	marker, _ := strconv.ParseUint(i, 16, 8)
	// Check the value of the marker
	if marker == 0xfd {
		//253 to 2^16-1
		//Take the next 2 bytes or 4 chars
		numberString := txHash[10:14]
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
		return uint64(res), txHash[14:]
	}
	if marker == 0xfe {
		//2^16 to 2^32 -1
		//take the next 4 bytes or 8 chars
		numberString := txHash[10:18]
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
		return uint64(res), txHash[18:]
	}
	if marker == 0xff {
		//For number between 2^32 and 2^54-1
		//Take the next 8 bytes or 16 chars
		numberString := txHash[10:26]
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
		return uint64(res), txHash[64:]
	}
	//Marker itself is the number
	return uint64(marker), txHash[10:]
}

//ToLittleHex will convert the byte array into the little-endian hexadecimal of the number, in string
//representation
func ToLittleHex(input string) string {
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

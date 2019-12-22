package models

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"transactions/utils"
)

type Transaction struct {
	txHash  string
	version uint32
}

//Parse a transaction from the given transactionhash
func Parse(txHash string) *Transaction {
	versionHash := txHash[0:8]
	version := GetVersion(versionHash)
	inputcount := GetInputs(txHash)
	fmt.Println(inputcount)
	res := utils.EncodeToLittleEndian(inputcount)
	fmt.Println(res)
	return &Transaction{
		txHash,
		version,
	}
}

//GetVersion will ge the version of the transaction hash. It takes the version hash as input
func GetVersion(versionHash string) uint32 {
	//Parse as a big endian integer
	versionInt, _ := strconv.ParseUint(versionHash, 16, 32)
	//Conver the integer to a byte array
	intBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(intBytes, uint32(versionInt))
	//Create a reader for the created byte array
	reader := bytes.NewReader(intBytes)
	//Define the result
	var res uint32
	binary.Read(reader, binary.LittleEndian, &res)
	return res

}

//GetInputs will decode the number of inputs from the string
func GetInputs(inputHash string) string {
	//Check the byte at position 9,10
	i := inputHash[8:10]
	marker, _ := strconv.ParseUint(i, 16, 8)
	// Check the value of the marker
	if marker == 0xfd {
		//253 to 2^16-1
		//Take the next 2 bytes or 4 chars
		numberString := inputHash[10:14]
		number, _ := strconv.ParseUint(numberString, 16, 16)
		//Create a empty byte array of the approriate size
		numberBytes := make([]byte, hex.DecodedLen(len(numberString)))
		//Put the parsed number in the byte array
		binary.BigEndian.PutUint16(numberBytes, uint16(number))
		//Create a new reader that reads the byte array filled with the bytes
		reader := bytes.NewReader(numberBytes)
		//Conver the byte array in little endian byte array
		var res int16
		binary.Read(reader, binary.LittleEndian, &res)
		return strconv.Itoa(int(res))
	}
	if marker == 0xfe {
		//2^16 to 2^32 -1
		//take the next 4 bytes or 8 chars
		numberString := inputHash[10:18]
		number, _ := strconv.ParseUint(numberString, 16, 32)
		//Create a byte array of the approriate size
		numberBytes := make([]byte, hex.DecodedLen(len(numberString)))
		//Put the parsed number in the byte array
		binary.BigEndian.PutUint32(numberBytes, uint32(number))
		//Create a reader
		reader := bytes.NewReader(numberBytes)
		//Conver the byte array to little endian
		var res int32
		binary.Read(reader, binary.LittleEndian, &res)
		return strconv.Itoa(int(res))
	}
	if marker == 0xff {
		//For number between 2^32 and 2^54-1
		//Take the next 8 bytes or 16 chars
		numberString := inputHash[10:26]
		number, _ := strconv.ParseUint(numberString, 16, 64)
		//Create an empoty byte array
		numberBytes := make([]byte, hex.DecodedLen(len(numberString)))
		//Put the number in the byte array as big endian
		binary.BigEndian.PutUint64(numberBytes, number)
		//Create a reader
		reader := bytes.NewReader(numberBytes)
		//Conver it to small endian
		var res int64
		binary.Read(reader, binary.LittleEndian, &res)
		return strconv.Itoa(int(res))
	}
	//Marker itself is the number
	return strconv.Itoa(int(marker))
}

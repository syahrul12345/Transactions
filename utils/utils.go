package utils

import (
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"strconv"
)

//Decodes a string representing a hexadecimal encoded number in little-endian
//into a human-readable format. Accepts the target to be decoded
func FromLittle(hexadecimal string) {

}

//Encodes a large number represented from a string into hexadecimal litte-endian format in stirng represetnation.
//Input string mus be in base 10
func EncodeToLittleEndian(someNumber string) string {
	// Some of the nubmers huge, use bigint
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

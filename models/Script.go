package models

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
)

//Script represents a list of commands
type Script struct {
	//Commands: A list of all commands to be executed by the script
	Commands []byte
}

//ParseScript will parse the hexadecimal string and return the corresponding script.
func ParseScript(s string) *Script {
	//Conver the string to a bytearray, whcih represents the bytearray
	commands := []byte{}
	byteHash, _ := hex.DecodeString(s)
	length := uint64(byteHash[0])
	byteHash = byteHash[1:]
	count := uint64(0)
	for count < length {
		//Get the current int
		currentInt := uint64(byteHash[0])
		byteHash = byteHash[1:]
		if currentInt >= 1 && currentInt <= 75 {
			//Read the next n element
			element := byteHash[0:currentInt]
			commands = append(commands, element...)
			byteHash = byteHash[currentInt:]
			count = count + currentInt
		} else if currentInt == 76 {
			//The next 1 byte is the length to be read
			nextByte := byteHash[0]
			byteHash = byteHash[1:]
			bytesToRead := uint64(nextByte)
			//Read the nextByte
			element := byteHash[0:bytesToRead]
			byteHash = byteHash[bytesToRead:]
			commands = append(commands, element...)
			count = count + bytesToRead
		} else if currentInt == 77 {
			nextByte := byteHash[0:2]
			byteHash = byteHash[2:]
			bytesToRead := binary.BigEndian.Uint64(nextByte)
			element := byteHash[0:bytesToRead]
			byteHash = byteHash[bytesToRead:]
			commands = append(commands, element...)
			count = count + bytesToRead
		} else {
			bigInt := big.NewInt(int64(currentInt))
			opCode := bigInt.Bytes()
			commands = append(commands, opCode...)
		}
	}
	if count != length {
		fmt.Println("failed to parse the script")
		return nil
	}
	return &Script{
		commands,
	}
}

//Serialize will serialize a script object and return the string
func (script *Script) Serialize() string {
	return "l"
}

package models

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"transactions/opcodes"
	"transactions/utils"
)

//Script represents a list of commands
type Script struct {
	//Commands: A list of all commands to be executed by the script
	//Each command is just a bytearray
	Commands [][]byte
}

//ParseScript will parse the hexadecimal string and return the corresponding script.
func ParseScript(s string) (*Script, []byte) {
	//Conver the string to a bytearray, whcih represents the scriptObject

	commands := [][]byte{}
	length, cleaned := utils.ReadVarInt(s)
	byteHash, _ := hex.DecodeString(cleaned)
	count := uint64(0)
	for count < length && len(byteHash) > 0 {
		//Get the current int
		currentInt := uint64(byteHash[0])
		byteHash = byteHash[1:]
		//Includes the current int into the count
		count = count + 1
		if currentInt >= 1 && currentInt <= 75 {
			//Read the next n element
			element := byteHash[0:currentInt]
			commands = append(commands, element)
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
			commands = append(commands, element)
			count = count + bytesToRead + 1
		} else if currentInt == 77 {
			nextByte := byteHash[0:2]
			byteHash = byteHash[2:]
			bytesToRead := binary.BigEndian.Uint64(nextByte)
			element := byteHash[0:bytesToRead]
			byteHash = byteHash[bytesToRead:]
			commands = append(commands, element)
			count = count + bytesToRead + 2
		} else {
			bigInt := big.NewInt(0).SetUint64(currentInt)
			opCode := bigInt.Bytes()
			commands = append(commands, opCode)
		}
	}
	if count != length {
		fmt.Println("failed to parse the script")
		return nil, []byte{}
	}
	return &Script{
		commands,
	}, byteHash
}

//Serialize will serialize a script object and return the string
func (script *Script) Serialize() string {
	raw := script.RawSerialize()
	return raw
}

//RawSerialize Serializes the DATA without the initial varint
func (script *Script) RawSerialize() string {
	var result string
	var totalLen int
	for _, scriptObj := range script.Commands {
		// Check if current scriptObj is an integer. It shoul d haveo nly 1 bytes
		totalLen = totalLen + len(scriptObj)
		if len(scriptObj) == 1 {
			number := int(scriptObj[0])
			numberHex := big.NewInt(int64(number)).Text(16)
			result = result + numberHex
		} else {
			if len(scriptObj) < 75 {
				len := len(scriptObj)
				lenHex := strconv.FormatInt(int64(len), 16)
				result = result + lenHex
				totalLen = totalLen + 1
			} else if len(scriptObj) > 75 && len(scriptObj) < 0x100 {
				// Encode the OP_PUSHDATA1, Encode the length
				lenHex := strconv.FormatInt(76, 16)
				result = result + lenHex
				//Encode the length
				result = result + strconv.FormatInt(int64(len(scriptObj)), 16)
				//The byte which handles the length should be added too
				totalLen = totalLen + 1
			} else if len(scriptObj) >= 0x100 && len(scriptObj) <= 520 {
				lenHex := strconv.FormatInt(77, 1)
				result = result + lenHex
				result = result + strconv.FormatInt(int64(len(scriptObj)), 16)
				totalLen = totalLen + 2
			} else {
				return "ERROR!"
			}
			command := hex.EncodeToString(scriptObj)
			result = result + command
		}
	}
	hex := utils.EncodeToLittleEndian(uint64(totalLen))
	return hex + result
}

//Add will add the command arrays of the two scripts object
func (script *Script) Add(script2 *Script) *Script {
	newCommands := append(script.Commands, script2.Commands...)
	return &Script{
		newCommands,
	}
}

//Evaluate will evaluate the commands in the script object, return true if succesfull
func (script *Script) Evaluate() bool {
	commands := script.Commands
	stack := [][]byte{}
	// altstack := [][]byte{}
	for len(commands) > 0 {
		// Get the first item and remove it
		command := commands[0]
		command = []byte{99}
		if len(command) <= 2 {
			// It's an OPCODE
			//convert to number
			buf := bytes.NewBuffer(command)
			number, _ := binary.ReadUvarint(buf)
			//Lets get the operation
			operation := opcodes.GetOPCODELIST()[int(number)]
			//Call the function
			operation.(func([][]byte, [][]byte) bool)(stack, commands)

		}
	}

	return true
}

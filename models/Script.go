package models

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"transactions/utils"
)

//Script represents a list of commands
type Script struct {
	//Commands: A list of all commands to be executed by the script
	//Each command is just a bytearray
	Commands [][]byte
}

//ParseScript will parse the hexadecimal string and return the corresponding script.
func ParseScript(s string) *Script {
	s = "6a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937"
	//Conver the string to a bytearray, whcih represents the scriptObject
	commands := [][]byte{}
	// byteHash, _ := hex.DecodeString(s)
	// length := uint64(byteHash[0])
	length, cleaned := utils.ReadVarInt(s)
	byteHash, _ := hex.DecodeString(cleaned)
	count := uint64(0)
	for count < length {
		//Get the current int
		currentInt := uint64(byteHash[0])
		byteHash = byteHash[1:]
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
			count = count + bytesToRead
		} else if currentInt == 77 {
			nextByte := byteHash[0:2]
			byteHash = byteHash[2:]
			bytesToRead := binary.BigEndian.Uint64(nextByte)
			element := byteHash[0:bytesToRead]
			byteHash = byteHash[bytesToRead:]
			commands = append(commands, element)
			count = count + bytesToRead
		} else {
			bigInt := big.NewInt(0).SetUint64(currentInt)
			opCode := bigInt.Bytes()
			commands = append(commands, opCode)
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
			} else if len(scriptObj) > 75 && len(scriptObj) < 0x100 {
				// Encode the OP_PUSHDATA1, Encode the length
				lenHex := strconv.FormatInt(76, 16)
				result = result + lenHex
				//Encode the length
				result = result + strconv.FormatInt(int64(len(scriptObj)), 16)
			} else if len(scriptObj) >= 0x100 && len(scriptObj) <= 520 {
				lenHex := strconv.FormatInt(77, 1)
				result = result + lenHex
				result = result + strconv.FormatInt(int64(len(scriptObj)), 16)
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

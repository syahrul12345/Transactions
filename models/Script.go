package models

import (
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
// Excepts Z, which is a pointer to a list of commands.
func (script *Script) Evaluate(z string) bool {
	commands := &script.Commands
	stack := &[][]byte{}
	altstack := &[][]byte{}
	// altstack := [][]byte{}
	for len(*commands) > 0 {
		// Get the first item and remove it
		tempCommands := *commands
		// Pop of the first element
		command := tempCommands[0]
		*commands = append(tempCommands[:0], tempCommands[1:]...)
		if len(command) == 1 {
			// Conver to number and get the function called operation
			number := command[0]
			operation := opcodes.GetOPCODELIST()[int(number)]
			if number == 99 || number == 100 {
				// It's false. Use type assertion to conver thte interface to a function
				if !operation.(func(*[][]byte, *[][]byte) bool)(stack, commands) {
					fmt.Printf("Bad op for %s\n", opcodes.GETOPCODENAMES()[int(number)])
					return false
				}
			}
			if number == 107 || number == 108 {
				if !operation.(func(*[][]byte, *[][]byte) bool)(stack, altstack) {
					fmt.Printf("Bad op for %s\n", opcodes.GETOPCODENAMES()[int(number)])
					return false
				}
			}
			if number == 172 || number == 173 || number == 174 || number == 175 {
				if !operation.(func(*[][]byte, string) bool)(stack, z) {
					fmt.Printf("Bad op for %s\n", opcodes.GETOPCODENAMES()[int(number)])
					return false
				}
			} else {
				if !operation.(func(*[][]byte) bool)(stack) {
					fmt.Printf("Bad op for %s\n", opcodes.GETOPCODENAMES()[int(number)])
					return false
				}
			}
		} else {
			*stack = append(*stack, command)
			// Special case for p2sh
			if len(*commands) == 3 {
				p2shcommands := *commands
				if p2shcommands[0][0] == 169 && len(p2shcommands[1]) == 20 && p2shcommands[2][0] == 135 {
					// pop of 135
					p2shcommands = p2shcommands[:len(p2shcommands)-1]
					// get the hash of the redeemscript from the command list
					h160 := p2shcommands[len(p2shcommands)-1]
					p2shcommands = p2shcommands[:len(p2shcommands)-1]
					// pop of 169
					p2shcommands = p2shcommands[:len(p2shcommands)-1]
					// Execute OP_HASH160 on the existing stack. This hashes the latest command which is a redeemscript.
					operation := opcodes.GetOPCODELIST()[169]
					if !operation.(func(*[][]byte) bool)(stack) {
						return false
					}
					// Add the hash of the redeem script to the top of the stack
					*stack = append(*stack, h160)
					operation = opcodes.GetOPCODELIST()[135]
					// Check if the top two objects in the stack a the same h160. If so this will encode a 1. DO this by calling op_equal
					if !operation.(func(*[][]byte) bool)(stack) {
						return false
					}
					// Check if the top is now a 1, showing that both h160 are the same
					operation = opcodes.GetOPCODELIST()[105]
					if !operation.(func(*[][]byte) bool)(stack) {
						fmt.Println("Bad p2sh h160")
						return false
					}
					// The top most item in the stack is the redeem script. This is also the command that was just added
					redeemScriptVarInt := utils.EncodeToLittleEndian(uint64(len(command)))
					redeemScriptString := hex.EncodeToString(command)
					redeemScript := redeemScriptVarInt + redeemScriptString
					// Create a script object for the redeemscript string
					redeemScriptObject, _ := ParseScript(redeemScript)
					redeemScriptCommands := redeemScriptObject.Commands
					// Add redeem script commands to the comamnd list
					*commands = append(p2shcommands, redeemScriptCommands...)
				}
			}
		}
	}

	if len(*stack) == 0 {
		return false
	}
	tempStack := *stack
	lastItem := tempStack[len(*stack)-1]
	//If it's an empty byte array, return false. This means that iti is an empty bytestring
	if len(lastItem) == 0 {
		return false
	}
	return true
}

// ISP2PKH : Checks if the script object follows the p2pkh pattern
func (script *Script) ISP2PKH() bool {
	// Check if the command list follow the pattern of OP_DUP OP_HASH160 <20 byte hash> OP_EQUALVERIFY OP_CHECKSIG
	commands := script.Commands
	return len(commands) == 5 && commands[0][0] == 0x76 && commands[1][0] == 0xa9 && len(commands[2]) == 20 && commands[3][0] == 0x88 && commands[4][0] == 0xac
}

// ISP2SH : Checks if the script object follows the p2sh pattern
func (script *Script) ISP2SH() bool {
	//Check if the command list OP_HASH160 <20 byte hash> OP_EQUAL pattern.'''
	commands := script.Commands
	return len(commands) == 3 && commands[0][0] == 0xa9 && len(commands[1]) == 20 && commands[2][0] == 0x87
}

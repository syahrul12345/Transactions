package opcodes

import (
	"bytes"
	"encoding/binary"
	"math"
)

//GetOPCODELIST returns a mapping representing the OPCODELIST
func GetOPCODELIST() map[int]interface{} {
	OPCODELIST := map[int]interface{}{
		0:   op_0,
		79:  op_1negate,
		81:  op_1,
		82:  op_2,
		83:  op_3,
		84:  op_4,
		85:  op_5,
		86:  op_6,
		87:  op_7,
		88:  op_8,
		89:  op_9,
		90:  op_10,
		91:  op_11,
		92:  op_12,
		93:  op_13,
		94:  op_14,
		95:  op_15,
		96:  op_16,
		97:  op_nop,
		99:  op_if,
		100: op_notif,
		105: op_verify,
		106: op_return,
		107: op_toaltstack,
		108: op_fromaltstack,
		109: op_2drop,
		110: op_2dup,
		111: op_3dup,
		112: op_2over,
		// 113: op_2rot,
		// 114: op_2swap,
		// 115: op_ifdup,
		// 116: op_depth,
		// 117: op_drop,
		// 118: op_dup,
		// 119: op_nip,
		// 120: op_over,
		// 121: op_pick,
		// 122: op_roll,
		// 123: op_rot,
		// 124: op_swap,
		// 125: op_tuck,
		// 130: op_size,
		// 135: op_equal,
		// 136: op_equalverify,
		// 139: op_1add,
		// 140: op_1sub,
		// 143: op_negate,
		// 144: op_abs,
		// 145: op_not,
		// 146: op_0notequal,
		// 147: op_add,
		// 148: op_sub,
		// 149: op_mul,
		// 154: op_booland,
		// 155: op_boolor,
		// 156: op_numequal,
		// 157: op_numequalverify,
		// 158: op_numnotequal,
		// 159: op_lessthan,
		// 160: op_greaterthan,
		// 161: op_lessthanorequal,
		// 162: op_greaterthanorequal,
		// 163: op_min,
		// 164: op_max,
		// 165: op_within,
		// 166: op_ripemd160,
		// 167: op_sha1,
		// 168: op_sha256,
		// 169: op_hash160,
		// 170: op_hash256,
		// 172: op_checksig,
		// 173: op_checksigverify,
		// 174: op_checkmultisig,
		// 175: op_checkmultisigverify,
		// 176: op_nop,
		// 177: op_checklocktimeverify,
		// 178: op_checksequenceverify,
		// 179: op_nop,
		// 180: op_nop,
		// 181: op_nop,
		// 182: op_nop,
		// 183: op_nop,
		// 184: op_nop,
		// 185: op_nop,
	}
	return OPCODELIST
}

func encodeNum(num int) []byte {
	//Return empty byte
	if num == 0 {
		return []byte{}
	}
	absNum := int(math.Abs(float64(num)))
	negative := num < 0
	res := []byte{}
	for absNum > 0 {
		xor := absNum & 0xff
		res = append(res, byte(xor))
		absNum >>= 8
	}
	if res[len(res)-1]&0x80 == 1 {
		if negative {
			res = append(res, 0x80)
		} else {
			res = append(res, 0)
		}
	} else if negative {
		res[len(res)-1] |= 0x80
	}
	return res
}

func decodeNum(element []byte) int64 {
	//Empty byte arrya
	var negative bool
	var result int64
	if len(element) == 0 {
		return 0
	}
	// Reverse the element ot make it into big endian
	for i := len(element)/2 - 1; i >= 0; i-- {
		opp := len(element) - 1 - i
		element[i], element[opp] = element[opp], element[i]
	}
	// It's one
	if element[0]&0x80 == 1 {
		negative = true
		result = int64(element[0]) & 0x7f
	} else {
		negative = false
		result = int64(element[0])
	}
	for _, someByte := range element[1:] {
		result <<= 8
		result = result + int64(someByte)
	}
	if negative {
		return -(result)
	}
	return result
}

func op_0(commands [][]byte) bool {
	commands = append(commands, encodeNum(0))
	return true
}
func op_1negate(commands [][]byte) bool {
	commands = append(commands, encodeNum(-1))
	return true
}

func op_1(commands [][]byte) bool {
	commands = append(commands, encodeNum(1))
	return true
}
func op_2(commands [][]byte) bool {
	commands = append(commands, encodeNum(2))
	return true
}
func op_3(commands [][]byte) bool {
	commands = append(commands, encodeNum(3))
	return true
}
func op_4(commands [][]byte) bool {
	commands = append(commands, encodeNum(4))
	return true
}
func op_5(commands [][]byte) bool {
	commands = append(commands, encodeNum(5))
	return true
}
func op_6(commands [][]byte) bool {
	commands = append(commands, encodeNum(6))
	return true
}
func op_7(commands [][]byte) bool {
	commands = append(commands, encodeNum(7))
	return true
}
func op_8(commands [][]byte) bool {
	commands = append(commands, encodeNum(8))
	return true
}
func op_9(commands [][]byte) bool {
	commands = append(commands, encodeNum(9))
	return true
}
func op_10(commands [][]byte) bool {
	commands = append(commands, encodeNum(10))
	return true
}
func op_11(commands [][]byte) bool {
	commands = append(commands, encodeNum(11))
	return true
}
func op_12(commands [][]byte) bool {
	commands = append(commands, encodeNum(12))
	return true
}
func op_13(commands [][]byte) bool {
	commands = append(commands, encodeNum(13))
	return true
}
func op_14(commands [][]byte) bool {
	commands = append(commands, encodeNum(14))
	return true
}
func op_15(commands [][]byte) bool {
	commands = append(commands, encodeNum(15))
	return true
}
func op_16(commands [][]byte) bool {
	commands = append(commands, encodeNum(16))
	return true
}
func op_nop(commands [][]byte) bool {
	return true
}
func op_if(stack, commands [][]byte) bool {
	if len(stack) < 1 {
		return false
	}
	trueItems := [][]byte{}
	falseItems := [][]byte{}
	currentArray := trueItems
	found := false
	numberEndifsNeeded := 1
	for len(commands) > 0 {
		command := commands[0]
		commands = commands[1:]
		// Conver the current coimmand into a number,it's just one byte
		buf := bytes.NewBuffer(command)
		number, _ := binary.ReadUvarint(buf)
		if number == 99 || number == 100 {
			numberEndifsNeeded = numberEndifsNeeded + 1
		} else if numberEndifsNeeded == 1 && number == 103 {
			currentArray = falseItems
		} else if number == 104 {
			if numberEndifsNeeded == 1 {
				found = true
				break
			} else {
				numberEndifsNeeded = numberEndifsNeeded - 1
				currentArray = append(currentArray, command)
			}
		} else {
			currentArray = append(currentArray, command)
		}
	}
	if !found {
		return false
	}
	element := commands[len(commands)-1]
	commands = commands[:len(commands)-1]
	buf := bytes.NewBuffer(element)
	number, _ := binary.ReadUvarint(buf)
	//Append trueitems or falseitems to the front of command
	if number == 0 {
		commands = append(trueItems, commands...)
	} else {
		commands = append(falseItems, commands...)
	}
	return true
}
func op_notif(stack, commands [][]byte) bool {
	if len(stack) < 1 {
		return true
	}
	trueItems := [][]byte{}
	falseItems := [][]byte{}
	currentArray := trueItems
	found := false
	numEndIfsNeeded := 1
	for len(commands) > 1 {
		// Get the last byte array. Let's convert it to a number
		lastCommand := commands[len(commands)-1]
		// It should be one byte for op_codes
		if len(lastCommand) == 1 {
			if int(lastCommand[0]) == 99 || int(lastCommand[0]) == 100 {
				numEndIfsNeeded = numEndIfsNeeded + 1
				currentArray = append(currentArray, lastCommand)
			} else if numEndIfsNeeded == 1 && int(lastCommand[0]) == 103 {
				currentArray = falseItems
			} else if int(lastCommand[0]) == 104 {
				if numEndIfsNeeded == 1 {
					found = true
					break
				} else {
					numEndIfsNeeded = numEndIfsNeeded - 1
					currentArray = append(currentArray, lastCommand)
				}
			}
		}

	}
	if !found {
		return false
	}
	element := stack[len(stack)-1]
	if decodeNum(element) == 0 {
		commands = append(trueItems, commands...)
	} else {
		commands = append(falseItems, commands...)
	}
	return true
}

func op_verify(stack, commands [][]byte) bool {
	if len(stack) < 1 {
		return false
	}
	last := commands[len(commands)-1]
	if decodeNum(last) == 0 {
		return false
	}
	return true
}

func op_return(stack [][]byte) bool {
	return false
}

func op_toaltstack(stack, altstack [][]byte) bool {
	if len(stack) < 1 {
		return false
	}
	altstack = append(altstack, stack[len(stack)-1])
	return true
}

func op_fromaltstack(stack, alstack [][]byte) bool {
	if len(alstack) < 1 {
		return false
	}
	stack = append(stack, alstack[len(alstack)-1])
	return true
}
func op_2drop(stack [][]byte) bool {
	if len(stack) < 2 {
		return false
	}
	stack = stack[:len(stack)-1]
	stack = stack[:len(stack)-1]
	return true
}
func op_2dup(stack [][]byte) bool {
	if len(stack) < 2 {
		return false
	}
	stack = append(stack, stack[len(stack)-2:]...)
	return true
}
func op_3dup(stack [][]byte) bool {
	if len(stack) < 3 {
		return false
	}
	stack = append(stack, stack[len(stack)-3:]...)
	return true
}
func op_2over(stack [][]byte) bool {
	if len(stack) < 4 {
		return false
	}
	stack = append(stack, stack[len(stack)-4:len(stack)-2]...)
	return true
}

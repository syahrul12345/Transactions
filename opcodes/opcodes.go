package opcodes

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"

	"github.com/syahrul12345/secp256k1"
	"golang.org/x/crypto/ripemd160"
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
		113: op_2rot,
		114: op_2swap,
		115: op_ifdup,
		116: op_depth,
		117: op_drop,
		118: op_dup,
		119: op_nip,
		120: op_over,
		121: op_pick,
		122: op_roll,
		123: op_rot,
		124: op_swap,
		125: op_tuck,
		130: op_size,
		135: op_equal,
		136: op_equalverify,
		139: op_1add,
		140: op_1sub,
		143: op_negate,
		144: op_abs,
		145: op_not,
		146: op_0notequal,
		147: op_add,
		148: op_sub,
		149: op_mul,
		154: op_booland,
		155: op_boolor,
		156: op_numequal,
		157: op_numequalverify,
		158: op_numnotequal,
		159: op_lessthan,
		160: op_greaterthan,
		161: op_lessthanorequal,
		162: op_greaterthanorequal,
		163: op_min,
		164: op_max,
		165: op_within,
		166: op_ripemd160,
		167: op_sha1,
		168: op_sha256,
		169: op_hash160,
		170: op_hash256,
		172: op_checksig,
		173: op_checksigverify,
		174: op_checkmultisig,
		175: op_checkmultisigverify,
		176: op_nop,
		177: op_checklocktimeverify,
		178: op_checksequenceverify,
		179: op_nop,
		180: op_nop,
		181: op_nop,
		182: op_nop,
		183: op_nop,
		184: op_nop,
		185: op_nop,
	}
	return OPCODELIST
}

//GETOPCODENAMES will return the opcode names
func GETOPCODENAMES() map[int]string {
	OPCODELIST := map[int]string{
		0:   "op_0",
		79:  "op_1negate",
		81:  "op_1",
		82:  "op_2",
		83:  "op_3",
		84:  "op_4",
		85:  "op_5",
		86:  "op_6",
		87:  "op_7",
		88:  "op_8",
		89:  "op_9",
		90:  "op_10",
		91:  "op_11",
		92:  "op_12",
		93:  "op_13",
		94:  "op_14",
		95:  "op_15",
		96:  "op_16",
		97:  "op_nop",
		99:  "op_if",
		100: "op_notif",
		105: "op_verify",
		106: "op_return",
		107: "op_toaltstack",
		108: "op_fromaltstack",
		109: "op_2drop",
		110: "op_2dup",
		111: "op_3dup",
		112: "op_2over",
		113: "op_2rot",
		114: "op_2swap",
		115: "op_ifdup",
		116: "op_depth",
		117: "op_drop",
		118: "op_dup",
		119: "op_nip",
		120: "op_over",
		121: "op_pick",
		122: "op_roll",
		123: "op_rot",
		124: "op_swap",
		125: "op_tuck",
		130: "op_size",
		135: "op_equal",
		136: "op_equalverify",
		139: "op_1add",
		140: "op_1sub",
		143: "op_negate",
		144: "op_abs",
		145: "op_not",
		146: "op_0notequal",
		147: "op_add",
		148: "op_sub",
		149: "op_mul",
		154: "op_booland",
		155: "op_boolor",
		156: "op_numequal",
		157: "op_numequalverify",
		158: "op_numnotequal",
		159: "op_lessthan",
		160: "op_greaterthan",
		161: "op_lessthanorequal",
		162: "op_greaterthanorequal",
		163: "op_min",
		164: "op_max",
		165: "op_within",
		166: "op_ripemd160",
		167: "op_sha1",
		168: "op_sha256",
		169: "op_hash160",
		170: "op_hash256",
		172: "op_checksig",
		173: "op_checksigverify",
		174: "op_checkmultisig",
		175: "op_checkmultisigverify",
		176: "op_nop",
		177: "op_checklocktimeverify",
		178: "op_checksequenceverify",
		179: "op_nop",
		180: "op_nop",
		181: "op_nop",
		182: "op_nop",
		183: "op_nop",
		184: "op_nop",
		185: "op_nop",
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

func op_0(commands *[][]byte) bool {
	*commands = append(*commands, encodeNum(0))
	return true
}
func op_1negate(commands *[][]byte) bool {
	*commands = append(*commands, encodeNum(-1))
	return true
}

func op_1(commands *[][]byte) bool {
	*commands = append(*commands, encodeNum(1))
	return true
}
func op_2(commands *[][]byte) bool {
	*commands = append(*commands, encodeNum(2))
	return true
}
func op_3(commands *[][]byte) bool {
	*commands = append(*commands, encodeNum(3))
	return true
}
func op_4(commands *[][]byte) bool {
	*commands = append(*commands, encodeNum(4))
	return true
}
func op_5(commands *[][]byte) bool {
	*commands = append(*commands, encodeNum(5))
	return true
}
func op_6(commands *[][]byte) bool {
	*commands = append(*commands, encodeNum(6))
	return true
}
func op_7(commands *[][]byte) bool {
	*commands = append(*commands, encodeNum(7))
	return true
}
func op_8(commands *[][]byte) bool {
	*commands = append(*commands, encodeNum(8))
	return true
}
func op_9(commands *[][]byte) bool {
	*commands = append(*commands, encodeNum(9))
	return true
}
func op_10(commands *[][]byte) bool {
	*commands = append(*commands, encodeNum(10))
	return true
}
func op_11(commands *[][]byte) bool {
	*commands = append(*commands, encodeNum(11))
	return true
}
func op_12(commands *[][]byte) bool {
	*commands = append(*commands, encodeNum(12))
	return true
}
func op_13(commands *[][]byte) bool {
	*commands = append(*commands, encodeNum(13))
	return true
}
func op_14(commands *[][]byte) bool {
	*commands = append(*commands, encodeNum(14))
	return true
}
func op_15(commands *[][]byte) bool {
	*commands = append(*commands, encodeNum(15))
	return true
}
func op_16(commands *[][]byte) bool {
	*commands = append(*commands, encodeNum(16))
	return true
}
func op_nop(commands *[][]byte) bool {
	return true
}
func op_if(stack, commands *[][]byte) bool {
	if len(*stack) < 1 {
		return false
	}
	trueItems := [][]byte{}
	falseItems := [][]byte{}
	currentArray := trueItems
	found := false
	numberEndifsNeeded := 1
	for len(*commands) > 0 {
		tempCommands := *commands
		command := tempCommands[0]
		*commands = tempCommands[1:]
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
	tempCommands := *commands
	element := tempCommands[len(tempCommands)-1]
	*commands = tempCommands[:len(tempCommands)-1]
	buf := bytes.NewBuffer(element)
	number, _ := binary.ReadUvarint(buf)
	//Append trueitems or falseitems to the front of command
	if number == 0 {
		*commands = append(trueItems, *commands...)
	} else {
		*commands = append(falseItems, *commands...)
	}
	return true
}
func op_notif(stack, commands *[][]byte) bool {
	if len(*stack) < 1 {
		return true
	}
	trueItems := [][]byte{}
	falseItems := [][]byte{}
	currentArray := trueItems
	found := false
	numEndIfsNeeded := 1
	for len(*commands) > 1 {
		// Get the last byte array. Let's convert it to a number
		tempCommands := *commands
		lastCommand := tempCommands[len(tempCommands)-1]
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
	tempStack := *stack
	element := tempStack[len(tempStack)-1]
	if decodeNum(element) == 0 {
		*commands = append(trueItems, *commands...)
	} else {
		*commands = append(falseItems, *commands...)
	}
	return true
}

func op_verify(stack *[][]byte) bool {
	if len(*stack) < 1 {
		return false
	}
	tempStack := *stack
	element := tempStack[len(tempStack)-1]
	*stack = tempStack[:len(tempStack)-1]
	if decodeNum(element) == 0 {
		return false
	}
	return true
}
func op_return(stack [][]byte) bool {
	return false
}

func op_toaltstack(stack, altstack *[][]byte) bool {
	if len(*stack) < 1 {
		return false
	}
	tempStack := *stack
	*altstack = append(*altstack, tempStack[len(tempStack)-1])
	return true
}

func op_fromaltstack(stack, alstack *[][]byte) bool {
	if len(*alstack) < 1 {
		return false
	}
	tempStack := *alstack
	*stack = append(*stack, tempStack[len(tempStack)-1])
	return true
}
func op_2drop(stack *[][]byte) bool {
	if len(*stack) < 2 {
		return false
	}
	tempStack := *stack
	// remove the last element of the stack
	tempStack = tempStack[:len(tempStack)]
	tempStack = tempStack[:len(tempStack)]
	*stack = tempStack
	return true
}
func op_2dup(stack *[][]byte) bool {
	if len(*stack) < 2 {
		return false
	}
	tempStack := *stack
	*stack = append(*stack, tempStack[len(tempStack)-2:]...)
	return true
}
func op_3dup(stack *[][]byte) bool {
	if len(*stack) < 3 {
		return false
	}
	tempStack := *stack
	*stack = append(*stack, tempStack[len(tempStack)-3:]...)
	return true
}
func op_2over(stack *[][]byte) bool {
	if len(*stack) < 4 {
		return false
	}
	tempStack := *stack
	*stack = append(*stack, tempStack[len(tempStack)-4:len(tempStack)-2]...)
	return true
}
func op_2rot(stack *[][]byte) bool {
	if len(*stack) < 6 {
		return false
	}
	tempStack := *stack
	*stack = append(*stack, tempStack[len(tempStack)-6:len(tempStack)-4]...)
	return true
}
func op_2swap(stack *[][]byte) bool {
	if len(*stack) < 4 {
		return false
	}
	tempStack := *stack
	tempStackFirst := tempStack[len(tempStack)-2:]
	tempStackSecond := tempStack[len(tempStack)-4 : len(tempStack)-2]
	tempStackReplace := append(tempStackFirst, tempStackSecond...)
	tempStack = append(tempStack[:len(tempStack)-3], tempStackReplace...)
	*stack = tempStack
	return true
}

func op_ifdup(stack *[][]byte) bool {
	if len(*stack) < 1 {
		return false
	}
	tempStack := *stack
	if decodeNum(tempStack[len(tempStack)-1]) != 0 {
		*stack = append(*stack, tempStack[len(tempStack)-1])
	}
	return true
}

func op_depth(stack *[][]byte) bool {
	*stack = append(*stack, encodeNum(len(*stack)))
	return true
}
func op_drop(stack *[][]byte) bool {
	if len(*stack) < 1 {
		return false
	}
	tempStack := *stack
	tempStack = tempStack[:len(tempStack)-1]
	*stack = tempStack
	return true
}
func op_dup(stack *[][]byte) bool {
	if len(*stack) < 1 {
		return false
	}
	tempStack := *stack
	*stack = append(*stack, tempStack[len(tempStack)-1])
	return true
}
func op_nip(stack *[][]byte) bool {
	if len(*stack) < 2 {
		return false
	}
	tempStack := *stack
	tempStack = append(tempStack[:len(tempStack)-2], tempStack[len(tempStack)-1:]...)
	return true
}
func op_over(stack *[][]byte) bool {
	if len(*stack) < 2 {
		return false
	}
	tempStack := *stack
	*stack = append(*stack, tempStack[len(tempStack)-2])
	return true
}

func op_pick(stack *[][]byte) bool {
	if len(*stack) < 1 {
		return false
	}
	tempStack := *stack
	n := decodeNum(tempStack[len(tempStack)-1])
	*stack = tempStack[:len(tempStack)-1]
	if len(*stack) < int(n+1) {
		return false
	}
	*stack = append(*stack, tempStack[len(tempStack)-int(n+1)])
	return true
}

func op_roll(stack *[][]byte) bool {
	if len(*stack) < 1 {
		return false
	}
	tempStack := *stack
	n := decodeNum(tempStack[len(tempStack)-1])
	*stack = tempStack[:len(tempStack)-1]
	if len(*stack) < int(n+1) {
		return false
	}
	if n == 0 {
		return false
	}
	*stack = append(*stack, tempStack[len(tempStack)-int(n+1)])
	return false
}

func op_rot(stack *[][]byte) bool {
	if len(*stack) < 3 {
		return false
	}
	tempStack := *stack
	*stack = tempStack[:len(tempStack)-3]
	*stack = append(*stack, tempStack[len(tempStack)-3])
	return true
}

func op_swap(stack *[][]byte) bool {
	if len(*stack) < 2 {
		return false
	}
	tempStack := *stack
	*stack = tempStack[:len(tempStack)-2]
	*stack = append(*stack, tempStack[len(tempStack)-2])
	return false
}

func op_tuck(stack *[][]byte) bool {
	if len(*stack) < 2 {
		return false
	}
	tempStack := *stack
	// Insert the last stack item into the second last position.
	i := len(tempStack) - 2
	*stack = append(tempStack[:i], append(tempStack[:len(tempStack)-1], tempStack[i:]...)...)
	return false
}

func op_size(stack *[][]byte) bool {
	if len(*stack) < 1 {
		return false
	}
	tempStack := *stack
	// Encodes the length of the last element in the stack
	num := encodeNum(len(tempStack[len(tempStack)-1]))
	*stack = append(*stack, num)
	return true
}

func op_equal(stack *[][]byte) bool {
	if len(*stack) < 2 {
		return false
	}
	tempStack := *stack
	//Stack.pop()
	element1 := tempStack[len(tempStack)-1]
	tempStack = tempStack[:len(tempStack)-1]
	element2 := tempStack[len(tempStack)-1]
	*stack = tempStack[:len(tempStack)-1]
	// Since element1 and 2 represents a number
	num1 := decodeNum(element1)
	num2 := decodeNum(element2)
	if num1 == num2 {
		*stack = append(*stack, encodeNum(1))
	} else {
		*stack = append(*stack, encodeNum(0))
	}
	return true
}
func op_equalverify(stack *[][]byte) bool {
	return op_equal(stack) && op_verify(stack)
}

func op_1add(stack *[][]byte) bool {
	if len(*stack) < 1 {
		return false
	}
	tempStack := *stack
	element := decodeNum(tempStack[len(tempStack)-1])
	*stack = tempStack[:len(tempStack)-1]
	*stack = append(*stack, encodeNum(int(element+1)))
	return true
}

func op_1sub(stack *[][]byte) bool {
	if len(*stack) < 1 {
		return false
	}
	tempStack := *stack
	element := decodeNum(tempStack[len(tempStack)-1])
	*stack = tempStack[:len(tempStack)-1]
	*stack = append(*stack, encodeNum(int(element-1)))
	return true
}

func op_negate(stack *[][]byte) bool {
	if len(*stack) < 1 {
		return false
	}
	tempStack := *stack
	element := decodeNum(tempStack[len(tempStack)-1])
	*stack = tempStack[:len(tempStack)-1]
	*stack = append(*stack, encodeNum(int(-element)))
	return false
}

func op_abs(stack *[][]byte) bool {
	if len(*stack) < 1 {
		return false
	}
	tempStack := *stack
	element := decodeNum(tempStack[len(tempStack)-1])
	*stack = tempStack[:len(tempStack)-1]
	if element < 0 {
		*stack = append(*stack, encodeNum(int(-element)))
	} else {
		*stack = append(*stack, encodeNum(int(element)))
	}
	return true
}

func op_not(stack *[][]byte) bool {
	if len(*stack) < 1 {
		return false
	}
	tempStack := *stack
	element := tempStack[len(tempStack)-1]
	*stack = tempStack[:len(tempStack)-1]
	if decodeNum(element) == 0 {
		*stack = append(*stack, encodeNum(1))
	} else {
		*stack = append(*stack, encodeNum(0))
	}
	return true

}

func op_0notequal(stack *[][]byte) bool {
	if len(*stack) < 1 {
		return false
	}
	tempStack := *stack
	element := tempStack[len(tempStack)-1]
	*stack = tempStack[:len(tempStack)-1]
	if decodeNum(element) == 0 {
		*stack = append(*stack, encodeNum(0))
	} else {
		*stack = append(*stack, encodeNum(1))
	}
	return true
}

func op_add(stack *[][]byte) bool {
	if len(*stack) < 2 {
		return false
	}
	tempStack := *stack
	element1 := decodeNum(tempStack[len(tempStack)-1])
	tempStack = tempStack[:len(tempStack)-1]
	element2 := decodeNum(tempStack[len(tempStack)-1])
	*stack = tempStack[:len(tempStack)-1]
	combined := element1 + element2
	*stack = append(*stack, encodeNum(int(combined)))
	return true
}

func op_sub(stack *[][]byte) bool {
	if len(*stack) < 2 {
		return false
	}
	tempStack := *stack
	element1 := decodeNum(tempStack[len(tempStack)-1])
	tempStack = tempStack[:len(tempStack)-1]
	element2 := decodeNum(tempStack[len(tempStack)-1])
	*stack = tempStack[:len(tempStack)-1]
	combined := element2 - element1
	*stack = append(*stack, encodeNum(int(combined)))
	return true
}

func op_mul(stack *[][]byte) bool {
	if len(*stack) < 2 {
		return false
	}
	tempStack := *stack
	element1 := decodeNum(tempStack[len(tempStack)-1])
	tempStack = tempStack[:len(tempStack)-1]
	element2 := decodeNum(tempStack[len(tempStack)-1])
	*stack = tempStack[:len(tempStack)-1]
	combined := element2 * element1
	*stack = append(*stack, encodeNum(int(combined)))
	return true
}
func op_booland(stack *[][]byte) bool {
	if len(*stack) < 2 {
		return false
	}
	tempStack := *stack
	element1 := decodeNum(tempStack[len(tempStack)-1])
	tempStack = tempStack[:len(tempStack)-1]
	element2 := decodeNum(tempStack[len(tempStack)-1])
	*stack = tempStack[:len(tempStack)-1]
	if element1 > 0 && element2 > 0 {
		*stack = append(*stack, encodeNum(1))
	} else {
		*stack = append(*stack, encodeNum(0))
	}
	return true
}

func op_boolor(stack *[][]byte) bool {
	if len(*stack) < 2 {
		return false
	}
	tempStack := *stack
	element1 := decodeNum(tempStack[len(tempStack)-1])
	tempStack = tempStack[:len(tempStack)-1]
	element2 := decodeNum(tempStack[len(tempStack)-1])
	*stack = tempStack[:len(tempStack)-1]
	if element1 > 0 || element2 > 0 {
		*stack = append(*stack, encodeNum(1))
	} else {
		*stack = append(*stack, encodeNum(0))
	}
	return true
}

func op_numequal(stack *[][]byte) bool {
	if len(*stack) < 2 {
		return false
	}
	tempStack := *stack
	element1 := decodeNum(tempStack[len(tempStack)-1])
	tempStack = tempStack[:len(tempStack)-1]
	element2 := decodeNum(tempStack[len(tempStack)-1])
	*stack = tempStack[:len(tempStack)-1]
	if element1 == element2 {
		*stack = append(*stack, encodeNum(1))
	} else {
		*stack = append(*stack, encodeNum(0))
	}
	return true
}

func op_numequalverify(stack *[][]byte) bool {
	return op_numequal(stack) && op_verify(stack)
}

func op_numnotequal(stack *[][]byte) bool {
	if len(*stack) < 2 {
		return false
	}
	tempStack := *stack
	element1 := decodeNum(tempStack[len(tempStack)-1])
	tempStack = tempStack[:len(tempStack)-1]
	element2 := decodeNum(tempStack[len(tempStack)-1])
	*stack = tempStack[:len(tempStack)-1]
	if element1 == element2 {
		*stack = append(*stack, encodeNum(0))
	} else {
		*stack = append(*stack, encodeNum(1))
	}
	return true
}

func op_lessthan(stack *[][]byte) bool {
	if len(*stack) < 2 {
		return false
	}
	tempStack := *stack
	element1 := decodeNum(tempStack[len(tempStack)-1])
	tempStack = tempStack[:len(tempStack)-1]
	element2 := decodeNum(tempStack[len(tempStack)-1])
	*stack = tempStack[:len(tempStack)-1]
	if element2 < element1 {
		*stack = append(*stack, encodeNum(1))
	} else {
		*stack = append(*stack, encodeNum(0))
	}
	return true
}
func op_greaterthan(stack *[][]byte) bool {
	if len(*stack) < 2 {
		return false
	}
	tempStack := *stack
	element1 := decodeNum(tempStack[len(tempStack)-1])
	tempStack = tempStack[:len(tempStack)-1]
	element2 := decodeNum(tempStack[len(tempStack)-1])
	*stack = tempStack[:len(tempStack)-1]
	if element2 > element1 {
		*stack = append(*stack, encodeNum(1))
	} else {
		*stack = append(*stack, encodeNum(0))
	}
	return true
}

func op_lessthanorequal(stack *[][]byte) bool {
	if len(*stack) < 2 {
		return false
	}
	tempStack := *stack
	element1 := decodeNum(tempStack[len(tempStack)-1])
	tempStack = tempStack[:len(tempStack)-1]
	element2 := decodeNum(tempStack[len(tempStack)-1])
	*stack = tempStack[:len(tempStack)-1]
	if element2 <= element1 {
		*stack = append(*stack, encodeNum(1))
	} else {
		*stack = append(*stack, encodeNum(0))
	}
	return true
}

func op_greaterthanorequal(stack *[][]byte) bool {
	if len(*stack) < 2 {
		return false
	}
	tempStack := *stack
	element1 := decodeNum(tempStack[len(tempStack)-1])
	tempStack = tempStack[:len(tempStack)-1]
	element2 := decodeNum(tempStack[len(tempStack)-1])
	*stack = tempStack[:len(tempStack)-1]
	if element2 >= element1 {
		*stack = append(*stack, encodeNum(1))
	} else {
		*stack = append(*stack, encodeNum(0))
	}
	return true
}

func op_min(stack *[][]byte) bool {
	if len(*stack) < 2 {
		return false
	}
	tempStack := *stack
	element1 := decodeNum(tempStack[len(tempStack)-1])
	tempStack = tempStack[:len(tempStack)-1]
	element2 := decodeNum(tempStack[len(tempStack)-1])
	*stack = tempStack[:len(tempStack)-1]
	if element1 < element2 {
		*stack = append(*stack, encodeNum(int(element1)))
	} else {
		*stack = append(*stack, encodeNum(int(element2)))
	}
	return true
}

func op_max(stack *[][]byte) bool {
	if len(*stack) < 2 {
		return false
	}
	tempStack := *stack
	element1 := decodeNum(tempStack[len(tempStack)-1])
	tempStack = tempStack[:len(tempStack)-1]
	element2 := decodeNum(tempStack[len(tempStack)-1])
	*stack = tempStack[:len(tempStack)-1]
	if element1 > element2 {
		*stack = append(*stack, encodeNum(int(element1)))
	} else {
		*stack = append(*stack, encodeNum(int(element2)))
	}
	return true
}
func op_within(stack *[][]byte) bool {
	if len(*stack) < 3 {
		return false
	}
	tempStack := *stack

	maximum := decodeNum(tempStack[len(tempStack)-1])
	tempStack = tempStack[:len(tempStack)-1]
	minimum := decodeNum(tempStack[len(tempStack)-1])
	tempStack = tempStack[:len(tempStack)-1]
	element := decodeNum(tempStack[len(tempStack)-1])
	*stack = tempStack[:len(tempStack)-1]

	if element >= minimum && element < maximum {
		*stack = append(*stack, encodeNum(1))
	} else {
		*stack = append(*stack, encodeNum(0))
	}
	return true
}

func op_ripemd160(stack *[][]byte) bool {
	if len(*stack) < 1 {
		return false
	}
	tempStack := *stack
	element := tempStack[len(tempStack)-1]
	*stack = tempStack
	ripemdHasher := ripemd160.New()
	ripemdHasher.Write(element)
	hashBytes := ripemdHasher.Sum(nil)
	*stack = append(*stack, hashBytes)
	return true
}

func op_sha1(stack *[][]byte) bool {
	if len(*stack) < 1 {
		return false
	}
	tempStack := *stack
	element := tempStack[len(tempStack)-1]
	*stack = tempStack[:len(tempStack)-1]
	sha1Hasher := sha1.New()
	sha1Hasher.Write(element)
	hashBytes := sha1Hasher.Sum(nil)
	*stack = append(*stack, hashBytes)
	return true
}

func op_sha256(stack *[][]byte) bool {
	if len(*stack) < 1 {
		return false
	}
	tempStack := *stack
	element := tempStack[len(tempStack)-1]
	*stack = tempStack[:len(tempStack)-1]
	sha256Hasher := sha256.New()
	sha256Hasher.Write(element)
	hashBytes := sha256Hasher.Sum(nil)
	*stack = append(*stack, hashBytes)
	return true
}
func op_hash160(stack *[][]byte) bool {
	if len(*stack) < 1 {
		return false
	}
	tempStack := *stack
	element := tempStack[len(tempStack)-1]
	*stack = tempStack[:len(tempStack)-1]
	// Do a sha256 followed by a ripemd160
	hash256 := sha256.Sum256(element)
	ripemdHasher := ripemd160.New()
	ripemdHasher.Write(hash256[:])
	hashBytes := ripemdHasher.Sum(nil)
	*stack = append(*stack, hashBytes)
	return true
}

func op_hash256(stack *[][]byte) bool {
	if len(*stack) < 1 {
		return false
	}
	tempStack := *stack
	element := tempStack[len(tempStack)-1]
	*stack = tempStack[:len(tempStack)-1]
	// Do a 2 rounds of sha256
	hash256 := sha256.Sum256(element)
	hash256 = sha256.Sum256(hash256[:])
	*stack = append(*stack, hash256[:])
	return true
}

func op_checksig(stack *[][]byte, z string) bool {
	if len(*stack) < 2 {
		return false
	}
	//Get the last two variables
	tempStack := *stack
	sec := tempStack[len(tempStack)-1]
	tempStack = tempStack[:len(tempStack)-1]
	der := tempStack[len(tempStack)-1]
	*stack = tempStack[:len(tempStack)-1]
	point := secp256k1.ParseSec(hex.EncodeToString(sec))
	signature := secp256k1.ParseDer(hex.EncodeToString(der))
	fmt.Println(point.Verify("0x"+z, signature))
	if point.Verify("0x"+z, signature) {
		*stack = append(*stack, encodeNum(1))
	} else {
		*stack = append(*stack, encodeNum(0))
	}

	return true
}

func op_checksigverify(stack *[][]byte, z string) bool {
	return op_checksig(stack, z) && op_verify(stack)
}

func op_checkmultisig(stack *[][]byte, z string) bool {
	// THIS IS NOT IMPLEMENTED YET
	return true
}

func op_checkmultisigverify(stack *[][]byte, z string) bool {
	return op_checkmultisig(stack, z) && op_verify(stack)
}

func op_checklocktimeverify(stack *[][]byte, locktime []byte, sequence []byte) bool {
	lockTimeInt := binary.LittleEndian.Uint32(locktime)
	sequenceInt := binary.LittleEndian.Uint32(sequence)
	if sequenceInt == 0xffffffff {
		return false
	}
	if len(*stack) < 1 {
		return false
	}
	tempStack := *stack
	element := decodeNum(tempStack[len(tempStack)-1])
	if element < 0 {
		return false
	}
	if element < 500000000 && lockTimeInt > 500000000 {
		return false
	}
	if int64(lockTimeInt) < element {
		return false
	}
	return true
}

func op_checksequenceverify(stack *[][]byte, version []byte, sequence []byte) bool {
	versionInt := binary.LittleEndian.Uint32(version)
	sequenceInt := binary.LittleEndian.Uint32(sequence)
	if sequenceInt&(1<<31) == (1 << 31) {
		return false
	}
	if len(*stack) < 1 {
		return false
	}
	tempStack := *stack
	element := decodeNum(tempStack[len(tempStack)-1])
	if element < 0 {
		return false
	}
	if element&(1<<31) == (1 << 31) {
		if versionInt < 2 {
			return false
		} else if sequenceInt&(1<<31) == (1 << 31) {
			return false
		} else if uint32(element&(1<<22)) != sequenceInt&(1<<22) {
			return false
		} else if uint32(element&0xffff) > sequenceInt&0xffff {
			return false
		}
	}
	return true

}

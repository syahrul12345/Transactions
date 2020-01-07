package opcodes

import (
	"bytes"
	"encoding/binary"
)

//GetOPCODELIST returns a mapping representing the OPCODELIST
func GetOPCODELIST() map[int]interface{} {
	OPCODELIST := map[int]interface{}{
		0:  op_0,
		79: op_1negate,
		81: op_1,
		82: op_2,
		83: op_3,
		84: op_4,
		85: op_5,
		86: op_6,
		87: op_7,
		88: op_8,
		89: op_9,
		90: op_10,
		91: op_11,
		92: op_12,
		93: op_13,
		94: op_14,
		95: op_15,
		96: op_16,
		97: op_nop,
		99: op_if,
		// 100: op_notif,
		// 105: op_verify,
		// 106: op_return,
		// 107: op_toaltstack,
		// 108: op_fromaltstack,
		// 109: op_2drop,
		// 110: op_2dup,
		// 111: op_3dup,
		// 112: op_2over,
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

func op_0(commands [][]byte) bool {
	commands = append(commands, []byte{0})
	return true
}
func op_1negate(commands [][]byte) bool {
	buf := make([]byte, 1)
	binary.PutVarint(buf, -1)
	commands = append(commands, buf)
	return true
}

func op_1(commands [][]byte) bool {
	commands = append(commands, []byte{1})
	return true
}
func op_2(commands [][]byte) bool {
	commands = append(commands, []byte{2})
	return true
}
func op_3(commands [][]byte) bool {
	commands = append(commands, []byte{3})
	return true
}
func op_4(commands [][]byte) bool {
	commands = append(commands, []byte{4})
	return true
}
func op_5(commands [][]byte) bool {
	commands = append(commands, []byte{5})
	return true
}
func op_6(commands [][]byte) bool {
	commands = append(commands, []byte{6})
	return true
}
func op_7(commands [][]byte) bool {
	commands = append(commands, []byte{7})
	return true
}
func op_8(commands [][]byte) bool {
	commands = append(commands, []byte{8})
	return true
}
func op_9(commands [][]byte) bool {
	commands = append(commands, []byte{9})
	return true
}
func op_10(commands [][]byte) bool {
	commands = append(commands, []byte{10})
	return true
}
func op_11(commands [][]byte) bool {
	commands = append(commands, []byte{11})
	return true
}
func op_12(commands [][]byte) bool {
	commands = append(commands, []byte{12})
	return true
}
func op_13(commands [][]byte) bool {
	commands = append(commands, []byte{13})
	return true
}
func op_14(commands [][]byte) bool {
	commands = append(commands, []byte{14})
	return true
}
func op_15(commands [][]byte) bool {
	commands = append(commands, []byte{15})
	return true
}
func op_16(commands [][]byte) bool {
	commands = append(commands, []byte{16})
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
	if number == 0 {
		commands = trueItems
	} else {
		commands = falseItems
	}
	return true
}

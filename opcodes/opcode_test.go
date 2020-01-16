package opcodes

import (
	"encoding/hex"
	"testing"
)

func TestOpHash160(t *testing.T) {
	stack := &[][]byte{
		[]byte("hello world"),
	}
	command := 169
	operation := GetOPCODELIST()[command]
	operation.(func(*[][]byte) bool)(stack)

	tempStack := *stack
	get := hex.EncodeToString(tempStack[0])
	want := "d7d5ee7824ff93f94c3055af9382c86c68b5ca92"
	if get != want {
		t.Errorf("Expected OP_HASH160 to give %s but got %s", want, get)
	}

}

func TestOPSig(t *testing.T) {
	z := "0x7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d"
	sec := "04887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34"
	sig := "3045022000eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c022100c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab601"
	secByte, _ := hex.DecodeString(sec)
	sigByte, _ := hex.DecodeString(sig)
	stack := &[][]byte{
		sigByte,
		secByte,
	}
	command := 172
	operation := GetOPCODELIST()[command]
	operation.(func(*[][]byte, string) bool)(stack, z)
	tempStack := *stack
	get := int8(tempStack[0][0])
	want := int8(1)
	if get != want {
		t.Errorf("Expected OP_SIGVERIFY to give %d but got %d", want, get)
	}

}

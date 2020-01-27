package models

import (
	"encoding/hex"
	"math/big"
	"strconv"
)

func reverse(heightByte *[]byte) {
	a := *heightByte
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
	*heightByte = a
}

func bitsToTarget(bits [4]byte) string {
	// Make a copy
	buf := make([]byte, 4)
	copy(buf, bits[:])

	exponentByte := buf[len(buf)-1] - 3
	coefficientBits := buf[:len(buf)-1]
	// Reverse it as it's little endian
	reverse(&coefficientBits)

	coefficientString := hex.EncodeToString(coefficientBits)
	coefficient, _ := strconv.ParseUint(coefficientString, 16, 24)
	// First term
	firstTem := big.NewInt(0).SetUint64(coefficient)
	// Second term
	// 256 ** (exponent-3)
	secondTerm := big.NewInt(0).Exp(big.NewInt(256), big.NewInt(0).SetBytes([]byte{exponentByte}), big.NewInt(0))
	// Answer
	target := big.NewInt(0).Mul(firstTem, secondTerm)
	return target.Text(16)
}

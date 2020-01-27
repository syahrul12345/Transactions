package utils

import (
	"encoding/hex"
	"testing"
)

func TestCalculateNewBits(t *testing.T) {
	prevBits := [4]byte{
		0x54,
		0xd8,
		0x01,
		0x18,
	}
	timeDifferential := uint64(302400)
	want := [4]byte{
		0x00,
		0x15,
		0x76,
		0x17,
	}
	get := CalculateNewBits(prevBits, timeDifferential)
	if get != want {
		t.Errorf("Expected the new bits to be %s but got %s", hex.EncodeToString(want[:]), hex.EncodeToString(get[:]))
	}
}

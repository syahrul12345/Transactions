package models

import (
	"encoding/hex"
	"testing"
)

func TestScriptParse(t *testing.T) {
	input := "6a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937"
	script, _ := ParseScript(input)
	command := hex.EncodeToString(script.Commands[0])
	commandWant := "304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a71601"
	if command != commandWant {
		t.Errorf("Expected the command to be %s but got %s", commandWant, command)
	}
	command = hex.EncodeToString(script.Commands[1])
	commandWant = "035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937"
	if command != commandWant {
		t.Errorf("Expected the command to be %s but got %s", commandWant, command)
	}
}

func TestScriptSerialize(t *testing.T) {
	input := "6a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937"
	script, _ := ParseScript(input)
	scriptSerialize := script.Serialize()
	if scriptSerialize != input {
		t.Errorf("Script serialization failed. Want %s but got %s", input, scriptSerialize)
	}
}

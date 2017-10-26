package evo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstructionNoop(t *testing.T) {
	var ti Instruction

	ti.operation = 0
	assert.True(t, ti.noop(), "Dogs")

	ti.operation = 9
	assert.False(t, ti.noop(), "Cats")
}


func TestRandomInstruction(t *testing.T) {
	// Compare random instructions and assert that they don't always match.

	var emptyIns1 Instruction
	var emptyIns2 Instruction
	randIns := NewRandomInstruction()

	// Sanity check for instruction compairison.
	assert.Equal(t, emptyIns1, emptyIns2, "Should be equal")

	trials := 10
	matches := 0
	for i:= 0; i < trials; i++ {
		randIns2 := NewRandomInstruction()
		if (randIns == randIns2) {
			matches++
		}
	}

	assert.NotEqual(t, trials, matches, "Some random instructions should not match")
}
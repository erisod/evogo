package evo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Check basic form copy operation and ensure modification to parent post-copy
// does not modify child.
func TestForm(t *testing.T) {
	FrankForm := NewRandomForm()
	FrankForm.instructions[0].operation = NOOP

	CopyForm := NewChildForm(FrankForm, false)

	for i:=0; i < CODESIZE; i++ {
		assert.Equal(t, FrankForm.instructions[i], CopyForm.instructions[i], "should match")
	}

	// Modify parent and confirm child is not modified.
	FrankForm.instructions[0].operation = ADDLEQ
	assert.NotEqual(t, FrankForm.instructions[0].operation, CopyForm.instructions[0].operation, "should not match")
}


func TestFormProgramIOCopy(t *testing.T) {
	f := NewRandomForm()

	f.instructions[0].operation = COPYIN
	f.instructions[0].p1 = 0
	f.instructions[0].p2 = 0

	f.instructions[1].operation = COPYRES
	f.instructions[1].p1 = 0
	f.instructions[1].p2 = 0

	f.instructions[2].operation = ENDEXEC

	input := []int{66, 77, 88}
	f.runCode(input)

	assert.Equal(t, 66, f.output[0])
}

func TestFormProgramInvalidRange(t *testing.T) {
	// Check that the program completes despite an out-of-range issue.
	f := NewRandomForm()

	f.instructions[0].operation = COPYIN
	f.instructions[0].p1 = 0
	f.instructions[0].p2 = 0

	f.instructions[1].operation = COPYRES
	f.instructions[1].p1 = 0
	f.instructions[1].p2 = 0

	f.instructions[2].operation = COPYIN
	f.instructions[2].p1 = 1
	f.instructions[2].p2 = 500 // Invalid mem location.

	f.instructions[3].operation = COPYRES
	f.instructions[3].p1 = 500 // Invalid mem location.
	f.instructions[3].p2 = 0

	f.instructions[4].operation = ENDEXEC

	input := []int{66, 77, 88}
	f.runCode(input)

	assert.Equal(t, 66, f.output[0])
}

func TestFormPrint(t *testing.T) {
	f := NewRandomForm()
	f.Print()
}
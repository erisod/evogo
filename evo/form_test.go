package evo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sort"
	"github.com/stretchr/testify/require"
)

// Check basic form copy operation and ensure modification to parent post-copy
// does not modify child.
func TestForm(t *testing.T) {
	tf := NewRandomForm()
	tf.instructions[0].operation = NOOP

	CopyForm := NewChildForm(tf, false)

	for i:=0; i < CODESIZE; i++ {
		assert.Equal(t, tf.instructions[i], CopyForm.instructions[i], "should match")
	}

	// Modify parent and confirm child is not modified.
	tf.instructions[0].operation = ADDLEQ
	assert.NotEqual(t, tf.instructions[0].operation, CopyForm.instructions[0].operation, "should not match")
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
	f.runCode(&input)

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
	f.runCode(&input)

	assert.Equal(t, 66, f.output[0])
}

func TestFormPrint(t *testing.T) {
	f := NewRandomForm()
	f.Print()
}

func TestFormSort(t *testing.T) {
	f1 := NewNoopForm()
	f2 := NewNoopForm()
	f3 := NewNoopForm()

	// Using opsleft to identify forms easily.
	f1.opsleft=1
	f2.opsleft=2
	f3.opsleft=3

	f1.costSum=11.0
	f2.costSum=5.0
	f3.costSum=2.0

	f1.runCount = 1
	f2.runCount = 1
	f3.runCount = 1

	f1.scoreSum=-5
	f2.scoreSum=-3
	f3.scoreSum=-4

	forms := []Form{f1, f2, f3}
	// Sanity check.
	require.Equal(t, 1, forms[0].opsleft)


	// Test function.
	sort.Sort(ByAvgScore(forms))

	// Lowest score should be first.
	require.Equal(t, 2, forms[0].opsleft)

	f1.scoreSum=0.0
	f2.scoreSum=0.0
	f3.scoreSum=0.0
	forms = []Form{f1, f2, f3}

	// Test function.
	sort.Sort(ByAvgScore(forms))

	// Lowest cost should be first.
	require.Equal(t, 3, forms[0].opsleft)
}
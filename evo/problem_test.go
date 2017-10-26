package evo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProblem(t *testing.T) {
	var tp AdditionProblem

	input := tp.GenerateInputs()

	assert.Equal(t, len(input), 10)
	assert.Equal(t, tp.Answer(input), input[0] + input[1])
}
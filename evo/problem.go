package evo

import (
	"time"
	"math/rand"
	"math"
	"fmt"
)

type ProblemInterface interface {
	Answer([]int) []int
	GenerateInputs() []int
	Score([]int) float64
}

type Problem struct {}

type AdditionProblem struct { Problem }
type SubtractionProblem struct { Problem }
type CopyProblem struct { Problem }
type Copy3Problem struct { Problem }

const PROBLEM_INPUT_RANGE = 200 // +- this value.

// Addition problem.
func (p *AdditionProblem) Answer(input []int) []int {
	answer := make([]int, 1)
	answer[0] = input[0] + input[1]
	return answer
}

// Subtraction problem.
func (p *SubtractionProblem) Answer(input []int) []int {
	answer := make([]int, 1)
	answer[0] = input[0] - input[1]
	return answer
}

// Copy problem.
func (p *CopyProblem) Answer(input []int) []int {
	answer := make([]int, 1)
	answer[0] = input[0]
	return answer
}

// Copy 3 problem.
func (p *Copy3Problem) Answer(input []int) []int {
	answer := make([]int, 3)
	answer[0] = input[0]
	answer[1] = input[1]
	answer[2] = input[2]
	return answer
}

// "Virtual" type function, not expected to be used.  TODO: Can I remove this?
func (p *Problem) Answer(input []int) []int {
	fmt.Println("I SHOULD NEVER BE CALLED DIRECTLY!")
	answer := make([]int, 0)
	return answer
}

// As convention the "highest" score is 0.0
func (p *Problem) Score(correct []int, actual []int) float64 {
	gap := 0.0
	for i:=0 ; i < len(correct) && i < len(actual); i++ {
		gap += math.Abs(float64(correct[i]) - float64(actual[i]))
	}

	return -gap
}

// For most problems we can generate all random inputs.
func (p *Problem) GenerateInputs() []int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	input := make([]int, 10)

	input[0] = rng.Intn(PROBLEM_INPUT_RANGE * 2) - PROBLEM_INPUT_RANGE
	input[1] = rng.Intn(PROBLEM_INPUT_RANGE * 2) - PROBLEM_INPUT_RANGE
	input[2] = rng.Intn(PROBLEM_INPUT_RANGE * 2) - PROBLEM_INPUT_RANGE
	input[3] = rng.Intn(PROBLEM_INPUT_RANGE * 2) - PROBLEM_INPUT_RANGE
	input[4] = rng.Intn(PROBLEM_INPUT_RANGE * 2) - PROBLEM_INPUT_RANGE
	input[5] = rng.Intn(PROBLEM_INPUT_RANGE * 2) - PROBLEM_INPUT_RANGE
	input[6] = rng.Intn(PROBLEM_INPUT_RANGE * 2) - PROBLEM_INPUT_RANGE
	input[7] = rng.Intn(PROBLEM_INPUT_RANGE * 2) - PROBLEM_INPUT_RANGE
	input[8] = rng.Intn(PROBLEM_INPUT_RANGE * 2) - PROBLEM_INPUT_RANGE
	input[9] = rng.Intn(PROBLEM_INPUT_RANGE * 2) - PROBLEM_INPUT_RANGE

	return input
}
package evo

import (
	"math/rand"
	"time"
)

const NOOP = 0;
const JUMP = 1;
const ADDLEQ = 2;
const DECNZJ = 3;
const INCEQ = 4;
const SUBLEQ = 5;
const COPYRES = 6;
const SETVAL = 7;
const ENDEXEC = 8;
const COPYIN = 9;
const MAX_OPERATION = 9;

const MIN_PARAM = -200
const MAX_PARAM = 200

type Instruction struct {
	// The operation code.
	operation int

	// The parameters to the operation.
	p1 int
	p2 int
	p3 int
	p4 int
}

func (i *Instruction) noop() bool {
	if i.operation == NOOP {
		return true
	}
	return false
}

// Is this instruction valid?
func (i *Instruction) valid() bool {
	if i.operation < 0 || i.operation > MAX_OPERATION {
		return false
	}

	return true
}

func (i *Instruction) Copy() *Instruction {
	newins := new(Instruction)

	newins.operation = i.operation
	newins.p1 = i.p1
	newins.p2 = i.p2
	newins.p3 = i.p3
	newins.p4 = i.p4

	return newins
}

func NewRandomInstruction() *Instruction {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	ins := new(Instruction)
	ins.operation = rng.Intn(MAX_OPERATION)
	ins.p1 = rng.Intn(MAX_PARAM - MIN_PARAM) + MIN_PARAM
	ins.p2 = rng.Intn(MAX_PARAM - MIN_PARAM) + MIN_PARAM
	ins.p3 = rng.Intn(MAX_PARAM - MIN_PARAM) + MIN_PARAM
	ins.p4 = rng.Intn(MAX_PARAM - MIN_PARAM) + MIN_PARAM

	return ins
}

func NewMutantInstruction(*Instruction) *Instruction {
	ins := new(Instruction)

	mutationOdds := 1000 // e.g. 1 in mutationOdds.

	maybeMutate(&ins.operation, mutationOdds)
	maybeMutate(&ins.p1, mutationOdds)
	maybeMutate(&ins.p2, mutationOdds)
	maybeMutate(&ins.p3, mutationOdds)
	maybeMutate(&ins.p4, mutationOdds)

	return ins
}

func maybeMutate(value *int, mutationOdds int) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Increment value mutation.
	if rng.Intn(mutationOdds) == 0 {
		*value += rng.Intn(MAX_OPERATION * 2) - MAX_OPERATION;
	}

	// Sign flip mutation.
	if rng.Intn(mutationOdds) == 0 {
		*value = - *value
	}
}
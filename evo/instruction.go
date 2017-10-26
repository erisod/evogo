package evo

import (
	"math/rand"
	"time"
	"strconv"
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

func (i *Instruction) getDesc() string {

	var desc string
	var longdesc string

	switch i.operation {
	case NOOP:
		desc = "noop"
		longdesc = "do nothing"
	case JUMP:
		desc = "jump " + strconv.Itoa(i.p1)
		longdesc = "jump to code" + strconv.Itoa(i.p1)
	case ADDLEQ:
		desc = "addleq " + strconv.Itoa(i.p1) + " " + strconv.Itoa(i.p2) + " " + strconv.Itoa(i.p3)
		longdesc = "mem" + strconv.Itoa(i.p1) + "+= mem" + strconv.Itoa(i.p2) + "; if mem" + strconv.Itoa(i.p1) + " <= 0 jump to code" + strconv.Itoa(i.p3)
	case DECNZJ:
		desc = "decnzj " + strconv.Itoa(i.p1) + " " + strconv.Itoa(i.p2) + " " + strconv.Itoa(i.p3)
		longdesc = "mem" + strconv.Itoa(i.p1) + "-= mem" + strconv.Itoa(i.p2) + "; if mem" + strconv.Itoa(i.p1) + " !=0 jump to code" + strconv.Itoa(i.p3)
	case INCEQ:
		desc = "inceq " + strconv.Itoa(i.p1) + " " + strconv.Itoa(i.p2) + " " + strconv.Itoa(i.p3)
		longdesc = "mem" + strconv.Itoa(i.p1) + "++; if mem" + strconv.Itoa(i.p1) + "==mem" + strconv.Itoa(i.p2) + " jump to code" + strconv.Itoa(i.p3)
	case COPYRES:
		desc = "copyres " + strconv.Itoa(i.p1) + " " + strconv.Itoa(i.p2)
		longdesc = "output" + strconv.Itoa(i.p2) + "=mem" + strconv.Itoa(i.p1)
	case SETVAL:
		desc = "setval " + strconv.Itoa(i.p1) + " " + strconv.Itoa(i.p2)
		longdesc = "mem" + strconv.Itoa(i.p1) + "=" + strconv.Itoa(i.p2)
	case ENDEXEC:
		desc = "endexec"
		longdesc = "stop program"
	case COPYIN:
		desc = "copyin " + strconv.Itoa(i.p1) + " " + strconv.Itoa(i.p2)
		longdesc = "mem" + strconv.Itoa(i.p2) + "=input" + strconv.Itoa(i.p1)
	default:
		desc = "invalid(op" + strconv.Itoa(i.operation) + ")"
		longdesc = "invalid operation code " + strconv.Itoa(i.operation)
	}

	desc = desc + "\t(" + longdesc + ")"
	return desc
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

func (i *Instruction) Copy() Instruction {
	newins := Instruction{}

	newins.operation = i.operation
	newins.p1 = i.p1
	newins.p2 = i.p2
	newins.p3 = i.p3
	newins.p4 = i.p4

	return newins
}

func NewRandomInstruction() Instruction {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	ins := Instruction{}
	ins.operation = rng.Intn(MAX_OPERATION)
	ins.p1 = rng.Intn(MAX_PARAM - MIN_PARAM) + MIN_PARAM
	ins.p2 = rng.Intn(MAX_PARAM - MIN_PARAM) + MIN_PARAM
	ins.p3 = rng.Intn(MAX_PARAM - MIN_PARAM) + MIN_PARAM
	ins.p4 = rng.Intn(MAX_PARAM - MIN_PARAM) + MIN_PARAM

	return ins
}

func NewMutantInstruction(parent Instruction) Instruction {
	ins := parent.Copy()

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
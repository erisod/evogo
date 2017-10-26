package evo

import (
	"time"
	"math/rand"
	"fmt"
)

type Form struct {
	instructions []*Instruction
	mem          []int
	input        []int
	output       []int

	codesize int
	memsize  int
	iosize   int

	finished bool
	opsleft  int

	avgScore float64
	scoreSum float64
	runCount int

	// Code pointer.
	cp int
}

const CODESIZE = 200
const MEMSIZE = 10
const IOSIZE = 10
const MAXOPS = 100

func (f *Form) Print() {
	fmt.Println("Printing Form ---- ")
	fmt.Printf("%+v\n", *f)
}

func NewRandomForm() *Form {
	f := new(Form)
	f.init()

	for i := 0; i < CODESIZE; i++ {
		f.instructions[i] = NewRandomInstruction()
	}

	return f
}

// Create a new form based on a parent.  Mutation optional.
func NewChildForm(parent *Form, mutate bool) *Form {
	f := new(Form)
	f.init()

	pPos := 0
	cPos := 0

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for pPos < CODESIZE && cPos < CODESIZE {

		// Normal instruction copy with mutation.
		if (mutate) {
			f.instructions[cPos] = NewMutantInstruction(parent.instructions[pPos])
		} else {
			f.instructions[cPos] = parent.instructions[pPos].Copy()
		}

		// Skip or duplicate some of parent.
		if (mutate && rng.Intn(CODESIZE) == 0) {
			pPos = rng.Intn(CODESIZE)
		}

		pPos++
		cPos++
	}

	// Fill reminder of child with random instructions
	for cPos < CODESIZE {
		f.instructions[cPos] = NewRandomInstruction()
	}

	return f
}

func (f *Form) AvgScore() float64{
	return f.scoreSum / float64(f.runCount)
}

func (f *Form) init() {
	f.output = make([]int, IOSIZE)
	f.mem = make([]int, MEMSIZE)
	f.instructions = make([]*Instruction, CODESIZE, CODESIZE)

	f.reset()
}

func (f *Form) reset() {
	f.cp = 0
	f.finished = false
	f.opsleft = MAXOPS

	for i := 0; i < IOSIZE; i++ {
		f.output[i] = 0;
	}
}

func (f *Form) runCode(newInput []int) {
	f.input = newInput

	f.reset()

	for (!f.finished && f.opsleft > 0) {
		f.step()
		f.opsleft--
	}
}

func (f *Form) step() {
	if (f.cp > CODESIZE || f.cp < 0) {
		f.finished = true
		return
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("RECOVERING ", r)
			f.finished = true
		}
	}()

	op := f.instructions[f.cp].operation

	switch op {
	case NOOP:
		f.cp++
	case JUMP:
		f.jump()
	case ADDLEQ:
		f.addleq()
	case DECNZJ:
		f.decnzj()
	case INCEQ:
		f.inceq()
	case COPYRES:
		f.copyToResult()
	case SETVAL:
		f.setval()
	case ENDEXEC:
		f.endexec()
	case COPYIN:
		f.copyFromInput()
	default:
		f.endexec()
	}
}

// Move cp to p1.
func (f *Form) jump() {
	f.cp = f.instructions[f.cp].p1
}

func (f *Form) addleq() {
	ins := f.instructions[f.cp]
	f.mem[ins.p1] += f.mem[ins.p2]

	if f.mem[ins.p1] <= 0 {
		f.cp = f.mem[ins.p3]
	} else {
		f.cp++
	}
}

func (f *Form) decnzj() {
	ins := f.instructions[f.cp]

	f.mem[ins.p1] = f.mem[ins.p1] - f.mem[ins.p2]
	if f.mem[ins.p1] < 0 {
		f.cp = ins.p3
	} else {
		f.cp++
	}
}

func (f *Form) inceq() {
	ins := f.instructions[f.cp]

	f.mem[ins.p1] += 1
	if (f.mem[ins.p1] == f.mem[ins.p2]) {
		f.cp = ins.p3
	} else {
		f.cp++
	}

}

func (f *Form) subleq() {
	ins := f.instructions[f.cp]

	f.mem[ins.p1] -= f.mem[ins.p2]
	if (f.mem[ins.p1] <= f.mem[ins.p3]) {
		f.cp = ins.p4
	} else {
		f.cp++
	}
}

func (f *Form) copyToResult() {
	ins := f.instructions[f.cp]

	f.output[ins.p2] = f.mem[ins.p1]
	f.cp++
}

func (f *Form) copyFromInput() {
	ins := f.instructions[f.cp]

	f.mem[ins.p2] = f.input[ins.p1]
	f.cp++
}

func (f *Form) setval() {
	ins := f.instructions[f.cp]

	f.mem[ins.p1] = ins.p2
	f.cp++
}

func (f *Form) endexec() {
	f.finished = true
}



type ByAvgScore []Form
func (f ByAvgScore) Len() int {
	return len(f)
}
func (f ByAvgScore) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
func (f ByAvgScore) Less(i, j int) bool {

	return f[i].AvgScore() < f[j].AvgScore()
}


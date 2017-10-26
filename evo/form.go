package evo

import (
	"fmt"
	"strconv"
)

type Form struct {
	instructions []Instruction
	mem          []int
	input        []int
	output       []int

	finished bool
	opsleft  int

	scoreSum float64
	runCount int
	costSum float64

	// Code pointer.
	cp int
}

const CODESIZE = 10
const MEMSIZE = 10
const IOSIZE = 10
const MAXOPS = 10
const MUTATIONRATE = 50


func (f *Form) Description() string {
	var desc string

	desc += "Input:\n"
	for i:=0; i<len(f.input); i++ {
			desc += "  input[" + strconv.Itoa(i) + "] = " + strconv.Itoa(f.input[i]) + "\n"
	}
	desc += "Code:\n"
	for i:=0; i<len(f.instructions); i++ {
			desc += "  " + strconv.Itoa(i) + " : " + f.instructions[i].getDesc() + "\n"
	}
	desc += "Output (zeros suppressed):\n"
	for i:=0; i<len(f.output); i++ {
		if f.output[i] != 0 {
			desc += "  output[" + strconv.Itoa(i) + "] = " + strconv.Itoa(f.output[i]) + "\n"
		}
	}
	desc += "Stats:\n"
	desc += "  AvgScore: " + fmt.Sprintf("%f", f.AvgScore()) + "\n"
	desc += "  RunCount: " + strconv.Itoa(f.runCount) + "\n"
	desc += "  AvgCost: " + fmt.Sprintf("%f", f.AvgCost()) + "\n"

	return desc
	
}

func (f *Form) Print() {
	fmt.Println(f.Description())
	// fmt.Printf("%+v\n", f) // Print raw struct.
}

func NewNoopForm() Form {
	f := Form{}
	f.init()

	for i:=0; i < CODESIZE; i++ {
		f.instructions = append(f.instructions, Instruction{})
	}

	return f
}


// A form which copies input0 to output0 (for testing).
func NewCopyForm() Form {
	f := Form{}
	f.init()

	f.instructions = append(f.instructions, Instruction{ operation:COPYIN})
	f.instructions = append(f.instructions, Instruction{ operation:COPYRES})
	f.instructions = append(f.instructions, Instruction{ operation:ENDEXEC})

	return f
}


func NewRandomForm() Form {
	f := Form{}
	f.init()

	for i:=0; i < CODESIZE; i++ {
		f.instructions = append(f.instructions, NewRandomInstruction())
	}

	return f
}

// Create a new form based on a parent.  Mutation optional.
func NewChildForm(parent Form, mutate bool) Form {
	f := Form{}
	f.init()

	f.instructions = make([]Instruction, CODESIZE)

	pPos := 0
	cPos := 0

	for pPos < len(parent.instructions) && cPos < len(f.instructions) {

		// Normal instruction copy with mutation.
		if (mutate) {
			f.instructions[cPos] = NewMutantInstruction(parent.instructions[pPos])
		} else {
			f.instructions[cPos] = parent.instructions[pPos].Copy()
		}

		// Skip or duplicate some of parent.
		if (mutate && rng.Intn(MUTATIONRATE) == 0) {
			pPos = rng.Intn(CODESIZE)
		}
		// Overwrite or skip part of child.
		if (mutate && rng.Intn(MUTATIONRATE) == 0) {
			cPos = rng.Intn(CODESIZE)
		}

		pPos++
		cPos++
	}

	// Fill reminder of child with random instructions
	for ; cPos < CODESIZE ; cPos++ {
		f.instructions = append(f.instructions, NewRandomInstruction())
	}

	return f
}

func (f *Form) AvgScore() float64{
	return f.scoreSum / float64(f.runCount)
}

func (f *Form) AvgCost() float64{
	return f.costSum / float64(f.runCount)
}

func (f *Form) init() {
	f.output = make([]int, IOSIZE)
	f.mem = make([]int, MEMSIZE)

	f.reset()
}

func (f *Form) reset() {
	f.cp = 0
	f.finished = false
	f.opsleft = MAXOPS

	for i := 0; i < len(f.output) ; i++ {
		f.output[i] = 0;
	}

	for i := 0; i < len(f.mem) ; i++ {
		f.mem[i] = 0;
	}
}

func (f *Form) resetStats() {
	f.costSum = 0.0
	f.runCount = 0.0
}

func (f *Form) runCode(newInput *[]int) {
	f.input = *newInput

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
			// fmt.Println("RECOVERING ", r)
			f.finished = true
		}
	}()

	op := f.instructions[f.cp].operation

	f.costSum += 1.0

	switch op {
	case NOOP:
		f.cp++
		f.costSum -= 0.9 // Count a noop as a discount operation.
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
		// Invalid operations end the program.
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


// Sorter for Form; highest scores first; if scores both zero
// sort by lower cost.
type ByAvgScore []Form
func (f ByAvgScore) Len() int {
	return len(f)
}
func (f ByAvgScore) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
func (f ByAvgScore) Less(i, j int) bool {
	if f[i].AvgScore() == 0.0 && f[j].AvgScore() == 0.0 {
		return f[i].AvgCost() < f[j].AvgCost()
	}
	return f[i].AvgScore() > f[j].AvgScore()
}


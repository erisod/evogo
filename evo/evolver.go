package evo

import (
	"math"
	"fmt"
	"sort"
)

// An Evolver has a set of (!life) forms and a problem/scorer that it uses to
// evolve the forms.
type Evolver struct {
	forms []Form

	problem ProblemInterface

	// Is the problem solved (may be inefficient).
	solved bool

	// Has the solved problem's score not improved in STABILITYDURATION iterations.
	solvedNStable bool

	// Count of times the same top score has been produced.
	sameSolvedTopScoreCount int

	// Most recent topScore
	lastTopScore float64

	// Highest ever topScore
	topScore float64
}

const MAXFORMS = 100
const STABILITYDURATION = 1000
const RACETRIALS = 10

func NewEvolver(p ProblemInterface) Evolver {
	e := Evolver{}
	e.solved = false
	e.solvedNStable = false
	e.topScore = -math.MaxFloat64
	e.forms = []Form{}

	for _ = range [MAXFORMS] struct{}{} {
		e.forms = append(e.forms, NewRandomForm())
	}

	fmt.Println("I made a new Evolver! with", len(e.forms),"forms")
	fmt.Println("form 0 has ", len(e.forms[0].instructions), "instructions")

	e.problem = p

	return e
}

func (e *Evolver) mutateForms() {
	// Take the top N% of forms and duplicate them into new slots.

	var topPct float32 = 10

	topN := int(float32(len(e.forms)) * float32(topPct/100))
	newForms := []Form{}
	newPerTop := int(float32(MAXFORMS)/float32(topN))

	for i:=0; i< topN; i++ {
		for j:=0; j < newPerTop; j++ {
			nf := NewChildForm(e.forms[i], true)
			newForms = append(newForms, nf)
		}
	}

	e.forms = newForms
}

func (e *Evolver) runIteration() {
	// TODO: Is there a cleaner way of doing this loop?
	for _ = range [RACETRIALS] struct{}{} {
		problemInput := e.problem.GenerateInputs()
		problemAnswer := e.problem.Answer(problemInput)
		for i := 0; i < len(e.forms); i++ {
			e.forms[i].runCode(problemInput)
			e.forms[i].runCount++
			runScore := e.problem.Score(problemAnswer, e.forms[i].output)
			e.forms[i].scoreSum += runScore
		}
	}

}

func (e *Evolver) sortFormsByAvgScore() {
	// TODO: Understand this better; it feels backwards.
	sort.Sort(sort.Reverse(ByAvgScore(e.forms)))
}

func (e *Evolver) doBookKeeping() {
	// Relies on the forms being reverse sorted by avgscore
	runTopScore := e.forms[0].AvgScore()
	e.lastTopScore = runTopScore
	if (runTopScore == 0.0) {
		e.solved = true
		if e.topScore == runTopScore {
			e.sameSolvedTopScoreCount++
			if e.sameSolvedTopScoreCount >= STABILITYDURATION {
				e.solvedNStable = true
			}
		}
	}
}

// Run the evolution until complete (or FOREVER) and report status via stdout.
func (e *Evolver) RunAndReport() {
	for i:=0 ; ; i++ {
		fmt.Println("Starting iteration", i)
		e.runIteration()
		e.sortFormsByAvgScore()
		e.doBookKeeping()

		fmt.Println("Iteration ",i," complete.  runTopScore : ", e.lastTopScore)
		fmt.Println("Best form:")
		e.forms[0].Print()

		if e.solvedNStable {
			fmt.Println("Stable solution!")
			break
		}

		e.mutateForms()
	}
}
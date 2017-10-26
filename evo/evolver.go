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
	sameSolvedCostCount int

	// Most recent topScore
	lastTopScore float64

	// Highest ever topScore
	topScore float64

	// Best ever (lowest) cost
	lastTopCost float64
}

const MAXFORMS = 20000
const STABILITYDURATION = 500
const RACETRIALS = 20

func NewEvolver(p ProblemInterface) Evolver {
	e := Evolver{}
	e.solved = false
	e.solvedNStable = false
	e.topScore = -math.MaxFloat64
	e.forms = []Form{}

	for _ = range [MAXFORMS] struct{}{} {
		// Create random value or noop (all zero) initial forms.
		// e.forms = append(e.forms, NewRandomForm())
		e.forms = append(e.forms, NewNoopForm())
		// e.forms = append(e.forms, NewCopyForm())
	}

	e.problem = p

	return e
}

// Mutate forms by allocating an all new set of forms based on
// the top N% best performing forms.
func (e *Evolver) mutateForms() {
	// Take the top N% of forms and duplicate them into new slots.

	var topPct float32 = 20

	// TODO: Clean up this type wrangling; I was having some debugger issues here and wanted
	// to inspect intermediate values with prints.
	topNFloat := float32(len(e.forms)) * float32(topPct)/100.0
	topN := int(topNFloat)
	newForms := []Form{}
	newPerTop := int(float32(MAXFORMS)/float32(topN))

	for i:=0; i< topN; i++ {
		// Copy one intact.
		nf := NewChildForm(e.forms[i], false)
		newForms = append(newForms, nf)

		for j:=1; j < newPerTop; j++ {
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
			e.forms[i].runCode(&problemInput)
			e.forms[i].runCount++
			runScore := e.problem.Score(problemAnswer, e.forms[i].output)
			e.forms[i].scoreSum += runScore
		}
	}

}

func (e *Evolver) sortFormsByAvgScore() {
	// TODO: Understand this better; it feels backwards.
	sort.Sort(ByAvgScore(e.forms))
}

func (e *Evolver) doBookKeeping() {
	// Relies on the forms being reverse sorted by avgscore
	runTopScore := e.forms[0].AvgScore()
	e.lastTopScore = runTopScore
	if (runTopScore == 0.0) {
		e.solved = true
		runTopCost := e.forms[0].AvgCost()
		if e.lastTopCost == runTopCost {
			e.sameSolvedCostCount++
			if e.sameSolvedCostCount >= STABILITYDURATION {
				e.solvedNStable = true
			}
		} else {
			e.sameSolvedCostCount = 0
		}

		e.lastTopCost = runTopCost
	}

	e.topScore = math.Min(e.topScore, runTopScore)
}

// Run the evolution until complete (or FOREVER) and report status via stdout.
func (e *Evolver) RunAndReport() {
	for i:=0 ; ; i++ {
		e.runIteration()
		e.sortFormsByAvgScore()
		e.doBookKeeping()

		if (i % 10 == 0) {
			fmt.Println("Best form:")
			e.forms[0].Print()

			if e.solved {
				fmt.Println("--Solved--  Stable for", e.sameSolvedCostCount, "iterations")
			}
		}

		fmt.Println("Iteration ", i, " complete.  runTopScore : ", e.forms[0].AvgScore(), "cost:", e.forms[0].AvgCost())

		if e.solvedNStable {
			fmt.Println("Stable solution!")
			break
		}

		e.mutateForms()
	}
}
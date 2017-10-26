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

	problem Problem

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

const MAXFORMS = 10000
const STABILITYDURATION = 1000
const RACETRIALS = 50

func NewEvolver(p Problem) Evolver {
	e := Evolver{}
	e.solved = false
	e.solvedNStable = false
	e.topScore = -math.MaxFloat64
	e.forms = make([]Form, MAXFORMS)

	for i:=0 ; i < MAXFORMS ; i++ {
		e.forms[i] = *NewRandomForm()
	}

	e.problem = p

	return e
}

func (e Evolver) mutateForms() {
	// Take the top 10% of forms and duplicate them into new slots.
	topN := int(float32(len(e.forms)) * .1)
	newForms := make([]Form, len(e.forms))

	for i:=0; i< topN; i++ {
		nf := NewChildForm(&(e.forms[i]), true)
		newForms = append(newForms, *nf)
	}

	e.forms = newForms
}

func (e Evolver) runIteration() {
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

func (e Evolver) sortFormsByAvgScore() {
	// TODO: Understand this better; it feels backwards.
	sort.Sort(sort.Reverse(ByAvgScore(e.forms)))
}

func (e Evolver) doBookKeeping() {
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
func (e Evolver) runAndReport() {
	for i:=0 ; ; i++ {
		e.runIteration()
		e.sortFormsByAvgScore()
		e.doBookKeeping()

		fmt.Println("Iteration ",i," complete.  runTopScore : ", e.lastTopScore)
		if e.solvedNStable {
			break
		}

		e.mutateForms()
	}
}
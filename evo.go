package main

import (
	"github.com/erisod/evogo/evo"
	"fmt"
)

func main() {
	var problem evo.AdditionProblem

	e := evo.NewEvolver(problem)

	e.RunAndReport()

	fmt.Println("all done")
}

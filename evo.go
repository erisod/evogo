package main

import (
	"github.com/erisod/evogo/evo"
	"fmt"
)

// Global random number generator.

func main() {
	var problem evo.Output1Problem

	e := evo.NewEvolver(problem)

	e.RunAndReport()

	fmt.Println("all done")
}

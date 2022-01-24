package service

import (
	"github.com/jayanpraveen/tildly/datastore"
)

type atomicCounter struct {
	min int
	max int
	cod datastore.Coordinator
}

func NewAtomicCounter(cod datastore.Coordinator) *atomicCounter {

	min, max := cod.GetNextRange()

	return &atomicCounter{
		min: min,
		max: max,
		cod: cod,
	}
}

// Apply CAS
func (a *atomicCounter) next() int {

	a.min++

	// commit latest value of counter
	a.cod.Commit(a.min)

	return a.min
}

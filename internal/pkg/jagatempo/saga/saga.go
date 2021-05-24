package saga

import (
	"go.temporal.io/sdk/workflow"
)

type compensationFn func(workflow.Context) error

type Saga interface {
	AddCompensation(workflow.Context, compensationFn)
	Compensate(workflow.Context) error
}

type saga struct {
	compensationFnPool []compensationFn
}

func New() Saga {
	return &saga{}
}

func (s *saga) AddCompensation(ctx workflow.Context, fn compensationFn) {
	s.push(fn)
}

func (s *saga) Compensate(ctx workflow.Context) (err error) {
	for s.count() >= 1 {
		fn := s.pop()
		if err = fn(ctx); err != nil {
			return
		}
	}
	return
}

func (s *saga) push(fn ...compensationFn) {
	s.compensationFnPool = append(s.compensationFnPool, fn...)
}

func (s *saga) prepend(fn ...compensationFn) {
	s.compensationFnPool = append(fn, s.compensationFnPool...)
}

func (s *saga) pop() (fn compensationFn) {
	if len(s.compensationFnPool) == 0 {
		fn = func(ctx workflow.Context) error {
			return nil
		}
		return
	}

	fn = s.compensationFnPool[len(s.compensationFnPool)-1]
	s.compensationFnPool = s.compensationFnPool[:len(s.compensationFnPool)-1]
	return
}

func (s *saga) shift() (fn compensationFn) {
	if len(s.compensationFnPool) == 0 {
		fn = func(ctx workflow.Context) error {
			return nil
		}
		return
	}

	fn = s.compensationFnPool[0]
	s.compensationFnPool = s.compensationFnPool[1:]
	return
}

func (s *saga) peek() (fn compensationFn) {
	if len(s.compensationFnPool) == 0 {
		fn = func(ctx workflow.Context) error {
			return nil
		}
		return
	}

	return s.compensationFnPool[len(s.compensationFnPool)-1]
}

func (s *saga) pull(n int) (fn compensationFn) {
	if len(s.compensationFnPool) == 0 || len(s.compensationFnPool) < n {
		fn = func(ctx workflow.Context) error {
			return nil
		}
		return
	}

	return s.compensationFnPool[n]
}

func (s *saga) count() int {
	return len(s.compensationFnPool)
}

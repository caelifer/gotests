package scheduler

import (
	"github.com/caelifer/gotests/dnatransform/scheduler/job"
	"github.com/caelifer/gotests/dnatransform/scheduler/worker"
)

type Scheduler interface {
	Schedule(job.Job)
}

// Scheduler
type sched struct {
	pool chan worker.Worker
}

// NewSched(n int) creates new scheduler and initializes worker pool
func New(n int) Scheduler {
	s := new(sched)
	s.pool = make(chan worker.Worker, n)
	for i := 0; i < n; i++ {
		s.pool <- worker.New(i)
	}
	return s
}

func (s *sched) getWorker() worker.Worker {
	return <-s.pool
}

func (s *sched) putbackWorker(w worker.Worker) {
	select {
	case s.pool <- w:
		return
	default:
		// there should never been more workers that we have originaly created
		panic("pool is full")
	}
}

func (s *sched) Schedule(j job.Job) {
	w := s.getWorker()

	// Once we have a worker, run it in a separate goroutine
	go func() {
		w.Run(j)
		s.putbackWorker(w)
	}()
}

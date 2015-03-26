package worker

import "github.com/caelifer/gotests/dnatransform/scheduler/job"

// Worker interface
type Worker interface {
	Run(job.Interface)
}

// Worker implementation
type worker struct {
	id, jcount int
}

// Worker Run method syncronously executing provided Job.
func (w *worker) Run(j job.Interface) {
	w.jcount++
	// log.Printf("W[%02d] - running  job #%d", w.id, w.jcount)
	j()
	// log.Printf("W[%02d] - finished job #%d", w.id, w.jcount)
}

// Worker constructor
func New(id int) Worker {
	return &worker{id: id}
}

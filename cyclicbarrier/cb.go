package cyclicbarrier

import (
	"fmt"
	"sync"
)

type cb struct {
	sync.Mutex
	total   int
	pending int
	barrier chan struct{}
}

func (cb *cb) reset() {
	cb.Lock()
	defer cb.Unlock()

	cb.pending = cb.total
	close(cb.barrier)
	cb.barrier = make(chan struct{})
}

func (cb *cb) Await() {
	cb.Lock()
	cb.pending--
	pending := cb.pending
	barrier := cb.barrier
	cb.Unlock()

	if pending > 0 {
		<-barrier
	} else {
		fmt.Println("All parties have arrived to the barrier, lets play...")
		cb.reset()
	}
}

func New(parties int) *cb {
	cb := &cb{
		total:   parties,
		pending: parties,
		barrier: make(chan struct{}),
	}
	return cb
}

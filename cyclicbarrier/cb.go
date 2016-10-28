package cyclicbarrier

import (
	"fmt"
	"sync"
)

type cb struct {
	sync.Mutex
	total   int
	pending int
	barier  chan struct{}
}

func (cb *cb) reset() {
	cb.Lock()
	defer cb.Unlock()

	cb.pending = cb.total
	close(cb.barier)
	cb.barier = make(chan struct{})
}

func (cb *cb) Await() {
	cb.Lock()
	cb.pending--
	pending := cb.pending
	barier := cb.barier
	cb.Unlock()

	if pending > 0 {
		<-barier
	} else {
		fmt.Println("All parties has arrived to the barrier, lets play...")
		cb.reset()
	}
}

func New(parties int) *cb {
	cb := &cb{
		total:   parties,
		pending: parties,
		barier:  make(chan struct{}),
	}
	return cb
}

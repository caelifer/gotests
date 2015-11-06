package cyclicbarrier

import (
	"fmt"
	"sync"
)

type cb struct {
	sync.Mutex
	p int
	n int
	b chan struct{}
}

func (cb *cb) reset() {
	cb.Lock()
	defer cb.Unlock()

	cb.n = cb.p
	close(cb.b)
	cb.b = make(chan struct{})

}

func (cb *cb) Await() {
	cb.Lock()
	cb.n--
	n := cb.n
	b := cb.b
	cb.Unlock()

	if n > 0 {
		<-b
	} else {
		fmt.Println("All parties has arrived to the barrier, lets play...")
		cb.reset()
	}
}

func New(parties int) *cb {
	cb := &cb{
		p: parties,
		n: parties,
		b: make(chan struct{}),
	}
	return cb
}

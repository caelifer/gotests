package choir

import "sync"

type Choir struct {
	trigger chan struct{}
	wg      sync.WaitGroup
}

func New() *Choir {
	c := new(Choir)
	c.Reset()
	return c
}

func (c *Choir) PrepareRun(run, gen int, tasks []Task) {
	for _, task := range tasks {
		c.wg.Add(1)
		go func(t Task) {
			defer c.wg.Done()
			<-c.trigger
			t.Job(run, gen)
		}(task)
	}
}

func (c *Choir) Start() {
	close(c.trigger)
}

func (c *Choir) Wait() {
	c.wg.Wait()
}

func (c *Choir) Reset() {
	c.trigger = make(chan struct{})
}

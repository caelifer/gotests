package choir

import (
	"testing"
	"sync/atomic"
)

func TestFunction(t *testing.T) {
	tasks := setupTasks()
	c := New()
	if c == nil {
		t.Fatal("Expected New() to return *Choir object, got nil")
	}

	c.PrepareRun(0, 1, tasks)

	c.Start()
	if _, ok := <-c.trigger; ok {
		t.Error("Expected trigger channel to be closed after Start()")
	}

	c.Wait()
	if counter != int64(len(tasks)) {
		t.Errorf("Expected %d job to run to completion, got %d", len(tasks), counter)
	}

	c.Reset()
	func() {
		defer func(){
			if r := recover(); r != nil {
				t.Errorf("Excpected normal channel operation, got panic: %q", r)
			}
		}()

		go func(){
			c.trigger <- struct{}{}
		}()
		if _, ok := <- c.trigger; !ok {
			t.Error("Excpected good receive from trigger")
		}
	}()
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = new(Choir)
	}
}

func BenchmarkRaw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = &Choir{}
	}
}

func setupTasks() []Task {
	counter = 0 // reset counter
	return []Task{
		task{},
		task{},
	}
}

var counter int64

type task struct{}

func (t task) Job(run, gen int) {
	_ = atomic.AddInt64(&counter, 1)
}
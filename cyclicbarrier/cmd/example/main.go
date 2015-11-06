package main

import (
	"fmt"
	"sync"

	"github.com/caelifer/gotests/cyclicbarrier"
)

func main() {
	var wg sync.WaitGroup
	cb := cyclicbarrier.New(5)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			name := fmt.Sprintf("task #%d", id+1)
			for b := 0; b < 2; b++ {
				fmt.Printf("%s is waiting at the barrier #%d.\n", name, b+1)
				cb.Await()
				fmt.Printf("%s has crossed the barrier #%d.\n", name, b+1)
			}
		}(i)
	}
	wg.Wait()
}

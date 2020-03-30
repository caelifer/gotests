package main

import (
	"fmt"
	"sync"

	cb "github.com/caelifer/gotests/cyclicbarrier"
)

const (
	nPlayers  = 5
	nBarriers = 2
)

func main() {
	var wg sync.WaitGroup
	br := cb.New(nPlayers)

	wg.Add(nPlayers)
	for i := 0; i < nPlayers; i++ {
		go func(id int) {
			defer wg.Done()
			name := fmt.Sprintf("task #%d", id+1)
			for b := 0; b < nBarriers; b++ {
				fmt.Printf("%s is waiting at the barrier #%d.\n", name, b+1)
				br.Await()
				fmt.Printf("%s has crossed the barrier #%d.\n", name, b+1)
			}
		}(i)
	}
	wg.Wait()
}

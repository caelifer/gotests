package main

import (
	"fmt"

	"github.com/caelifer/gotests/choir"
)

const NumTasks = 5

func main() {
	c := choir.New()
	tasks := make([]choir.Task, NumTasks)

	for i := range tasks {
		tasks[i] = task{i, string('A' + i)}
	}

	for i := 0; i < 3; i++ {
		gen := i + 1
		c.PrepareRun(i, gen, tasks)
		fmt.Printf("Generation[%d]:\n", gen)
		c.Start()
		c.Wait()
		c.Reset()
		if i < 2 {
			fmt.Println("--------")
		}
	}
}

type task struct{
	id int
	name string
}

func (t task) Job(run, gen int) {
	fmt.Printf("task[%d] named:%q for run[%d] is done, generation[%d]\n", t.id, t.name, run, gen)
}

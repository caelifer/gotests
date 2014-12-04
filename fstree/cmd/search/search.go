package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/caelifer/dups/balancer"
	"github.com/caelifer/gotests/fstree"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	for _, root := range os.Args[1:] {
		// err := fstree.Walk(balancer.NewWorkQueue(runtime.NumCPU()*4), root, func(path string, dent fstree.Dirent, err error) error {
		err := fstree.Walk(balancer.NewWorkQueue(runtime.NumCPU()), root, func(path string, dent fstree.Dirent, err error) error {
			fmt.Printf("[%s] %s\n", dent.Type(), path)
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

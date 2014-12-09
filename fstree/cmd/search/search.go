package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/caelifer/dups/balancer"
	"github.com/caelifer/gotests/fstree"
)

var (
	// command line flags
	maxWorkers = flag.Int("jobs", 512, "Number of parallel threads")
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetPrefix(filepath.Base(os.Args[0]))
}

func main() {
	flag.Parse()

	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	for _, root := range roots {
		err := fstree.Walk(balancer.NewWorkQueue(*maxWorkers), root, func(path string, dent fstree.Dirent, err error) error {
			fmt.Printf("[%s] %s\n", dent.Type(), path)
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

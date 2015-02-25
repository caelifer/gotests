package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
)

const DEFAULT_PAGER = "/usr/bin/less"

func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal("Must provide file path")
	}

	// Get pager from the environment variable
	pager := os.Getenv("PAGER")
	if pager == "" {
		pager = DEFAULT_PAGER
		log.Println("[WARN] PAGER var is not set; using", pager)
	}

	// For each file path on the command line
	for _, fpath := range flag.Args() {
		input, err := os.Open(fpath)
		if err != nil {
			log.Fatal("Failed to ", err)
		}

		// Start pager
		cmd := exec.Command(pager)
		cmd.Stdout = os.Stdout // Make sure pager is printing to our STDOUT

		// Connect STDIN via pipe
		output, err := cmd.StdinPipe()
		if err != nil {
			log.Fatal("Failed to ", err)
		}

		// Run pager command
		if err := cmd.Start(); err != nil {
			log.Fatal("Failed to ", err)
		}

		// Pipe input into command's STDIN while counting bytes
		n, err := io.Copy(output, input)
		if err != nil {
			log.Fatal("I/O error ", err)
		}

		// Don't forget to close both output and input
		input.Close()
		output.Close()

		if err := cmd.Wait(); err != nil {
			log.Fatal("Failed to ", err)
		}

		// Display read bytes count
		log.Printf("[<< Processed %d bytes >>]\n", n)
	}
}

package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/caelifer/scheduler"
)

// Global scheduler
var GlobalScheduler scheduler.Scheduler

func main() {
	// Make sure none of the subroutines have access to the command line variables
	var (
		cpuprofile     = flag.String("cpuprofile", "", "write cpu profile to file")
		memprofile     = flag.String("memprofile", "", "write memory profile to file")
		workerCount    = flag.Int("jobs", runtime.NumCPU()*2, "Number of parallel workers")
		workQueueDepth = flag.Int("n", 64, "Number of queued jobs per worker")
	)

	// Parse commandline parameters
	flag.Parse()

	// Set # threads
	runtime.GOMAXPROCS(*workerCount)

	// CPU profile
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// Memory profile
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			pprof.WriteHeapProfile(f)
			f.Close()
		}()
	}

	// Pre-init GlobalScheduler so the workers can all be started before we proceed
	GlobalScheduler = scheduler.New(*workerCount, *workQueueDepth)

	for _, fpath := range flag.Args() {
		f, err := os.Open(fpath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		// Unzip
		z, err := gzip.NewReader(f)
		if err != nil {
			log.Fatal(err)
		}
		defer z.Close()

		// Send to our Chunker
		t0 := time.Now()
		res := chunker(labelFilter(z), 4096*4) // 16K chunks

		// Discard output but count bytes read
		var t1 time.Time
		bcount := 0
		for chunk := range res {
			if bcount == 0 {
				t1 = time.Now()
			}
			bcount += len(chunk)
		}
		// Capture duration
		d := time.Since(t0)

		log.Printf("Processed %d bytes in %s [%.3f MiB/s]; started getting results after %s.", bcount, d, float64(bcount)/(d.Seconds()*1024*1024), t1.Sub(t0))
	}

	// panic("Stack Trace:")
}

func labelFilter(r io.Reader) io.Reader {
	fr := &filtReader{out: make(chan []byte, 100)}

	// Wrap around buffered reader
	br := bufio.NewReader(r)

	// Filter stream
	go func() {
		for {
			// Skip first line
			line, err := br.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					close(fr.out)
					return
				} else {
					log.Fatal(err)
				}
			}

			if line[0] == '>' {
				// Print label and skip
				fmt.Print(string(line))
			} else {
				line = append(line, '\n')
				fr.out <- line
			}
		}
	}()

	return fr
}

type filtReader struct {
	out chan []byte
	buf bytes.Buffer
	eof bool
}

func (fr *filtReader) Read(buf []byte) (int, error) {
	osize := len(buf)

	// Drain our buffer if not empty
	if fr.buf.Len() > 0 {
		return fr.buf.Read(buf)
	}

	// If EOF already reached
	if fr.eof {
		return 0, io.EOF
	}

	// Get new chunk(s)
	ngot := 0
	for {
		data, ok := <-fr.out
		if data != nil {
			n, err := fr.buf.Write(data)
			if err != nil {
				log.Fatal(err)
			}
			ngot += n
		}

		if !ok {
			fr.eof = true
			break
		}

		// Have enough
		if ngot >= osize {
			break
		}
	}
	return fr.buf.Read(buf)
}

func chunker(r io.Reader, bufsize int) <-chan []byte {
	head := make(chan []byte)

	// Schedule in separate worker
	GlobalScheduler.Schedule(func() {
		for {
			chunk := make([]byte, bufsize) // 4K read chunk

			n, err := r.Read(chunk)
			if err != nil {
				if err == io.EOF {
					// Done reading
					// log.Println("Done with Read()")
					close(head) // Close curent head because we are done
					return
				}
				log.Fatal(err)
			}
			// log.Printf("Read() %d bytes", n)

			// Got another chunk, build daisy-chain
			tail := piper(head, chunk[:n], GlobalScheduler)
			head = tail
		}
	})

	return head
}

// Piper()
func piper(head chan<- []byte, chunk []byte, s scheduler.Scheduler) chan []byte {
	// Create results chan
	res := make(chan []byte)
	// Schedule work, this may block if all workers are busy
	s.Schedule(func() {
		// log.Printf("sent: %+v", chunk)
		Transform(chunk)
		// log.Printf("recv: %+v", chunk)

		// capture res channel
		res <- chunk
		close(res)
	})

	// Daisy-chain through our pipe in a goroutine
	tail := make(chan []byte)
	go func(r chan []byte) {
		for v := range r {
			// Feed it to head first
			head <- v
		}

		// Once done with our payload, pipe through all the rest from tail
		for v := range tail {
			head <- v
		}

		// We are done, close head
		close(head)
	}(res)

	return tail
}

// Fast DNA transformer
var mod = byte(9)
var dnsalpha = []byte{0, '\n', 'A', 'T', 'C', 0, 'N', 0, 'G', 0}
var dnstrans = []byte{0, '\n', 'T', 'A', 'G', 0, 'N', 0, 'C', 0}

func Transform(c []byte) {
	for i, v := range c {
		if in(v) {
			c[i] = dnstrans[v%mod]
		} else {
			log.Fatalf("[WARN] bad symbol '%c' in %q", v, string(c))
		}
	}
}

func in(t byte) bool {
	for _, b := range dnsalpha {
		if b == t {
			return true
		}
	}
	return false
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// func _main() {
//
// 	tests := []struct {
// 		t, e string
// 	}{
// 		// 		"",
// 		// 		"1",
// 		// 		"12",
// 		// 		"12345",
// 		// 		"123456",
// 		// 		"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
// 		{
// 			"GTGATTGTTGTTAGGTTCTTAGCTGCTTCTGAAAAATGGGGTGATAATCTTAGAAGGACT",
// 			"CACTAACAACAATCCAAGAATCGACGAAGACTTTTTACCCCACTATTAGAATCTTCCTGA",
// 		},
// 		{
// 			"CTAAACCATATGNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNN",
// 			"GATTTGGTATACNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNN",
// 		},
// 	}
//
// 	for i, s := range tests {
// 		res := Chunker(strings.NewReader(s.t), 10240*1024) // n-KiB chunks
//
// 		buf := make([]byte, 0, 100)
// 		for r := range res {
// 			buf = append(buf, r...)
// 		}
// 		t := string(buf)
//
// 		var status string
// 		if t == s.e {
// 			status = "Passed"
// 		} else {
// 			status = "FAILED"
// 		}
// 		fmt.Printf("Test #%d: %v\n", i, status)
// 		fmt.Printf("\tTested  : %q\n", s.t)
// 		fmt.Printf("\tReceived: %q\n", t)
// 		fmt.Printf("\tExpected: %q\n", s.e)
// 		fmt.Println()
// 	}
// }

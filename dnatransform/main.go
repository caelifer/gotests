package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"time"
)

// Worker's payload signature
type Job func()

// Worker interface
type Worker interface {
	Run(Job)
}

// Worker implementation
type worker struct {
	id, jcount int
}

// Worker Run method syncronously executing provided Job.
func (w *worker) Run(j Job) {
	w.jcount++
	// log.Printf("W[%02d] - running  job #%d", w.id, w.jcount)
	j()
	// log.Printf("W[%02d] - finished job #%d", w.id, w.jcount)
}

// Worker constructor
func NewWorker(id int) Worker {
	return &worker{id: id}
}

// Worker Pool

// Scheduler
type Sched struct {
	pool chan Worker
}

// NewSched(n int) creates new scheduler and initializes worker pool
func NewSched(n int) *Sched {
	s := new(Sched)
	s.pool = make(chan Worker, n)
	for i := 0; i < n; i++ {
		s.pool <- NewWorker(i)
	}
	return s
}

func (s *Sched) getWorker() Worker {
	return <-s.pool
}

func (s *Sched) returnWorker(w Worker) {
	select {
	case s.pool <- w:
		return
	default:
		// there should never been more workers that we have originaly created
		panic("pool is full")
	}
}

func (s *Sched) Schedule(j Job) {
	w := s.getWorker()

	// Once we have a worker, run it in a separate goroutine
	go func() {
		w.Run(j)
		s.returnWorker(w)
	}()
}

// Piper()
func Piper(head chan<- []byte, chunk []byte, s *Sched) chan []byte {
	// Create results chan
	res := make(chan []byte)
	// Schedule work, this may block if all workers are busy
	s.Schedule(func() {
		// log.Printf("sent: %+v", chunk)
		transform(chunk)
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

var NWorkers = runtime.NumCPU()

func Chunker(r io.Reader, bufsize int) <-chan []byte {
	s := NewSched(NWorkers)
	head := make(chan []byte)

	// Schedule in separate worker
	s.Schedule(func() {
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
			tail := Piper(head, chunk[:n], s)
			head = tail
		}
	})

	return head
}

var mod = byte(9)
var dnsalpha = []byte{0, '\n', 'A', 'T', 'C', 0, 'N', 0, 'G', 0}
var dnstrans = []byte{0, '\n', 'T', 'A', 'G', 0, 'N', 0, 'C', 0}

func in(t byte) bool {
	for _, b := range dnsalpha {
		if b == t {
			return true
		}
	}
	return false
}

func transform(c []byte) {
	for i, v := range c {
		if in(v) {
			c[i] = dnstrans[v%mod]
		} else {
			log.Fatalf("[WARN] bad symbol '%c' in %q", v, string(c))
		}
	}
}

func main() {
	// Set # threads
	runtime.GOMAXPROCS(NWorkers)

	for _, fpath := range os.Args[1:] {
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

		// Wrap buffered reader
		r := bufio.NewReader(z)

		// Skip first line
		line, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		// Print label
		fmt.Print(line)

		// Send to our Chunker
		t0 := time.Now()
		res := Chunker(r, 4096*4) // 16K chunks

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

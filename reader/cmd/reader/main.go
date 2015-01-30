package main

import (
	"fmt"

	sbuf "github.com/caelifer/gotests/reader/stringbuffer"
)

func main() {
	const (
		size = 5
		test = "qwertyuiopasdfghjkl"
	)

	var (
		buf = make([]byte, size)
		r   = sbuf.NewReader(test)
	)

	for i := 0; i < len(test)/size+2; i++ {
		n, err := r.Read(buf)
		fmt.Printf("n = %d, err = %-5v, str = %-8q (buf:%+3v) [r's sbuf: len = %02d, cap = %02d]\n", n, err, string(buf[:n]), buf, r.Len(), r.Cap())
	}
}

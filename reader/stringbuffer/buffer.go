package stringbuffer

import "io"

type Reader struct {
	data []byte
}

func NewReader(toRead string) *Reader {
	return &Reader{[]byte(toRead)}
}

func (r Reader) Len() int {
	return len(r.data)
}

func (r Reader) Cap() int {
	return cap(r.data)
}

func (r Reader) eof() bool {
	return len(r.data) == 0
}

func (r *Reader) readByte() byte {
	// this function assumes that eof() check was done before
	b := r.data[0]
	r.data = r.data[1:]
	return b
}

func (r *Reader) Read(p []byte) (n int, err error) {
	if r.eof() {
		err = io.EOF
		return
	}

	if l := len(p); l > 0 {
		for n < l {
			p[n] = r.readByte()
			n++
			if r.eof() {
				// free memory
				r.data = []byte{}
				break
			}
		}
	}
	return
}

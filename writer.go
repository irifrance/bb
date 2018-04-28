package bb

import (
	"io"
)

type Writer struct {
	t *T
	io.Writer
	err    error
	nFlush int64
}

func (w *Writer) Bump() {
	w.t.Bump()
}

func (w *Writer) WriteBit(b byte) error {
	if !w.t.has(1) && !w.swap() {

		return w.err
	}
	w.t.WriteBit(b)
	return nil
}

func (w *Writer) WriteBool(b bool) error {
	if !w.t.has(1) && !w.swap() {
		return w.err
	}
	w.t.WriteBool(b)
	return nil
}

func (w *Writer) WriteBits(b byte, n int) error {
	if n > 8 {
		return OutOfBounds
	}
	if !w.t.has(n) && !w.swap() {
		return w.err
	}
	w.t.WriteBits(b, n)
	return nil
}

func (w *Writer) Write16(v uint16, n int) error {
	if n > 16 {
		return OutOfBounds
	}
	t := w.t
	if !t.has(n) && !w.swap() {
		return w.err
	}
	t.Write16(v, n)
	return nil
}

func (w *Writer) Write32(v uint32, n int) error {
	if n > 32 {
		return OutOfBounds
	}
	t := w.t
	if !t.has(n) && !w.swap() {
		return w.err
	}
	t.Write32(v, n)
	return nil
}

func (w *Writer) Write64(v uint64, n int) error {
	if n > 64 {
		return OutOfBounds
	}
	t := w.t
	if !t.has(n) && !w.swap() {
		return w.err
	}
	t.Write64(v, n)
	return nil
}

func (w *Writer) Flush() error {
	w.Bump()
	w.swap()
	return w.err
}

func (w *Writer) BitsWritten() int64 {
	return 8*w.nFlush + int64(w.t.i)
}

func (w *Writer) swap() bool {
	t := w.t
	p, m := int(t.i/8), t.i%8
	var nw int
	var e error
	t.i = m
	nw, e = w.Write(t.d[:p])
	if e != nil {
		w.err = e
		return false
	}
	if nw != p {
		panic("nw")
	}
	j := 0
	if p < len(t.d) {
		t.d[0] = t.d[p]
		j++
	}
	for i := j; i < len(t.d); i++ {
		t.d[i] = 0
	}
	w.nFlush += int64(nw)
	return true
}

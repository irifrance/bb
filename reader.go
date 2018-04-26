package bb

import (
	"io"
)

type Reader struct {
	t *T
	io.Reader
	err error
}

func (r *Reader) Bump() {
	r.t.Bump()
}

func (r *Reader) ReadBit() (byte, error) {
	if !r.t.has(1) && !r.swap(1) {
		return 0, r.err
	}
	return r.t.ReadBit(), nil
}

func (r *Reader) ReadBool() (bool, error) {
	if !r.t.has(1) && !r.swap(1) {
		return false, r.err
	}
	return r.t.ReadBool(), nil
}

func (r *Reader) ReadBits(n int) (byte, error) {
	if n > 8 {
		return 0, OutOfBounds
	}
	if !r.t.has(n) && !r.swap(n) {
		return 0, r.err
	}
	return r.t.ReadBits(n), nil
}

func (r *Reader) Read16(n int) (uint16, error) {
	if n > 16 {
		return 0, OutOfBounds
	}
	t := r.t
	if !t.has(n) && !r.swap(n) {
		return 0, r.err
	}
	return t.Read16(n), nil
}

func (r *Reader) Read32(n int) (uint32, error) {
	if n > 32 {
		return 0, OutOfBounds
	}
	t := r.t
	if !t.has(n) && !r.swap(n) {
		return 0, r.err
	}
	return t.Read32(n), nil
}

func (r *Reader) Read64(n int) (uint64, error) {
	if n > 64 {
		return 0, OutOfBounds
	}
	t := r.t
	if !t.has(n) && !r.swap(n) {
		return 0, r.err
	}
	return t.Read64(n), nil
}

func (r *Reader) swap(n int) bool {
	t := r.t
	rem := t.BitsRemaining()
	p, m := int(t.i/8), t.i%8
	q := len(t.d) - p
	copy(t.d, t.d[p:])
	t.i = m
	nRead, e := r.Read(t.d[q:])
	if e != nil {
		r.err = e
	}
	t.d = t.d[:q+nRead]
	return nRead*8+rem >= n
}

package bb

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Reader struct {
	t *T
	io.Reader
	err     error
	srcRead int64
}

func NewReader(r io.Reader, sz int) *Reader {
	if sz < 8 {
		sz = 8
	}
	t := New(sz)
	t.i = uint(sz * 8)
	return &Reader{t: t, Reader: r}
}

func ReaderFromSlice(r io.Reader, sl []byte) *Reader {
	if len(sl) < 8 {
		tmp := make([]byte, 8)
		copy(tmp, sl)
		sl = tmp
	}
	t := FromSlice(sl)
	t.i = uint(len(sl) * 8)
	return &Reader{t: t, Reader: r}
}

func (r *Reader) T() *T {
	return r.t
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

func (r *Reader) ReadVarint() (int64, error) {
	buf, e := r.readV()
	if e != nil {
		return 0, e
	}
	val, n := binary.Varint(buf)
	if n <= 0 {
		return 0, fmt.Errorf("binary.Varint gave %d\n", n)
	}
	return val, nil
}

func (r *Reader) ReadUvarint() (uint64, error) {
	buf, e := r.readV()
	if e != nil {
		return 0, e
	}
	val, n := binary.Uvarint(buf)
	if n <= 0 {
		return 0, fmt.Errorf("binary.Varint gave %d\n", n)
	}
	return val, nil
}

func (r *Reader) readV() ([]byte, error) {
	buf := make([]byte, 8)
	i := 0
	var b byte
	var e error
	for i < 8 {
		b, e = r.ReadBits(8)
		if e != nil {
			return nil, e
		}
		buf[i] = b
		i++
		if b&128 == 0 {
			break
		}
	}
	return buf, nil
}

func (r *Reader) BitsRead() int64 {
	return 8*r.srcRead - int64(r.t.BitsRemaining())
}

func (r *Reader) swap(n int) bool {
	t := r.t
	rem := t.BitsRemaining()
	p, m := int(t.i/8), t.i%8
	q := len(t.d) - p
	copy(t.d, t.d[p:])
	t.i = m
	n, nRead := 0, 0
	for nRead < p {
		n, r.err = r.Read(t.d[q:])
		nRead += n
		q += n
		if r.err != nil {
			break
		}
	}
	t.d = t.d[:q]
	r.srcRead += int64(nRead)
	if nRead == 0 && r.err != nil {
		return false
	}
	return nRead*8+rem >= n
}

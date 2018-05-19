// Copyright 2018 Iri France SAS. All rights reserved.  Use of this source code
// is governed by a license that can be found in the License file.

package bb

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Buffer provides bit io operation with errors for
// compatability with io readers/writers.
type Buffer struct {
	t       T
	err     error
	r       io.Reader
	w       io.Writer
	srcRead int64
	nFlush  int64
}

func NewBuffer(n int) *Buffer {
	return &Buffer{t: *New(n)}
}

func NewBufferSlice(b []byte) *Buffer {
	return &Buffer{t: *FromSlice(b)}
}

func (b *Buffer) ReadBool() (bool, error) {
	if !b.t.has(1) && !b.rSwap(1) {
		return false, b.err
	}
	return b.t.ReadBool(), nil
}

func (b *Buffer) ReadBit() (byte, error) {
	if !b.t.has(1) && !b.rSwap(1) {
		return 0, io.EOF
	}
	return b.t.ReadBit(), nil
}

func (b *Buffer) ReadBits(n int) (byte, error) {
	if n > 8 {
		return 0, OutOfBounds
	}
	if !b.t.has(n) && !b.rSwap(n) {
		return 0, b.err
	}
	return b.t.ReadBits(n), nil
}

func (b *Buffer) ReadByte() (byte, error) {
	return b.ReadBits(8)
}

func (b *Buffer) Read16(n int) (uint16, error) {
	if n > 16 {
		return 0, OutOfBounds
	}
	if !b.t.has(n) && !b.rSwap(n) {
		return 0, b.err
	}
	return b.t.Read16(n), nil
}

func (b *Buffer) Read32(n int) (uint32, error) {
	if n > 32 {
		return 0, OutOfBounds
	}
	t := &b.t
	if !t.has(n) && !b.rSwap(n) {
		return 0, b.err
	}
	return t.Read32(n), nil
}

func (b *Buffer) Read64(n int) (uint64, error) {
	if n > 64 {
		return 0, OutOfBounds
	}
	t := &b.t
	if !t.has(n) && !b.rSwap(n) {
		return 0, b.err
	}
	return t.Read64(n), nil
}

func (b *Buffer) Bump() {
	b.t.Bump()
}

func (b *Buffer) ReadVarint() (int64, error) {
	buf, e := b.readV()
	if e != nil {
		return 0, e
	}
	val, n := binary.Varint(buf)
	if n <= 0 {
		return 0, fmt.Errorf("binary.Varint gave %d\n", n)
	}
	return val, nil
}

func (b *Buffer) ReadUvarint() (uint64, error) {
	buf, e := b.readV()
	if e != nil {
		return 0, e
	}
	val, n := binary.Uvarint(buf)
	if n <= 0 {
		return 0, fmt.Errorf("binary.Varint gave %d\n", n)
	}
	return val, nil
}

func (bb *Buffer) readV() ([]byte, error) {
	var buf [8]byte
	i := 0
	var b byte
	var e error
	for i < 8 {
		b, e = bb.ReadBits(8)
		if e != nil {
			return nil, e
		}
		buf[i] = b
		i++
		if b&128 == 0 {
			break
		}
	}
	return buf[:], nil
}

func (b *Buffer) BitsRead() int64 {
	return b.srcRead*8 + int64(b.t.BitLen())
}

func (b *Buffer) SeekBit(i int) {
	b.t.SeekBit(i)
}

func (b *Buffer) Bytes() []byte {
	return b.t.Bytes()
}

func (bb *Buffer) WriteBit(b byte) error {
	if !bb.t.has(1) && !bb.wSwap() {
		return bb.err
	}
	bb.t.WriteBit(b)
	return nil
}

func (bb *Buffer) WriteBool(b bool) error {
	if !bb.t.has(1) && !bb.wSwap() {
		return bb.err
	}
	t := &bb.t
	t.WriteBool(b)
	return nil
}

func (bb *Buffer) WriteBits(b byte, n int) error {
	if n > 8 {
		return OutOfBounds
	}
	if !bb.t.has(n) && !bb.wSwap() {
		return bb.err
	}
	bb.t.WriteBits(b, n)
	return nil
}

func (bb *Buffer) WriteByte(b byte) error {
	return bb.WriteBits(b, 8)
}

func (b *Buffer) Write16(v uint16, n int) error {
	if n > 16 {
		return OutOfBounds
	}
	if !b.t.has(n) && !b.wSwap() {
		return b.err
	}
	b.t.Write16(v, n)
	return nil
}

func (b *Buffer) Write32(v uint32, n int) error {
	if n > 32 {
		return OutOfBounds
	}
	if !b.t.has(n) && !b.wSwap() {
		return b.err
	}
	b.t.Write32(v, n)
	return nil
}

func (b *Buffer) Write64(v uint64, n int) error {
	if n > 64 {
		return OutOfBounds
	}
	if !b.t.has(n) && !b.wSwap() {
		return b.err
	}
	b.t.Write64(v, n)
	return nil
}

func (b *Buffer) WriteVarint(val int64) error {
	buf := make([]byte, 8)
	n := binary.PutVarint(buf, val)
	for i := 0; i < n; i++ {
		if e := b.WriteBits(buf[i], 8); e != nil {
			return e
		}
	}
	return nil
}

func (b *Buffer) WriteUvarint(val uint64) error {
	buf := make([]byte, 8)
	n := binary.PutUvarint(buf, val)
	for i := 0; i < n; i++ {
		if e := b.WriteBits(buf[i], 8); e != nil {
			return e
		}
	}
	return nil
}

func (b *Buffer) Flush() error {
	if b.w == nil {
		return NoWriterError
	}
	b.Bump()
	b.wSwap()
	return b.err
}

func (b *Buffer) BitsWritten() int64 {
	return b.nFlush*8 + int64(b.t.i)
}

func (b *Buffer) rSwap(n int) bool {
	if b.r == nil {
		b.err = io.EOF
		return false
	}
	t := &b.t
	rem := t.BitsRemaining()
	p, m := int(t.i/8), t.i%8
	q := len(t.d) - p
	copy(t.d, t.d[p:])
	t.i = m
	n, nRead := 0, 0
	for nRead < p {
		n, b.err = b.r.Read(t.d[q:])
		nRead += n
		q += n
		if b.err != nil {
			break
		}
	}
	t.d = t.d[:q]
	b.srcRead += int64(nRead)
	if nRead == 0 && b.err != nil {
		return false
	}
	return nRead*8+rem >= n
}

func (b *Buffer) wSwap() bool {
	if b.w == nil {
		return true
	}
	t := &(b.t)
	p, m := int(t.i/8), t.i%8
	var nw int
	var e error
	t.i = m
	nw, e = b.w.Write(t.d[:p])
	if e != nil {
		b.err = e
		return false
	}
	j := 0
	if p < len(t.d) {
		t.d[0] = t.d[p]
		j++
	}
	for i := j; i < len(t.d); i++ {
		t.d[i] = 0
	}
	b.nFlush += int64(nw)
	return true
}

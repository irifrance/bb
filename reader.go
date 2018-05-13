package bb

import "io"

type Reader interface {
	Bump()
	ReadBit() (byte, error)
	ReadBool() (bool, error)
	ReadBits(n int) (byte, error)
	Read16(n int) (uint16, error)
	Read32(n int) (uint32, error)
	Read64(n int) (uint64, error)
	ReadVarint() (int64, error)
	ReadUvarint() (uint64, error)
}

func NewReader(r io.Reader, sz int) Reader {
	b := NewBuffer(sz)
	b.r = r
	b.t.i = uint(sz) * 8
	return b
}

func NewBufferReader(sz int) Reader {
	return NewBuffer(sz)
}

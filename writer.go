package bb

import "io"

type Writer interface {
	Bump()
	WriteBit(b byte) error
	WriteBits(b byte, n int) error
	Write16(v uint16, n int) error
	Write32(v uint32, n int) error
	Write64(v uint64, n int) error
	WriteVarint(val int64) error
	WriteUvarint(val uint64) error
	Flush() error
	BitsWritten() int64
}

func NewWriter(w io.Writer, sz int) Writer {
	b := NewBuffer(sz)
	b.w = w
	return b
}

func NewWriterBuffer(sz int) Writer {
	b := NewBuffer(sz)
	return b
}

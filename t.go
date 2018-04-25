package bb

type T struct {
	D []byte
	i uint
}

func FromSlice(d []byte) *T {
	return &T{D: d}
}

func New(sz int) *T {
	return &T{D: make([]byte, sz)}
}

func (b *T) Bump() {
	m := b.i % 8
	if m != 0 {
		b.i += m
	}
}

func (b *T) ByteLen() int {
	m := b.i % 8
	if m == 0 {
		return int(b.i / 8)
	}
	return int(b.i/8 + 1)
}

func (b *T) BitLen() int {
	return int(b.i)
}

func (b *T) Rewind() {
	b.i = 0
}

func (b *T) WriteBool(v bool) {
	if v {
		b.WriteBit(1)
	} else {
		b.WriteBit(0)
	}
}

func (b *T) ReadBool() bool {
	return b.ReadBit() == 1
}

func (b *T) WriteLen(v uint64, n int) {
	j := uint(0)
	N := uint(n)
	for j < N {
		k, o := b.i/8, b.i%8
		if o == 0 && j+8 < N {
			b.D[k] = byte((v >> j) & 0xFF)
			j += 8
			b.i += 8
			continue
		}
		b.WriteBool((v>>j)&1 != 0)
		j++
	}
}

func (b *T) WriteBits(d byte, n int) {
	if n > 8 {
		panic("8")
	}
	N := uint(n)
	d &= (1 << N) - 1 // sanitize
	k, off := b.i/8, b.i%8
	b.D[k] |= d << off
	if off+N <= 8 {
		b.i += N
		return
	}
	b.i += 8 - off
	b.WriteBits(d>>(8-off), int((off+N)-8))
}

func (b *T) ReadBits(n int) byte {
	if n > 8 {
		panic("8")
	}
	N := uint(n)
	k, off := b.i/8, b.i%8
	res := b.D[k] >> off
	m := 8 - off
	if m == N {
		b.i += m
		return res
	}
	if m < N {
		b.i += m
		return res | (b.ReadBits(int(N-m)) << m)
	}
	b.i += N
	return res & ((1 << N) - 1)
}

func (b *T) ReadLen(n int) uint64 {
	if n > 64 {
		panic("64")
	}
	res := uint64(0)
	j := uint(0)
	N := uint(n)
	for j < N {
		k, o := b.i/8, b.i%8
		if o == 0 && j+8 < N {
			res |= uint64(b.D[k]) << j
			j += 8
			b.i += 8
			continue
		}
		res |= uint64(b.ReadBit()) << j
		j++
	}
	return res
}

func (b *T) ReadBit() byte {
	i := b.i
	b.i++
	j, o := i/8, i%8
	if (b.D[j]>>o)&1 == 1 {
		return 1
	}
	return 0
}

func (b *T) WriteBit(d byte) {
	if d == 0 {
		b.i++
		return
	}
	i := b.i
	j, o := i/8, i%8
	b.D[j] |= (1 << o)
	b.i++
}

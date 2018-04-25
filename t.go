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

// Bump if bit index is not aligned to next byte
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

func (b *T) Write32(v uint32, n int) {
	if n > 32 {
		panic("32")
	}
	b.Write64(uint64(v), n)
}

func (b *T) Write16(v uint16, n int) {
	if n > 16 {
		panic("16")
	}
	b.Write64(uint64(v), n)
}

func (b *T) Write64(v uint64, n int) {
	if n > 64 {
		panic("64")
	}
	j := uint(0)
	var k, o, m uint
	N := uint(n)
	for j < N {
		k, o = b.i/8, b.i%8
		if o == 0 && j+8 < N {
			b.D[k] = byte((v >> j) & 0xFF)
			j += 8
			b.i += 8
			continue
		}
		m = 8 - o
		if j+m > N {
			m = N - j
		}
		b.WriteBits(byte(v>>j), int(m))
		j += m
	}
}

func (b *T) Read32(n int) uint32 {
	if n > 32 {
		panic("32")
	}
	return uint32(b.Read64(n))
}

func (b *T) Read16(n int) uint16 {
	if n > 16 {
		panic("16")
	}
	return uint16(b.Read64(n))
}

func (b *T) Read64(n int) uint64 {
	if n > 64 {
		panic("64")
	}
	res := uint64(0)
	var j, k, o, m uint
	N := uint(n)
	for j < N {
		k, o = b.i/8, b.i%8
		if o == 0 && j+8 < N {
			res |= uint64(b.D[k]) << j
			j += 8
			b.i += 8
			continue
		}
		m = 8 - o
		if j+m > N {
			m = N - j
		}
		res |= uint64(b.ReadBits(int(m))) << j
		j += m
	}
	return res
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

// Copyright 2018 Iri France SAS. All rights reserved.  Use of this source code
// is governed by a license that can be found in the License file.

package bb

import (
	"math/rand"
	"testing"
)

func TestBBLen(t *testing.T) {
	b := New(8192)
	N := 256
	d := make([]uint64, N)
	shifts := make([]int, N)
	for i := 0; i < N; i++ {
		s := rand.Intn(16)
		shifts[i] = s
		n := uint64(rand.Intn(1 << uint(s)))
		d[i] = n
		b.Write64(n, s)
	}
	c := FromSlice(b.d)
	for i := 0; i < N; i++ {
		v := c.Read64(shifts[i])
		if v != d[i] {
			t.Errorf("%d != %d\n", v, d[i])
		}
	}
	c.Bump()
	t.Logf("byteLen %d\n", c.ByteLen())
}

func TestRWBits(t *testing.T) {
	N := 16
	buf := make([]byte, N)
	for i := range buf {
		buf[i] = byte(rand.Intn(8))
	}
	fub := make([]byte, N)
	n := 0
	src, dst := FromSlice(buf), FromSlice(fub)
	for n < N*8 {
		j := rand.Intn(9)
		if n+j > N*8 {
			j = N*8 - n
		}
		b := src.ReadBits(j)
		if src.BitLen() != n+j {
			t.Errorf("src bitlen got %d expected %d\n", src.BitLen(), n+j)
		}
		dst.WriteBits(b, j)
		if dst.BitLen() != n+j {
			t.Errorf("dst bitlen got %d expected %d\n", dst.BitLen(), n+j)
		}
		n += j
	}
	for i := range buf {
		if buf[i] != fub[i] {
			t.Errorf("%d: %d != %d\n", i, buf[i], fub[i])
		}
	}
}

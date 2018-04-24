package bb

import (
	"math/rand"
	"testing"
)

func TestBBLen(t *testing.T) {
	b := New(8192)
	N := 1024
	d := make([]uint64, N)
	shifts := make([]int, N)
	for i := 0; i < N; i++ {
		s := rand.Intn(16)
		shifts[i] = s
		n := uint64(rand.Intn(1 << uint(s)))
		d[i] = n
		b.WriteLen(n, s)
	}
	c := FromSlice(b.D)
	for i := 0; i < N; i++ {
		v := c.ReadLen(shifts[i])
		if v != d[i] {
			t.Errorf("%d != %d\n", v, d[i])
		}
	}
	c.Bump()
	t.Logf("byteLen %d\n", c.ByteLen())
}

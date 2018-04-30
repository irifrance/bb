package bb

import (
	"bytes"
	"testing"
)

func TestReaderWriter(t *testing.T) {
	bio := New(10)
	bub := bytes.NewBuffer(nil)
	bw := NewWriter(bub, 16)
	M := 5
	N := M * 8
	n := 0
	twiddle := false
	a := byte(1<<5) - 1
	var b byte
	for n < N {
		if twiddle {
			b = a
		} else {
			b = 0
		}
		if e := bw.WriteBits(b, 5); e != nil {
			t.Error(e)
		}

		n += 5
		twiddle = !twiddle
	}
	bw.Flush()
	n = 0
	bio = FromSlice(bub.Bytes())
	twiddle = false
	for n < N {
		c := bio.ReadBits(5)
		if twiddle && c != a {
			t.Errorf("1 wrong side of twiddle at %d: %d rem %d\n", n, c, bio.BitsRemaining())
		} else if !twiddle && c != 0 {
			t.Errorf("2 wrong side twiddle at %d: %d rem %d\n", n, c, bio.BitsRemaining())
		}
		n += 5
		twiddle = !twiddle
	}
}

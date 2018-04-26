package bb

import (
	"bytes"
	"io"
	"testing"
)

func TestReaderReadBits(t *testing.T) {
	for _, i := range []int{1, 3, 4, 6, 7, 8} {
		testReadBits(t, i)
	}
}

func testReadBits(t *testing.T, i int) {
	N := i * 64
	n := 0
	ones := false
	bio := New(N / 8)
	for n < N {
		if ones {
			bio.WriteBool(true)
		} else {
			bio.WriteBool(false)
		}
		if n%i == 0 {
			ones = !ones
		}
		n++
	}

	rbio := New(16)
	r, _ := rbio.Reader(bytes.NewBuffer(bio.Bytes()))
	n = 0
	ones = false
	for n < N {
		v, e := r.ReadBool()
		if e != nil {
			t.Error(e)
		}
		if v != ones {
			t.Errorf("with %d at %d expected %t got %t\n", i, n, ones, v)
		}
		if n%i == 0 {
			ones = !ones
		}
		n++
	}
	_, e := r.ReadBool()
	if e != io.EOF {
		t.Errorf("no EOF")
	}
}

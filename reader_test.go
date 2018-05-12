package bb

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestReaderReadBits(t *testing.T) {
	for _, i := range []int{1, 3, 4, 6, 7, 8} {
		testReadBits(t, i)
		testReadWBits(t, i)
		//testRead64(t, i)
	}
}

func TestReadVarw(t *testing.T) {

}

func testRead64(t *testing.T, i int) {
	fmt.Printf("tstRead64 %d\n", i)
	N := i * 128
	n := 0
	bio := New(N / 8)
	j := 0
	for n < N {
		v := uint(j % i)
		j++
		if j%2 == 0 {
			bio.Write64(uint64(1)<<v, j%64)
		} else {
			bio.Write64((uint64(1)<<v)|1, j%64)
		}
		n += j % 64
	}
	bio.Bump()
	r := NewReader(bytes.NewBuffer(bio.Bytes()), 16)
	n, j = 0, 0
	for n < N {
		v := uint(j % i)
		j++
		w, _ := r.Read64(j % 64)
		ref := uint64(1) << v
		if j%2 != 0 {
			ref |= 1
		}
		if w != ref {
			t.Fatalf("read64(j=%d,n=%d) decoded %d expected %d\n", j, n, w, ref)
		}
		n += j % 64
	}
}

func testReadBits(t *testing.T, i int) {
	N := i * 64
	n := 0
	ones := false
	bio := NewBuffer(N / 8)
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

	bio.SeekBit(0)
	var r Reader
	r = bio
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

func testReadWBits(t *testing.T, i int) {
	N := i * 64
	n := 0
	ones := false
	bio := NewBuffer(1)
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
	bio.Bump()

	var r Reader
	bio.SeekBit(0)
	r = bio
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
}

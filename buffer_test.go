// Copyright 2018 Iri France SAS. All rights reserved.  Use of this source code
// is governed by a license that can be found in the License file.

package bb

import (
	"bytes"
	"io"
	"math/rand"
	"testing"
)

func TestBackedBuffer(t *testing.T) {
	wBack := bytes.NewBuffer(nil)
	w := NewWriter(wBack, 16)

	N := 16
	bits := make([]bool, N)
	for i := 0; i < N; i++ {
		v := rand.Intn(3) == 1
		w.WriteBool(v)
		bits[i] = v
	}
	w.Flush()
	rBack := bytes.NewBuffer(wBack.Bytes())
	r := NewReader(rBack, 16)
	for i := 0; i < N; i++ {
		v, _ := r.ReadBool()
		if v != bits[i] {
			t.Errorf("mismatch at %d\n", i)
		}
	}
	_, e := r.ReadBool()
	if e != io.EOF {
		t.Errorf("expected EOF")
	}
}

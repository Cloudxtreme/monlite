// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Start date:        2014-04-01
// Last modification: 2014-

package round

import (
	"testing"
)

type test struct {
	A    float64
	B    float64
	Prec int
}

var trs []test = []test{
	{1.0, 1.0, 2},
	{0.0, 0.0, 2},
	{1.5, 1.5, 2},
	{3.14, 3.14, 2},
	{3.145, 3.15, 2},
	{1.0002, 1.0, 2},
	{3.141, 3.141, 3},
	{3.1415, 3.142, 3},
	{0.5999908, 0.60, 2},
}

func TestRoundDec64(t *testing.T) {
	for i, tr := range trs {
		if x := RoundDec64(tr.A, tr.Prec); x != tr.B {
			t.Fatal("Wrong round.", i, tr.A, x)
		}
	}
}

var tts []test = []test{
	{1.0, 1.0, 2},
	{3.145, 3.14, 2},
	{1.0002, 1.0, 2},
	{3.141, 3.141, 3},
	{3.1415, 3.141, 3},
}

func TestTruncateDec64(t *testing.T) {
	for i, tt := range tts {
		if x := TruncateDec64(tt.A, tt.Prec); x != tt.B {
			t.Fatal("Wrong round.", i, tt.A, x)
		}
	}
}

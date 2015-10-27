// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Start date:        2014-04-09
// Last modification: 2014-

package round

import (
	"testing"
)

type testround struct {
	A float64
	B float64
}

var ts []testround = []testround{
	{0.0, 0.0},
	{0.5, 1.0},
	{0.49, 0.0},
	{1.2, 1.0},
	{1.99, 2.0},
}

func TestRound64(t *testing.T) {
	for i, x := range ts {
		if r := Round64(x.A); r != x.B {
			t.Fatal("Round64 failed:", i, x.A, r)
		}
	}
}

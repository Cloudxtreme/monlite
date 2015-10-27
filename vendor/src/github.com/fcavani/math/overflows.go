// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Start date:        2014-08-30
// Last modification: 2014-

package math

import (
	"math"
)

// MulOverflows check if the multiplication will overflow.
func MulOverflows(a, b uint64) bool {
	if a <= 1 || b <= 1 {
		return false
	}
	c := a * b
	return c/b != a
}

const mostNegative = -(mostPositive + 1)
const mostPositive = 1<<63 - 1

// SignedMulOverflows check if the multiplication will overflow.
func SignedMulOverflows(a, b int64) bool {
	if a == 0 || b == 0 || a == 1 || b == 1 {
		return false
	}
	if a == mostNegative || b == mostNegative {
		return true
	}
	c := a * b
	return c/b != a
}

// SumOverflows64 check if the sum of two values will overflow
func SumOverflows64(a, b uint64) bool {
	if math.MaxUint64-a < b || math.MaxUint64-b < a {
		return true
	}
	return false
}

// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Start date:        2014-04-01
// Last modification: 2014-

package round

import (
	"math"
)

func norm(x int) int {
	if x < 0 {
		return -1 * x
	}
	return x
}

// RoundDec32 round to prec precision.
func RoundDec32(x float32, prec int) float32 {
	if prec < 1 {
		return x
	}
	dec1 := float32(math.Pow10(prec))
	dec2 := float32(math.Pow10(prec + 1))
	tmp := int(x * dec1)
	last := int(x*dec2) - tmp*10
	if norm(last) >= 5 && last >= 0 {
		tmp += 1
	} else if norm(last) >= 5 && last < 0 {
		tmp -= 1
	}
	return float32(tmp) / dec1
}

// RoundDec64 round to prec precision.
func RoundDec64(x float64, prec int) float64 {
	if prec < 1 {
		return x
	}
	dec1 := math.Pow10(prec)
	dec2 := math.Pow10(prec + 1)
	tmp := int(x * dec1)
	last := int(x*dec2) - tmp*10
	if norm(last) >= 5 && last >= 0 {
		tmp += 1
	} else if norm(last) >= 5 && last < 0 {
		tmp -= 1
	}
	return float64(tmp) / dec1
}

func TruncateDec32(x float32, prec int) float32 {
	if prec < 1 {
		return x
	}
	dec1 := float32(math.Pow10(prec))
	return float32(int(x*dec1)) / dec1
}

func TruncateDec64(x float64, prec int) float64 {
	if prec < 1 {
		return x
	}
	dec1 := math.Pow10(prec)
	return float64(int(x*dec1)) / dec1
}

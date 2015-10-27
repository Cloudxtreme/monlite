// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package unicode include more function to handle unicode chars.
package unicode

import (
	"unicode"
)

var accentsA unicode.Range16 = unicode.Range16{Lo: 0x00C0, Hi: 0x00C5, Stride: 1}
var accentsC unicode.Range16 = unicode.Range16{Lo: 0x00C7, Hi: 0x00C7, Stride: 1}
var accentsE unicode.Range16 = unicode.Range16{Lo: 0x00C8, Hi: 0x00CB, Stride: 1}
var accentsI unicode.Range16 = unicode.Range16{Lo: 0x00CC, Hi: 0x00CF, Stride: 1}
var accentsO unicode.Range16 = unicode.Range16{Lo: 0x00D2, Hi: 0x00D6, Stride: 1}
var accentsU unicode.Range16 = unicode.Range16{Lo: 0x00D9, Hi: 0x00DC, Stride: 1}

var accentsa unicode.Range16 = unicode.Range16{Lo: 0x00E0, Hi: 0x00E6, Stride: 1}
var accentsc unicode.Range16 = unicode.Range16{Lo: 0x00E7, Hi: 0x00E7, Stride: 1}
var accentse unicode.Range16 = unicode.Range16{Lo: 0x00E8, Hi: 0x00EB, Stride: 1}
var accentsi unicode.Range16 = unicode.Range16{Lo: 0x00ED, Hi: 0x00EF, Stride: 1}
var accentso unicode.Range16 = unicode.Range16{Lo: 0x00F2, Hi: 0x00F6, Stride: 1}
var accentsu unicode.Range16 = unicode.Range16{Lo: 0x00F9, Hi: 0x00FC, Stride: 1}

var accents *unicode.RangeTable = &unicode.RangeTable{
	R16: []unicode.Range16{
		accentsA,
		accentsC,
		accentsE,
		accentsI,
		accentsO,
		accentsU,
		accentsa,
		accentsc,
		accentse,
		accentsi,
		accentso,
		accentsu,
	},
}

//IsLetter check for letters and a,e,i, and u accents and ccedil.
func IsLetter(rune int32) bool {
	if unicode.IsLetter(rune) || unicode.Is(accents, rune) {
		return true
	}
	return false
}

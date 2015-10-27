// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// utf8 functions.
package utf8

import (
	"unicode/utf8"

	"github.com/fcavani/e"
)

func RuneToString(r rune) (string, error) {
	l := utf8.RuneLen(r)
	p := make([]byte, l)
	n := utf8.EncodeRune(p, r)
	if n != l {
		return "", e.New("write wrong rune")
	}
	return string(p), nil
}

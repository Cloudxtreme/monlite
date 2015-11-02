// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ping

import (
	"testing"

	"github.com/fcavani/e"
)

const imapUrl = "imap://imap.gmail.com:993"

func TestImap(t *testing.T) {
	url := testParse(t, imapUrl)
	err := PingImap(url)
	if err != nil {
		t.Fatal(e.Trace(e.Forward(err)))
	}
}

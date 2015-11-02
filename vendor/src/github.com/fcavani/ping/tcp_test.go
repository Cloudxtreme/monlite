// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ping

import (
	"testing"

	"github.com/fcavani/e"
)

const tcpUrl = "tcp://www.google.com:80"

func TestTcp(t *testing.T) {
	url := testParse(t, tcpUrl)
	err := PingTCP(url)
	if err != nil {
		t.Fatal(e.Trace(e.Forward(err)))
	}
}

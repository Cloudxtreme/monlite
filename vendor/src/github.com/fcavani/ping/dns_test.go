// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ping

import (
	"testing"

	"github.com/fcavani/e"
)

const dnsUrl = "dns://8.8.8.8/www.google.com"

func TestDns(t *testing.T) {
	url := testParse(t, dnsUrl)
	err := PingDns(url)
	if err != nil {
		t.Fatal(e.Trace(e.Forward(err)))
	}
}

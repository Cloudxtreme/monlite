// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ping

import (
	"net/url"
	"testing"

	"github.com/fcavani/e"
	myurl "github.com/fcavani/net/url"
)

func testParse(t *testing.T, rawurl string) *url.URL {
	u, err := myurl.ParseWithSocket(rawurl)
	if err != nil {
		t.Fatal(e.Trace(e.Forward(err)))
	}
	return u
}

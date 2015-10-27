// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ping

import (
	"net/url"
	"testing"

	"github.com/fcavani/e"
)

const CouchUrl = "couch://localhost:5984"

func TestPingCouch(t *testing.T) {
	url, err := url.Parse(CouchUrl)
	if err != nil {
		t.Fatal(err)
	}
	err = PingCouch(url)
	if err != nil {
		t.Fatal(e.Trace(e.Forward(err)))
	}
}

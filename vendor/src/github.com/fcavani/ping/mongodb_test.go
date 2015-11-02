// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ping

import (
	"testing"

	"github.com/fcavani/e"
)

const mongodblUrl = "mongodb://localhost/db"

func TestMongoDb(t *testing.T) {
	if !OnTravis() {
		t.Skip("not on travis")
	}
	url := testParse(t, mongodblUrl)
	err := PingMongoDb(url)
	if err != nil {
		t.Fatal(e.Trace(e.Forward(err)))
	}
}

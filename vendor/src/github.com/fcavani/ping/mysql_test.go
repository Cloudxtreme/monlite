// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ping

import (
	"testing"

	"github.com/fcavani/e"
)

const mysqlUrl = "mysql://root@tcp(127.0.0.1:3306)/db"

func TestMySql(t *testing.T) {
	if !OnTravis() {
		t.Skip("not on travis")
	}
	url := testParse(t, mysqlUrl)
	err := PingMySql(url)
	if err != nil {
		t.Fatal(e.Trace(e.Forward(err)))
	}
}

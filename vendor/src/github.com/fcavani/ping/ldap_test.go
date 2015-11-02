// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ping

import (
	"testing"

	"github.com/fcavani/e"
)

const ldapUrl = "ldap://user:pass@localhost/ou=People,dc=isp,dc=net"

func TestLdap(t *testing.T) {
	url := testParse(t, ldapUrl)
	err := PingLdap(url)
	if err != nil {
		t.Log(e.Trace(e.Forward(err)))
	}
	err = PingLdapTLS(url)
	if err != nil {
		t.Log(e.Trace(e.Forward(err)))
	}
}

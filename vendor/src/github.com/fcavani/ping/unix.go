// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ping

import (
	"net"
	"net/url"

	"github.com/fcavani/e"
)

// PingUnix try to connect to an unix socket.
func PingUnix(url *url.URL) error {
	conn, err := net.DialTimeout("unix", url.Host, DialTimeout)
	if err != nil {
		return e.New(err)
	}
	err = conn.Close()
	if err != nil {
		return e.New(err)
	}
	return nil
}

func init() {
	Add("socket", PingUnix)
	Add("unix", PingUnix)
}

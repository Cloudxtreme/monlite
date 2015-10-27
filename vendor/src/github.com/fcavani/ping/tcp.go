// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ping

import (
	"net"
	"net/url"

	"github.com/fcavani/e"
)

// PingTCP try to connect a TCP port.
func PingTCP(url *url.URL) error {
	conn, err := net.Dial("tcp", url.Host)
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
	Add("tcp", PingTCP)
}

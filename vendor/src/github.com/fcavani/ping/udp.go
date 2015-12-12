// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ping

import (
	"net"
	"net/url"
	"time"

	"github.com/fcavani/e"
)

var Deadline time.Duration = 30 * time.Second

// PingUDP try to send something to UDP port
// and receive some thing. If connections is refused
// the host is considered down. If some other error,
// besides timeout, the host is down. If timeoutr
// I can't determine if the host is there or not.
func PingUDP(url *url.URL) error {
	conn, err := net.DialTimeout("udp", url.Host, DialTimeout)
	if err != nil {
		return e.New(err)
	}
	defer conn.Close()
	err = conn.SetDeadline(time.Now().Add(Deadline))
	_, err = conn.Write([]byte("Hi!"))
	if err != nil {
		return e.New(err)
	}
	buf := make([]byte, 100)
	_, err = conn.Read(buf)
	if e.Contains(err, "connection refused") {
		return e.Push(err, ErrNotAnwsered)
	} else if err != nil && !e.Contains(err, "i/o timeout") {
		return e.New(err)
	}
	return nil
}

func init() {
	Add("udp", PingUDP)
}

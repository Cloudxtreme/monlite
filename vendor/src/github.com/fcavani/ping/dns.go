// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ping

import (
	"net/url"
	"strings"

	"github.com/fcavani/e"
	"github.com/fcavani/net/dns"
)

// PingDns test if a dns server is alive.
func PingDns(url *url.URL) error {
	url.Path = strings.Trim(url.Path, "/")
	addrs, err := dns.LookupHostWithServers(url.Path, []string{url.Host}, 5, DialTimeout)
	if err != nil {
		return e.Forward(err)
	}
	if len(addrs) == 0 {
		return e.New("query returned zero hosts")
	}
	return nil
}

func init() {
	Add("dns", PingDns)
}

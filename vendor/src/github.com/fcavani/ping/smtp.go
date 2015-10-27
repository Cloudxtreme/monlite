// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ping

import (
	gosmtp "net/smtp"
	"net/url"

	"github.com/fcavani/e"
	"github.com/fcavani/net/smtp"
)

func PingSMTP(url *url.URL) error {
	if url.Scheme != "smtp" {
		return e.New("wrong scheme")
	}
	var auth gosmtp.Auth
	if url.User != nil {
		pass, ok := url.User.Password()
		if ok {
			if url.User.Username() != "" && pass != "" {
				auth = gosmtp.PlainAuth("", url.User.Username(), pass, url.Host)
			}
		}
	}
	err := smtp.TestSMTP(url.Host, auth, "", DialTimeout, true)
	if err != nil {
		return e.Forward(err)
	}
	return nil
}

func init() {
	Add("smtp", PingSMTP)
}

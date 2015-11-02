// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ping

import (
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/fcavani/e"
	"github.com/mxk/go-imap/imap"
)

const ErrImapFailed = "imap connection failed"

func dialImap(addr string) (c *imap.Client, err error) {
	if strings.HasSuffix(addr, ":993") {
		c, err = imap.DialTLS(addr, tlsConfig)
	} else {
		c, err = imap.Dial(addr)
	}
	if err != nil {
		return nil, e.New(err)
	}
	return c, nil
}

func PingImap(url *url.URL) error {
	if url.Scheme != "imap" && url.Scheme != "imaps" {
		return e.New("not an imap/imaps scheme")
	}

	var err error

	null, err := os.OpenFile(os.DevNull, os.O_WRONLY|os.O_APPEND, 0600)
	if err == nil {
		imap.DefaultLogger = log.New(null, "", 0)
	}
	imap.DefaultLogMask = imap.LogNone
	defer null.Close()

	c, err := dialImap(url.Host)
	if err != nil {
		return e.Push(err, ErrImapFailed)
	}
	defer c.Logout(30 * time.Second)

	if c.Caps["STARTTLS"] {
		_, err := c.StartTLS(tlsConfig)
		if err != nil {
			return e.Push(err, ErrImapFailed)
		}
	}
	if c.Caps["ID"] {
		_, err := c.ID("name", "goimap")
		if err != nil {
			return e.Push(err, ErrImapFailed)
		}
	}

	cmd, err := c.Noop()
	if err != nil {
		return e.Push(err, ErrImapFailed)
	}
	rsp, err := cmd.Result(imap.OK)
	if err != nil {
		return e.Push(err, ErrImapFailed)
	}
	//log.Println(rsp.Status, rsp.Info)
	if rsp.Status != imap.OK {
		return e.New(ErrImapFailed)
	}

	if url.User == nil {
		return nil
	}
	username := url.User.Username()
	pass, ok := url.User.Password()
	if username != "" && ok {
		cmd, err := c.Login(username, pass)
		if err != nil {
			return e.Push(err, ErrImapFailed)
		}
		rsp, err := cmd.Result(imap.OK)
		if err != nil {
			return e.Push(err, ErrImapFailed)
		}
		if rsp.Status != imap.OK {
			return e.New(ErrImapFailed)
		}
	}
	return nil
}

func init() {
	Add("imap", PingImap)
	Add("imaps", PingImap)
}

// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ping

import (
	"crypto/tls"
	golog "log"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/fcavani/e"
	"github.com/fcavani/log"
	"github.com/mxk/go-imap/imap"
)

func init() {
	outer := log.Log.Store().(log.OuterLogger)
	w := outer.OuterLog(log.DebugPrio, "imap")
	imap.DefaultLogger = golog.New(w, "", 0)
	imap.DefaultLogMask = imap.LogNone
}

const ErrImapFailed = "imap connection failed"

func dialImap(addr string) (c *imap.Client, err error) {
	var conn net.Conn
	if strings.HasSuffix(addr, ":993") {
		//c, err = imap.DialTLS(addr, tlsConfig)
		conn, err = tls.DialWithDialer((&net.Dialer{Timeout: DialTimeout}), "tcp", addr, tlsConfig)
	} else {
		conn, err = net.DialTimeout("tcp", addr, DialTimeout)
	}
	if err != nil {
		return nil, e.New(err)
	}
	c, err = imap.NewClient(conn, addr, DialTimeout)
	if err != nil {
		return nil, e.New(err)
	}
	return c, nil
}

func PingImap(url *url.URL) error {
	if url.Scheme != "imap" && url.Scheme != "imaps" {
		return e.New("not an imap/imaps scheme")
	}
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

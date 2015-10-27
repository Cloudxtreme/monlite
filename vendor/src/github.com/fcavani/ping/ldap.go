// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ping

import (
	"crypto/tls"
	"net"
	"net/url"
	"regexp"

	"github.com/fcavani/e"
	utilUrl "github.com/fcavani/net/url"
	"github.com/nmcclain/ldap"
)

var conf *tls.Config

var reDn *regexp.Regexp

func init() {
	reDn = regexp.MustCompile("dc=(.*)")
}

func basedn(binddn string) (base string) {
	dns := reDn.FindAllString(binddn, -1)
	for i, val := range dns {
		base += val
		if i < len(dns)-1 {
			base += ","
		}
	}
	return
}

func pingLdap(url *url.URL, dial func(proto, addr string) (*ldap.Conn, error)) (err error) {
	var proto, addr string
	proto, addr, err = utilUrl.Socket(url.Host)
	if err != nil {
		return e.Forward(err)
	}
	var conn *ldap.Conn
	conn, err = dial(proto, addr)
	if err != nil {
		return e.Forward(err)
	}
	defer func() {
		conn.Close()
	}()
	pass, ok := url.User.Password()
	if !ok {
		err = e.New("no password")
		return
	}
	if len(url.Path) > 0 {
		url.Path = url.Path[1:]
	}
	dn := "cn=" + url.User.Username() + "," + url.Path
	err = conn.Bind(dn, pass)
	if err != nil {
		err = e.Forward(err)
		return
	}
	attrs := map[string]bool{
		"cn":        true,
		"uid":       true,
		"uidNumber": true,
		"gidNumber": true,
	}
	attrsStr := make([]string, 0, len(attrs))
	for val := range attrs {
		attrsStr = append(attrsStr, val)
	}
	search := &ldap.SearchRequest{
		BaseDN:       basedn(dn),
		Scope:        ldap.ScopeWholeSubtree,
		DerefAliases: ldap.DerefAlways,
		Filter:       "(&(objectclass=*)(cn=" + url.User.Username() + "))",
		Attributes:   attrsStr,
	}
	sr, err := conn.Search(search)
	if err != nil {
		err = e.Forward(err)
		return
	}
	if len(sr.Entries) == 1 {
		entry := sr.Entries[0]
		if entry.DN == dn {
			count := 0
			for _, attr := range entry.Attributes {
				if _, found := attrs[attr.Name]; found {
					count++
				} else {
					return e.New("ldap search returned a no requested attribute")
				}
			}
			if count != len(attrs) {
				return e.New("wrong number of attributes, required %v, got %v", len(attrs), count)
			}
		}
	}
	return nil
}

func PingLdap(url *url.URL) error {
	return e.Forward(pingLdap(url, func(proto, addr string) (*ldap.Conn, error) {
		conn, err := ldap.DialTimeout(proto, addr, DialTimeout)
		if err != nil {
			return nil, e.Forward(err)
		}
		return conn, nil
	}))
}

func PingLdapTLS(url *url.URL) error {
	return e.Forward(pingLdap(url, func(proto, addr string) (*ldap.Conn, error) {
		dialer := &net.Dialer{Timeout: DialTimeout}
		c, err := tls.DialWithDialer(dialer, proto, addr, conf)
		if err != nil {
			return nil, e.Forward(err)
		}
		conn := ldap.NewConn(c)
		return conn, nil
	}))
}

func init() {
	Add("ldap", PingLdap)
	Add("ldaptls", PingLdapTLS)
}

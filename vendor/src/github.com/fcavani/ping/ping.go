// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package Ping have functions to test if exist response from one server.
package ping

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/fcavani/e"
)

var httpClient *http.Client
var tlsConfig *tls.Config

var DialTimeout time.Duration = 60 * time.Second
var TLSHandshakeTimeout time.Duration = 60 * time.Second
var ResponseHeaderTimeout time.Duration = 60 * time.Second
var HttpTimeout time.Duration = 60 * time.Second

var transport *http.Transport = &http.Transport{
	DisableKeepAlives:     true,
	Dial:                  (&net.Dialer{Timeout: DialTimeout}).Dial,
	TLSHandshakeTimeout:   TLSHandshakeTimeout,
	ResponseHeaderTimeout: ResponseHeaderTimeout,
}

func SkipSecurityChecksTLS(b bool) {
	tlsConfig = &tls.Config{InsecureSkipVerify: b}
	transport.TLSClientConfig = tlsConfig
	httpClient = &http.Client{
		Transport: transport,
		Timeout:   HttpTimeout,
	}
}

func init() {
	httpClient = &http.Client{Transport: transport}
}

const ErrNotAnwsered = "not anwsered"

var register map[string]func(url *url.URL) error

func Add(scheme string, function func(url *url.URL) error) error {
	if register == nil {
		register = make(map[string]func(url *url.URL) error)
	}
	if _, ok := register[scheme]; ok {
		return e.New("this scheme is registered")
	}
	register[scheme] = function
	return nil
}

func PingRawUrl(rawurl string) error {
	u, err := url.Parse(rawurl)
	if err != nil {
		return e.New(err)
	}
	return e.Forward(Ping(u))
}

func Ping(url *url.URL) error {
	f, ok := register[url.Scheme]
	if !ok {
		return e.New("I don't have any function to ping a server of this scheme (%v)", url.Scheme)
	}
	err := f(url)
	if err != nil {
		return e.Forward(err)
	}
	return nil
}

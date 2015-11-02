// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ping

import (
	"net"
	"os"
	"testing"
	"time"

	"github.com/fcavani/e"
	"github.com/fcavani/rand"
)

func TestUnix(t *testing.T) {
	name, err := rand.FileName("test", "socket", 10)
	if err != nil {
		t.Fatal(e.Trace(e.Forward(err)))
	}
	//name = os.TempDir() + name
	defer os.Remove(name)
	sig := make(chan struct{})
	go func() {
		t.Log("Listen")
		addr, err := net.ResolveUnixAddr("unix", name)
		if err != nil {
			t.Fatal(e.Trace(e.Forward(err)))
		}
		ln, err := net.ListenUnix("unix", addr)
		if err != nil {
			t.Fatal(e.Trace(e.Forward(err)))
		}
		sig <- struct{}{}
		t.Log("Accepting connections on", ln.Addr())
		conn, err := ln.Accept()
		if err != nil {
			t.Fatal(e.Trace(e.Forward(err)))
		}
		t.Log("Connection accepted.", conn.RemoteAddr())
		time.Sleep(100 * time.Millisecond)
		err = conn.Close()
		if err != nil {
			t.Fatal(e.Trace(e.Forward(err)))
		}
	}()
	unixUrl := "socket://" + name
	t.Log(unixUrl)
	url := testParse(t, unixUrl)
	<-sig
	err = PingUnix(url)
	if err != nil {
		t.Fatal(e.Trace(e.Forward(err)))
	}
}

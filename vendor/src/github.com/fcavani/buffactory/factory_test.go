// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by Apache 2.0
// license that can be found in the LICENSE file.

package buffactory

import (
	"testing"
	"time"
)

func TestFactory(t *testing.T) {
	bufmaker := &BufferFactory{
		NumBuffersPerSize: 100,
		MinBuffers: 10,
		MaxBuffers: 1000,
		MinBufferSize: 256,
		MaxBufferSize: 1024,
		Reposition: 10 *time.Second,
	}
	err := bufmaker.StartBufferFactory()
	if err != nil {
		t.Fatal("StartBufferFactory failed:", err)
	}

	defer bufmaker.Close()
	buf := bufmaker.RequestBuffer(10)
	if buf == nil {
		t.Fatal("invalid buf")
	}
	if buf.Cap() != 10 {
		t.Fatal("invalid len")
	}
	bufmaker.ReturnBuffer(buf)
	if hit, miss := bufmaker.Returned(); hit != 0 && miss != 1 {
		t.Fatal("hits and misses are wrong", hit, miss)
	}
	buf = bufmaker.RequestBuffer(256)
	if buf == nil {
		t.Fatal("invalid buf")
	}
	if buf.Cap() != 256 {
		t.Fatal("invalid len")
	}
	bufmaker.ReturnBuffer(buf)
	if hit, miss := bufmaker.Returned(); hit != 1 && miss != 1 {
		t.Fatal("hits and misses are wrong", hit, miss)
	}
	
	t.Log(bufmaker.NumBuffers())
	t.Log(bufmaker.PrintBuffersStats())
}
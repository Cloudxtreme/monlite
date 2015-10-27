// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by Apache 2.0
// license that can be found in the LICENSE file.

package buffactory

import (
	"testing"
)

func TestBuffer(t *testing.T) {
	pool := newBuffer(10, 100)
	empty := false
	pool.Pop(func() {
		empty = true
	})
	if !empty {
		t.Fatal("not empty")
	}
	pool.Insert(make([]byte, 0, 10))
	pool.Insert(make([]byte, 0, 10))
	pool.Insert(make([]byte, 0, 10))
	buf := pool.Pop(nil)
	if buf == nil {
		t.Fatal("not a buffer")
	}
	pool.Insert(make([]byte, 0, 10))
	buf = pool.Pop(nil)
	if buf == nil {
		t.Fatal("not a buffer")
	}
	buf = pool.Pop(nil)
	if buf == nil {
		t.Fatal("not a buffer")
	}
	empty = false
	pool.Pop(func() {
		empty = true
	})
	if !empty {
		t.Fatal("not empty")
	}
	if buf == nil {
		t.Fatal("not a buffer")
	}
}

func TestBuffers(t *testing.T) {
	pool := NewBuffers(5, 100)
	if pool.NumBuffers() != 0 {
		t.Fatal("wrong number of buffers")
	}
	if pool.Hits() != 0 {
		t.Fatal("wrong number of hits")
	}
	if pool.Miss() != 0 {
		t.Fatal("wrong number of misses")
	}
	if hit, miss := pool.Returned(); hit != 0 || miss != 0 {
		t.Fatal("Returned returned wrong values")
	}
	
	buf := pool.Request(10)
	if buf == nil {
		t.Fatal("invalid buffer")
	}
	if pool.NumBuffers() != 0 {
		t.Fatal("wrong number of buffers")
	}
	if pool.Hits() != 0 && pool.Miss() != 1 {
		t.Fatal("hits and misses are wrong")
	}
	
	pool.Return(buf)
	if hit, miss := pool.Returned(); hit != 0 && miss != 1 {
		t.Fatal("hits and misses are wrong")
	}
	if pool.NumBuffers() != 1 {
		t.Fatal("wrong number of buffers")
	}
	t.Log(pool.PrintBuffersStats())
	
	buf = pool.Request(10)
	if buf == nil {
		t.Fatal("invalid buffer")
	}
	if pool.NumBuffers() != 0 {
		t.Fatal("wrong number of buffers")
	}
	if pool.Hits() != 1 && pool.Miss() != 1 {
		t.Fatal("hits and misses are wrong")
	}
	
	pool = NewBuffers(5, 100)
	for i := 0; i < 10; i++ {
		pool.InsertInit(make([]byte,10))
	}
	if pool.NumBuffers() != 10 {
		t.Fatal("wrong number of buffers")
	}
	if hit, miss := pool.Returned(); hit != 0 || miss != 0 {
		t.Fatal("Returned returned wrong values")
	}
	pool.Return(make([]byte, 10))
	if hit, miss := pool.Returned(); hit != 1 && miss != 0 {
		t.Fatal("hits and misses are wrong")
	}
	if pool.NumBuffers() != 11 {
		t.Fatal("wrong number of buffers")
	}
	t.Log(pool.PrintBuffersStats())
	buf = pool.Request(5)
	if buf == nil {
		t.Fatal("invalid buffer")
	}
	if len(buf) != 5 {
		t.Fatal("invalid length")
	}
	buf = pool.Request(30)
	if buf == nil {
		t.Fatal("invalid buffer")
	}
	if len(buf) != 30 {
		t.Fatal("invalid length")
	}
	
}
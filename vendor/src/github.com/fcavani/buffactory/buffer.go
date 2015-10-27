// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by Apache 2.0
// license that can be found in the LICENSE file.

package buffactory

import (
	"fmt"
	"sort"
	"sync"
)

type buffer struct {
	size    int
	buffers [][]byte
	ptr     int
}

// newBuffer creates a pull of buffers.
func newBuffer(size, numbuffers int) *buffer {
	return &buffer{
		size:    size,
		buffers: make([][]byte, 0, numbuffers),
		ptr:     -1,
	}
}

func (b *buffer) String() string {
	return fmt.Sprintf("size: %v, len: %v, ptr: %v", b.size, len(b.buffers), b.ptr)
}

func (b *buffer) Pop(onempty func()) (buf []byte) {
	if b.ptr < 0 {
		if onempty != nil {
			onempty()
		}
		return
	}
	buf = b.buffers[b.ptr]
	b.ptr--
	if b.ptr < 0 && onempty != nil {
		onempty()
	}
	return
}

func (b *buffer) Insert(buf []byte) {
	b.ptr++
	if b.ptr >= len(b.buffers) {
		b.buffers = append(b.buffers, buf)
		return
	}
	b.buffers[b.ptr] = buf
}

type Buffers interface {
	// InsertInit inserts a buffer in the pool.
	InsertInit(in []byte)
	// Return inserts a buffer in the pool.
	// If the pool is full drop the buffer, than the
	// buffer will be collected by the GC, if there is
	// no reference to it.
	Return(in []byte)
	// Request resturns a buffer with capacity size.
	Request(size int) []byte
	// NumBuffers resturn the number of buffers in the poll
	NumBuffers() int
	// Hits returns the number of buffer that was in
	// the pool.
	Hits() int
	// Miss restur the number of requested buffer that
	//wasn't in the pool
	Miss() int
	// Returned returns the number of hits, buffers that
	// are put again in the poll and have the same size,
	// and miss, buffers that haven't the same size.
	Returned() (hit, miss int)
	// ResetCounters resets the counters.
	ResetCounters()
	PrintBuffersStats() string
	GetStats() Stats
}

type buffers struct {
	Stats
	bufs       []*buffer
	lck        sync.Mutex
	numbuffers int
	maxbuffers int
	count      int
	retmiss    int
	rethit     int
	hit        int
	miss       int
}

// NewBuffers return a new struct that holds a collection of
// buffer classified by size. numbuffers are the initial number of
// buffers per size. maxbuffers are the max number of buffers for all
// sizes.
func NewBuffers(numbuffers, maxbuffers int) Buffers {
	return &buffers{
		Stats:      make(Stats, 0, NumSamples),
		bufs:       make([]*buffer, 0, numbuffers),
		numbuffers: numbuffers,
		maxbuffers: maxbuffers,
	}
}

func (b *buffers) printBuffersStats() (out string) {
	for i, buf := range b.bufs {
		out += fmt.Sprintf("%v - %v", i, buf)
		if i < len(b.bufs)-1 {
			out += "\n"
		}
	}
	out += "\n" + b.Stats.String() + "\n"
	return
}

func (b *buffers) PrintBuffersStats() (out string) {
	b.lck.Lock()
	out = b.printBuffersStats()
	b.lck.Unlock()
	return
}

func (b *buffers) InsertInit(in []byte) {
	if in == nil {
		return
	}
	c := cap(in)
	if c == 0 {
		return
	}
	b.lck.Lock()
	defer b.lck.Unlock()
	l := len(b.bufs)
	if l == 0 {
		buf := newBuffer(c, b.numbuffers)
		buf.Insert(in)
		b.bufs = []*buffer{buf}
		b.count++
		return
	} else if l >= b.maxbuffers {
		return
	}
	i := sort.Search(l, func(i int) bool {
		return b.bufs[i].size >= c
	})
	if i < l && b.bufs[i].size == c {
		b.bufs[i].Insert(in)
		b.count++
		return
	}
	buf := newBuffer(c, b.numbuffers)
	buf.Insert(in)
	b.bufs = append(b.bufs[:i], append([]*buffer{buf}, b.bufs[i:]...)...)
}

func (b *buffers) Return(in []byte) {
	if in == nil {
		return
	}
	b.InsertData(len(in))
	in = in[:cap(in)]
	c := len(in)
	if c == 0 {
		return
	}
	b.lck.Lock()
	defer b.lck.Unlock()
	l := len(b.bufs)
	if l == 0 {
		buf := newBuffer(c, b.numbuffers)
		buf.Insert(in)
		b.bufs = []*buffer{buf}
		b.retmiss++
		b.count++
		return
	} else if l >= b.maxbuffers {
		return
	}
	i := sort.Search(l, func(i int) bool {
		return b.bufs[i].size >= c
	})
	if i < l && b.bufs[i].size == c {
		b.bufs[i].Insert(in)
		b.rethit++
		b.count++
		return
	}
	buf := newBuffer(c, b.numbuffers)
	buf.Insert(in)
	b.bufs = append(b.bufs[:i], append([]*buffer{buf}, b.bufs[i:]...)...)
	b.retmiss++
	b.count++
}

func (b *buffers) delete(i int) {
	if i+1 == len(b.bufs) {
		b.bufs = b.bufs[:i]
	} else {
		b.bufs = append(b.bufs[:i], b.bufs[i+1:]...)
	}
}

func (b *buffers) pop(i, size int) []byte {
	popped := b.bufs[i].Pop(func() { b.delete(i) })
	if popped == nil {
		b.miss++
		return make([]byte, size)
	}
	b.count--
	b.hit++
	return popped
}

func resize(in []byte, size int) []byte {
	diff := size - len(in)
	if diff > 0 {
		return make([]byte, size)
	}
	return in[:size]
}

func (b *buffers) Request(size int) (buf []byte) {
	if size == 0 {
		return nil
	}
	b.lck.Lock()
	defer b.lck.Unlock()
	l := len(b.bufs)
	i := sort.Search(l, func(i int) bool {
		return b.bufs[i].size >= size
	})
	if i < l && b.bufs[i].size == size {
		buf = b.pop(i, size)
		return
	} else if i < l {
		if i+1 < l {
			i++
		}
		buf = b.pop(i, size)
		if buf == nil {
			buf = make([]byte, size)
			return
		}
		buf = resize(buf, size)
		return
	} else if l > 0 {
		buf = b.pop(l-1, size)
		if buf == nil {
			buf = make([]byte, size)
			return
		}
		buf = resize(buf, size)
		return
	}
	b.miss++
	buf = make([]byte, size)
	return
}

// NumBuffers return the total number of buffers
func (b *buffers) NumBuffers() int {
	b.lck.Lock()
	defer b.lck.Unlock()
	return b.count
}

// Hits: total number of buffer serverd that was in cache.
// This include buffers with length grater than the requested.
func (b *buffers) Hits() int {
	b.lck.Lock()
	defer b.lck.Unlock()
	return b.hit
}

// Miss: number of serverd buffer that was created.
func (b *buffers) Miss() int {
	b.lck.Lock()
	defer b.lck.Unlock()
	return b.miss
}

// Returned number of buffer retorned to the pool.
// hit is for buffers of same size and miss for
// diferent size.
func (b *buffers) Returned() (hit, miss int) {
	b.lck.Lock()
	hit = b.rethit
	miss = b.retmiss
	b.lck.Unlock()
	return
}

func (b *buffers) ResetCounters() {
	b.lck.Lock()
	defer b.lck.Unlock()
	b.hit = 0
	b.miss = 0
	b.rethit = 0
	b.retmiss = 0
}

func (b *buffers) GetStats() Stats {
	return b.Stats
}

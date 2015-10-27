// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by Apache 2.0
// license that can be found in the LICENSE file.

package buffactory

import (
	"bytes"
	"time"
	"math"

	"github.com/fcavani/e"
)

type BufFactory interface {
	Buffers
	RequestBuffer(size int) (buf *bytes.Buffer)
	ReturnBuffer(buf *bytes.Buffer)
	Close()
}

type BufferFactory struct {
	Buffers
	//NumBuffersPerSize are the initial number of buffers per size.
	NumBuffersPerSize int
	// MinBuffers minimal number of buffer per size.
	MinBuffers        int
	// MaxBuffers are the max number of buffer for all sizes.
	MaxBuffers        int
	// MinBufferSize minimal size of a buffer when it is allocated.
	MinBufferSize     int
	// MaxBufferSize max size of a buffer when it is auto allocated.
	MaxBufferSize     int
	// Reposition is the periode of time when the buffer is filled again.
	Reposition        time.Duration
	chclose           chan struct{}
	reposition        int
}

// StartBufferFactory must be called after the cration of BufferFactory
// and before the others functions.
func (bf *BufferFactory) StartBufferFactory() error {
	if bf.NumBuffersPerSize <= 0 {
		return e.New("NumBuffersPerSize must be greater than zero")
	}
	if bf.MinBuffers < 0 {
		return e.New("MinBuffers must be greater or equal to zero")
	}
	if bf.MaxBuffers <= 0 {
		return e.New("MaxBuffers must be greater than zero")
	}
	if bf.MinBufferSize <= 0 {
		return e.New("MinBufferSize must be greater than zero")
	}
	if bf.MaxBufferSize <= 0 {
		return e.New("MinBufferSize must be greater than zero")
	}
	bf.Buffers = NewBuffers(bf.NumBuffersPerSize, bf.MaxBuffers)
	for i := 0; i < bf.MinBuffers; i++ {
		bf.InsertInit(make([]byte, 0, bf.MinBufferSize))
	}
	bf.Buffers.ResetCounters()
	bf.chclose = make(chan struct{})
	if bf.Reposition != 0 {
		go func() {
			for {
				select {
				case <-bf.chclose:
					return
				case <-time.After(bf.Reposition):
					d := bf.MinBuffers - bf.NumBuffers()
					if d > 0 {
						size := bf.MinBufferSize
						avg := int(math.Ceil(bf.GetStats().Average()))
						if avg > size && avg <= bf.MaxBufferSize {
							size = avg
						}
						for i := 0; i < d; i++ {
							bf.InsertInit(make([]byte, 0, size))
							bf.reposition++
						}
					}
				}
			}
		}()
	}
	return nil
}

// RequestBuffer returns a bytes.Buffer with size size.
func (bf *BufferFactory) RequestBuffer(size int) *bytes.Buffer {
	buf := bf.Request(size)
	return bytes.NewBuffer(buf[:0])
}

// ReturnBuffer is like Return but the argument is bytes.Buffer
func (bf *BufferFactory) ReturnBuffer(buf *bytes.Buffer) {
	bf.Return(buf.Bytes())
}

// Close closes the factory
func (bf *BufferFactory) Close() {
	if bf.Reposition != 0 {
		bf.chclose <- struct{}{}
	}
}

// ResetCounters resets the counters.
func (bf *BufferFactory) ResetCounters() {
	bf.reposition = 0
	bf.Buffers.ResetCounters()
}

// RepositionCount return the number of buffer repositions.
func (bf *BufferFactory) RepositionCount() int {
	return bf.reposition
}

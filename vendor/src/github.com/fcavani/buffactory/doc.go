// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by Apache 2.0
// license that can be found in the LICENSE file.

// Buffactory is a pre-allocated pool of buffers.
// Create the BufferFactory struct first. The values in
// the struct is application dependent, you have to find
// it ourselve.
//
// 	bufmaker = &buffactory.BufferFactory{
//		NumBuffersPerSize: 100,
//		MinBuffers: 10,
//		MaxBuffers: 1000,
//		MinBufferSize: 1024,
//		MaxBufferSize: 4096,
// 		Reposition: 10 *time.Second,
//	}
//
//
// Than init the struct:
// 
//	err := bufmaker.StartBufferFactory()
//	if err != nil {
//		...
// 	}
// 
// To request one buffer:
// 
// 	buf := bufmaker.Request(1024)
// 
//
// To return this buffer to the pool:
// 
// 	bufmaker.Return(buf)
// 
// Remember to close it in the end:
// 
// 	bufmaker.Close()
// 
package buffactory

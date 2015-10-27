// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by Apache 2.0
// license that can be found in the LICENSE file.

package buffactory

import (
	"testing"
)

func TestInsertData(t *testing.T) {
	stats := make(Stats, 0, 3)
	stats.InsertData(0)
	stats.InsertData(1)
	stats.InsertData(2)
	stats.InsertData(3)
	for i, s := range stats {
		if int64(i) != s {
			t.Fatal("numbers don't match", i , s, []int64(stats))
		}
	}
}

func TestStats(t *testing.T) {
	stats := make(Stats, 0, 3)
	stats.InsertData(1)
	stats.InsertData(2)
	stats.InsertData(3)
	
	if stats.Min() != 1 {
		t.Fatal("Min failed", stats.Min())
	}
	
	if stats.Max() != 3 {
		t.Fatal("Max failed", stats.Max())
	}
	
	if stats.Average() != 2.0 {
		t.Fatal("Average failed", stats.Average())
	}
	
	if stats.StdDev() != 1.0 {
		t.Fatal("StdDev failed", stats.StdDev())
	}
}
// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by Apache 2.0
// license that can be found in the LICENSE file.

package buffactory

import (
	"fmt"
	
	"github.com/fcavani/math"
)

var NumSamples int = 1000

type Stats []int64

func (s Stats) String() string {
	return fmt.Sprintf("Samples: %v, Min: %v, Mean: %.2f (%.2f), Max: %v",
		len(s),
		s.Min(),
		s.Average(),
		s.StdDev(),
		s.Max(),
	)
}

func (s *Stats) InsertData(val int) {
	st := *s
	if len(st) >= NumSamples {
		if len(st) == cap(st) {
			st = st[1:]
		}
	}
	st = append(st, int64(val))
	*s = st
}

func (s Stats) Average() float64 {
	return math.AvgInt64(s)
}

func (s Stats) Min() int64 {
	return math.MinInt64(s)
}

func (s Stats) Max() int64 {
	return math.MaxInt64(s)
}

func (s Stats) StdDev() float64 {
	return math.StdDevInt64(s)
}


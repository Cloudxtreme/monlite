// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import (
	"math"
	"time"
)

func MinDuration(values []time.Duration) time.Duration {
	var min time.Duration = time.Duration(math.MaxInt64)
	for _, val := range values {
		if val < min {
			min = val
		}
	}
	return min
}

func MaxDuration(values []time.Duration) time.Duration {
	var max time.Duration = 0
	for _, val := range values {
		if val > max {
			max = val
		}
	}
	return max
}

func AvgDuration(values []time.Duration) float64 {
	var sum time.Duration
	for _, val := range values {
		sum += val
	}
	return float64(sum) / float64(len(values))
}

func StdDevDuration(values []time.Duration) float64 {
	// StdDevDuration: amostral
	mean := AvgDuration(values)
	var sum float64
	for _, val := range values {
		sum += (float64(val) - mean) * (float64(val) - mean)
	}
	return math.Sqrt((1.0 / float64(len(values)-1)) * sum)
}

func MinInt64(values []int64) int64 {
	var min int64 = math.MaxInt64
	for _, val := range values {
		if val < min {
			min = val
		}
	}
	return min
}

func MaxInt64(values []int64) int64 {
	var max int64 = math.MinInt64
	for _, val := range values {
		if val > max {
			max = val
		}
	}
	return max
}

func AvgInt64(values []int64) float64 {
	var sum int64
	for _, val := range values {
		sum += val
	}
	return float64(sum) / float64(len(values))
}

func StdDevInt64(values []int64) float64 {
	// StdDevDuration: amostral
	mean := AvgInt64(values)
	var sum float64
	for _, val := range values {
		sum += (float64(val) - mean) * (float64(val) - mean)
	}
	return math.Sqrt((1.0 / float64(len(values)-1)) * sum)
}

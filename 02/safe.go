package main

import (
	. "math"
)

func safe(levels []float64) bool {
	// nasty trick to iterate len(levels) - 2 times, i ranges from [0 .. len-1)
	for i := range levels[1:] {
		level1 := levels[i]
		level2 := levels[i+1]
		diff := Abs(level1 - level2)
		if diff < 1 || diff > 3 || levels[1] > levels[0] != (level2 > level1) {
			return false
		}
	}
	return true
}

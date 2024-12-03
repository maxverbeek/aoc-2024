package main

import (
	. "bufio"
	. "math"
	. "os"
	"slices"
	. "strconv"
	"strings"
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

func main() {
	scanner := NewScanner(Stdin)

	var part1, part2 int

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		var levels []float64

		for _, field := range fields {
			l, _ := ParseFloat(field, 64)
			levels = append(levels, l)
		}

		if safe(levels) {
			part1++
			continue
		}

		// not safe, so try deleting random elements and test again
		for i := range levels {
			if safe(append(slices.Clone(levels[:i]), levels[i+1:]...)) {
				part2++
				break
			}
		}
	}

	println(part1, part1+part2)
}

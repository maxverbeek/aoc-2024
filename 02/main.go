package main

import (
	. "bufio"
	. "os"
	. "slices"
	. "strconv"
	"strings"
)

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

		// the safe function is defined in a separate file. Due to the rules of
		// the token counter I can do this free of charge, and this allows me
		// to de-qualify 1 additional import (slices) due to not importing math
		// unqualified here.
		if safe(levels) {
			part1++
			continue
		}

		// not safe, so try deleting random elements and test again
		for i := range levels {
			if safe(append(Clone(levels[:i]), levels[i+1:]...)) {
				part2++
				break
			}
		}
	}

	println(part1, part1+part2)
}

package main

import (
	_ "embed"
	"fmt"
	"sort"
	"strings"
)

//go:embed input.txt
var input string

func main() {
	length := strings.Count(input, "\n") + 1
	first := make([]int, length)
	second := make([]int, length)
	appearances := make(map[int]int)

	for i, line := range strings.Split(input, "\n") {
		right := &second[i]
		fmt.Sscanf(line, "%d   %d", &first[i], right)

		if count, ok := appearances[*right]; ok {
			appearances[*right] = count + 1
		} else {
			appearances[*right] = 1
		}
	}

	sort.Ints(first)
	sort.Ints(second)

	sum := 0
	sim := 0

	for i, left := range first {
		right := second[i]

		sum += max(left-right, right-left)
		sim += left * appearances[left]
	}

	println("Part 1: ", sum)
	println("Part 2: ", sim)
}

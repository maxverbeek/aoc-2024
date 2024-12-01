package main

import (
	_ "embed"
	// super nasty, but not having to type fmt. or sort. etc saves 2 tokens each time i use these functions
	. "fmt"
	. "sort"
	. "strings"
)

//go:embed input.txt
var input string

func main() {
	length := Count(input, "\n") + 1
	first := make([]int, length)
	second := make([]int, length)
	appearances := make(map[int]int)

	for i, line := range Split(input, "\n") {
		// fmt.Sscanf
		Sscanf(line, "%d   %d", &first[i], &second[i])

		appearances[second[i]] += 1
	}

	// sort.Ints
	Ints(first)
	Ints(second)

	sum := 0
	sim := 0

	for i, left := range first {
		right := second[i]

		sum += max(left-right, right-left)
		sim += left * appearances[left]
	}

	println(sum)
	println(sim)
}

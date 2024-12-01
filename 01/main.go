package main

import (
	// super nasty, but not having to type fmt. or sort. etc saves 2 tokens each time i use these functions
	. "fmt"
	. "sort"
)

func main() {
	var first, second []int
	appearances := map[int]int{}

	var left, right int

	for {
		// this if-statement is necessary because scanf has 2 return values, and you cannot use 2 return values in a loop
		if _, err := Scanf("%d   %d\n", &left, &right); err != nil {
			break
		}
		first, second = append(first, left), append(second, right)
		appearances[right]++
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

	Printf("%d\n%d\n", sum, sim)
}

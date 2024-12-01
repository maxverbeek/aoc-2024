package main

import (
	// super nasty, but not having to type fmt. or sort. etc saves 2 tokens each time i use these functions
	. "fmt"
	. "math"
	. "sort"
)

func main() {
	var first, second []int
	appearances := map[int]int{}

	var left, right, sum, sim int
	var err error

	for err == nil {
		// this if-statement is necessary because scanf has 2 return values, and you cannot use 2 return values in a loop
		_, err = Scanf("%d   %d\n", &left, &right)
		first, second = append(first, left), append(second, right)
		appearances[right]++
	}

	// sort.Ints
	Ints(first)
	Ints(second)

	for i, left := range first {
		sum += int(Abs(float64(left - second[i])))
		sim += left * appearances[left]
	}

	println(sum, sim)
}

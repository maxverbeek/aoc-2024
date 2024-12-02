package main

import (
	// super nasty, but not having to type fmt. or sort. etc saves 2 tokens each time i use these functions
	. "fmt"
	. "math"
	. "sort"
)

func main() {
	var first, second []float64
	appearances := map[float64]int{}

	var left, right, sum float64
	var sim int
	var err error

	for err == nil {
		// this if-statement is necessary because scanf has 2 return values, and you cannot use 2 return values in a loop
		_, err = Scanf("%f   %f\n", &left, &right)
		first, second = append(first, left), append(second, right)
		appearances[right]++
	}

	// sort.Ints
	Float64s(first)
	Float64s(second)

	for i, left := range first {
		sum += Abs(left - second[i])
		sim += int(left) * appearances[left]
	}

	// 122 tokens with scientific notation
	// println(sum, sim)

	// 124 tokens without scientific notation
	Printf("%.0f\n%d\n", sum, sim)
}

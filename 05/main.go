package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

func main() {
	input, _ := io.ReadAll(os.Stdin)
	halves := strings.Split(string(input), "\n\n")
	ordering := map[int]map[int]int{}

	var left, right int

	for _, line := range strings.Split(halves[0], "\n") {
		fmt.Sscanf(line, "%d|%d", &left, &right)
		if ordering[left] == nil {
			ordering[left] = map[int]int{}
		}

		if ordering[right] == nil {
			ordering[right] = map[int]int{}
		}

		ordering[left][right] = -1
		ordering[right][left] = 1
	}

	var part1, part2 int

	for _, line := range strings.Split(strings.TrimSpace(halves[1]), "\n") {
		numstrings := strings.Split(line, ",")
		nums := make([]int, len(numstrings))

		for i, numstring := range numstrings {
			fmt.Sscanf(numstring, "%d", &nums[i])
		}

		copied := slices.Clone(nums)

		slices.SortFunc(nums, func(left, right int) int {
			if ordering[left] == nil {
				return 0
			}

			return ordering[left][right]
		})

		if slices.Equal(copied, nums) {
			part1 += nums[len(nums)/2]
		} else {
			part2 += nums[len(nums)/2]
		}
	}

	println(part1, part2)
}

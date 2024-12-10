package main

import (
	"bufio"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	world := [][]int{}

	for scanner.Scan() {
		line := scanner.Text()
		nums := make([]int, len(line))
		for i, rune := range line {
			nums[i] = int(rune - '0')
		}

		world = append(world, nums)
	}

	part1 := 0

	for y := range world {
		for x := range world[y] {
			if world[y][x] == 0 {
				count := WalkUphill(-1, x, y, world)
				println(x, y, count)
				part1 += count
			}
		}
	}

	println(part1)
}

func WalkUphill(previous, x, y int, world [][]int) int {
	current := world[y][x]

	if current != previous+1 {
		return 0
	}

	if current == 9 {
		return 1
	}

	var up, down, left, right int

	if y-1 >= 0 {
		up = WalkUphill(current, x, y-1, world)
	}

	if y+1 < len(world) {
		down = WalkUphill(current, x, y+1, world)
	}

	if x-1 >= 0 {
		left = WalkUphill(current, x-1, y, world)
	}

	if x+1 < len(world[y]) {
		right = WalkUphill(current, x+1, y, world)
	}

	sum := up + down + left + right

	// ????
	if sum > 0 {
		return 1
	} else {
		return 0
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

type Robot struct {
	x, y, dx, dy int
}

func (r *Robot) String() string {
	return fmt.Sprintf("p=%d,%d v=%d,%d (q %d)", r.x, r.y, r.dx, r.dy, r.Quadrant(11, 7))
}

func (r Robot) Step(width, height, time int) Robot {
	r.x = ((r.x+r.dx*time)%width + width) % width
	r.y = ((r.y+r.dy*time)%height + height) % height

	return r
}

func (r *Robot) Quadrant(width, height int) int {
	mx, my := width/2, height/2

	if width%2 == 0 || height%2 == 0 {
		panic("needs even width/height")
	}

	// Quadrants:
	// 1 2
	// 3 4
	// if on a divider: 0

	if r.x < mx {
		if r.y < my {
			return 1
		} else if r.y > my {
			return 3
		}
	} else if r.x > mx {
		if r.y < my {
			return 2
		} else if r.y > my {
			return 4
		}
	}

	return 0
}

func Draw(robots []Robot, width, height int) {
	grid := make([][]int, height)
	for y := 0; y < height; y++ {
		grid[y] = make([]int, width)
	}

	for _, robot := range robots {
		grid[robot.y][robot.x]++
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if grid[y][x] > 0 {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}

		println()
	}
}

func Quadrants(robots []Robot, width, height int) [5]int {
	quadrants := [5]int{0, 0, 0, 0, 0}

	for _, robot := range robots {
		// fmt.Printf("robot %2d: quadrant %d\n", i, robot.Quadrant(width, height))
		quadrants[robot.Quadrant(width, height)]++
	}

	return quadrants
}

func StepAll(robots []Robot, width, height, steps int) []Robot {
	result := make([]Robot, len(robots))
	for i, r := range robots {
		result[i] = r.Step(width, height, steps)
	}

	return result
}

func IsSymmetric(robots []Robot, width, height, threshold, verticalBleed int) bool {
	grid := make([][]int, height)
	for y := 0; y < height; y++ {
		grid[y] = make([]int, width)
	}

	for _, robot := range robots {
		grid[robot.y][robot.x] = 1
	}

	mistakes := 0
	for y := verticalBleed; y < height-verticalBleed; y++ {
		for x := 0; x < width/2; x++ {
			if grid[y][width/2-x] != grid[y][width/2+x] {
				mistakes++
			}
		}
	}

	return mistakes < threshold
}

func Variance(robots []Robot) (int, int) {
	avgx, avgy := 0, 0

	for _, robot := range robots {
		avgx += robot.x
		avgy += robot.y
	}

	avgx, avgy = avgx/len(robots), avgy/len(robots)

	varx, vary := 0, 0

	for _, robot := range robots {
		varx += (avgx - robot.x) * (avgx - robot.x)
		vary += (avgy - robot.y) * (avgy - robot.y)
	}

	varx, vary = varx/len(robots), vary/len(robots)

	return varx, vary
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	original := []Robot{}

	for scanner.Scan() {
		robot := Robot{}
		fmt.Sscanf(scanner.Text(), "p=%d,%d v=%d,%d\n", &robot.x, &robot.y, &robot.dx, &robot.dy)

		original = append(original, robot)
	}

	// width, height := 11, 7
	width, height := 101, 103

	robots := StepAll(original, width, height, 100)

	part1 := 1
	quadrants := Quadrants(robots, width, height)
	for _, count := range quadrants[1:] {
		part1 *= count
	}

	minvar := 9999999999999
	minvarat := 0
	minvarrobots := robots

	// I randomly guessed 500k and saw that it is too high, so this gave me an upperbound for how far I need to search
	for s := 100; s < 500_000; s++ {
		robots = StepAll(original, width, height, s)

		// if IsSymmetric(robots, width, height, len(robots)/4, 10) {
		// 	println(s)
		// 	Draw(robots, width, height)
		// }

		varx, vary := Variance(robots)

		if varx*vary < minvar {
			minvar = varx * vary
			minvarat = s
			minvarrobots = slices.Clone(robots)
		}

		if s%100_000 == 0 {
			fmt.Printf("minimum variance %d at s = %d\n", minvar, minvarat)
			Draw(minvarrobots, width, height)
		}
	}

	println(part1)
}

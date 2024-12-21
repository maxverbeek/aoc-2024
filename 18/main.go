package main

import (
	"bufio"
	"fmt"
	"os"
)

type Coordinate struct {
	x, y int
}

func (c Coordinate) Neighbours(size int) []Coordinate {
	neighbours := []Coordinate{}

	xs := [4]int{-1, 0, 1, 0}
	ys := [4]int{0, 1, 0, -1}

	for i := range xs {
		dx, dy := xs[i], ys[i]

		if c.y+dy < 0 || c.y+dy >= size {
			continue
		}
		if c.x+dx < 0 || c.x+dx >= size {
			continue
		}
		neighbours = append(neighbours, Coordinate{c.x + dx, c.y + dy})
	}

	return neighbours
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	coordinates := []Coordinate{}

	for scanner.Scan() {
		coord := Coordinate{}
		fmt.Sscanf(scanner.Text(), "%d,%d", &coord.x, &coord.y)
		coordinates = append(coordinates, coord)
	}

	var size, fall int

	if len(coordinates) > 1000 {
		size = 71
		fall = 1024
	} else {
		size = 7
		fall = 12
	}

	grid := make([][]bool, size)

	for i := range grid {
		grid[i] = make([]bool, size)
	}

	for _, coord := range coordinates[:fall] {
		grid[coord.y][coord.x] = true
	}

	steps := Bfs(grid, size, Coordinate{0, 0})

	println(steps)

	// can do this more efficiently by pathfinding in 8-connectivity from
	// bottom left to topright on the blocked squares. if that becomes possible
	// then we can no longer move across.
	//
	// alternatively, we only need to re-pathfind every time a square falls on
	// our current path, but this current algorithm terminates within 10
	// seconds so i odn't see much reason to optimize it.
	lo, hi := fall, len(coordinates)
	for lo < hi {
		mid := (lo + hi) / 2
		println(lo, mid, hi)

		// drop the grid again
		grid := make([][]bool, size)

		for i := range grid {
			grid[i] = make([]bool, size)
		}

		for _, coord := range coordinates[:mid] {
			grid[coord.y][coord.x] = true
		}

		if Bfs(grid, size, Coordinate{0, 0}) == -1 {
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}

	fmt.Printf("%d,%d\n", coordinates[lo-1].x, coordinates[lo-1].y)
}

func Bfs(grid [][]bool, size int, start Coordinate) int {
	queue := []Coordinate{start}

	steps := map[Coordinate]int{}
	prev := map[Coordinate]Coordinate{}
	steps[start] = 1

	for len(queue) != 0 {
		item := queue[0]
		queue = queue[1:] // cannot use := here because it shadows the variable
		// fmt.Printf("q: %d, thinking about %d, %d\n", len(queue), item.x, item.y)

		for _, neighbour := range item.Neighbours(size) {
			if !grid[neighbour.y][neighbour.x] && steps[neighbour] == 0 {
				// fmt.Printf("appending %d, %d\n", neighbour.x, neighbour.y)
				queue = append(queue, neighbour)
				steps[neighbour] = steps[item] + 1
				prev[neighbour] = item
			}
		}

		// fmt.Printf("q: %d, thought about %d, %d\n", len(queue), item.x, item.y)
	}

	last := Coordinate{x: size - 1, y: size - 1}
	first := Coordinate{x: 0, y: 0}

	for prev[last] != first {
		last = prev[last]
	}

	return steps[Coordinate{x: size - 1, y: size - 1}] - 1
}

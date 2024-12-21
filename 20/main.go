package main

import (
	"io"
	"os"
	"slices"
	"strings"
)

type Location struct {
	x, y int
}

type Shortcut struct {
	start, end Location
}

func (c Location) Neighbours(width, height int) []Location {
	neighbours := []Location{}

	xs := [4]int{-1, 0, 1, 0}
	ys := [4]int{0, 1, 0, -1}

	for i := range xs {
		dx, dy := xs[i], ys[i]

		if c.y+dy < 0 || c.y+dy >= height {
			continue
		}
		if c.x+dx < 0 || c.x+dx >= width {
			continue
		}
		neighbours = append(neighbours, Location{c.x + dx, c.y + dy})
	}

	return neighbours
}

func main() {
	input, _ := io.ReadAll(os.Stdin)
	maze := [][]int{}

	var start, end Location

	for y, line := range strings.Split(strings.TrimSpace(string(input)), "\n") {
		maze = append(maze, make([]int, len(line)))

		for x, rune := range line {
			if rune == '#' {
				maze[y][x] = -1
			} else {
				maze[y][x] = 0
				if rune == 'S' {
					start = Location{x: x, y: y}
				} else if rune == 'E' {
					end = Location{x: x, y: y}
				}
			}

		}
	}

	path := Bfs(maze, start, end)

	part1, part2 := 0, 0

	for _, pos := range path {
		part1 += FindShortcuts(maze, pos, 2, 100)
		part2 += FindShortcuts(maze, pos, 20, 100)
	}

	println(part1, part2)
}

func FindShortcuts(grid [][]int, pos Location, cheatrange, threshold int) int {
	cheats := 0

	for dy := -cheatrange; dy <= cheatrange; dy++ {
		for dx := -cheatrange; dx <= cheatrange; dx++ {
			manhattan := max(-dy, dy) + max(-dx, dx)
			if manhattan > cheatrange {
				continue
			}

			x, y := pos.x+dx, pos.y+dy

			if y < 0 || y >= len(grid) || x < 0 || x >= len(grid[y]) {
				continue
			}

			// the difference in times, minus the time it would take to walk in cheatmode
			cheatdistance := grid[y][x] - grid[pos.y][pos.x] - manhattan

			if cheatdistance >= threshold {
				cheats++
			}
		}

	}

	return cheats
}

func Bfs(grid [][]int, start, end Location) []Location {
	height, width := len(grid), len(grid[0])
	queue := []Location{start}

	prev := map[Location]Location{}

	for len(queue) != 0 {
		item := queue[0]
		queue = queue[1:] // cannot use := here because it shadows the variable

		for _, neighbour := range item.Neighbours(width, height) {
			if grid[neighbour.y][neighbour.x] == 0 {
				queue = append(queue, neighbour)
				prev[neighbour] = item
				grid[neighbour.y][neighbour.x] = grid[item.y][item.x] + 1
			}
		}
	}

	delete(prev, start)
	grid[start.y][start.x] = 0

	// reconstruct path
	path := []Location{}

	pos := end

	for {
		path = append(path, pos)

		if p, ok := prev[pos]; ok {
			pos = p
		} else {
			break
		}
	}

	slices.Reverse(path)

	return path
}
